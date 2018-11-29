package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Import Database Migrate Postgres suppose
	"github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import Postgres Support
	"github.com/rs/xid"
	config "github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/snowzach/gorestapi/conf"
	"github.com/snowzach/gorestapi/store/postgres/migrations"
)

// Client is the database client
type Client struct {
	logger *zap.SugaredLogger
	dbName string
	db     *sqlx.DB
	newID  func() string
}

// New returns a new database client
func New() (*Client, error) {

	logger := zap.S().With("package", "storage.psql")

	var err error

	var dbCreds string       // Regular credentials
	var dbCreateCreds string // Optional admin credentials to create database
	var dbURL string         // Postgres connect url for db migration tool
	var dbURLOptions string  // OPtions for connection

	// Username
	if username := config.GetString("storage.username"); username != "" {
		dbCreds = username + ":" + config.GetString("storage.password")
	} else {
		return nil, fmt.Errorf("No username specified")
	}

	// Check to see if we have an admin password specified (that we will use to create the database if it does not exist)
	pgPassword := config.GetString("POSTGRES_PASSWORD")
	if pgPassword != "" {
		pgUser := config.GetString("POSTGRES_USER")
		if pgUser == "" {
			pgUser = "postgres"
		}
		dbCreateCreds = fmt.Sprintf("%s:%s", pgUser, pgPassword)
	} else {
		dbCreateCreds = dbCreds
	}

	// Host + Port
	if hostname := config.GetString("storage.host"); hostname != "" {
		dbURL += "@" + hostname
	} else {
		return nil, fmt.Errorf("No hostname specified")
	}
	if port := config.GetString("storage.port"); port != "" {
		dbURL += ":" + port
	}

	// Database Name
	dbName := config.GetString("storage.database")
	if dbName == "" {
		return nil, fmt.Errorf("No database specified")
	}

	// SSL Mode
	if sslMode := config.GetString("storage.sslmode"); sslMode != "" {
		// dbConnection += fmt.Sprintf("sslmode=%s ", sslMode)
		dbURLOptions += fmt.Sprintf("?sslmode=%s", sslMode)
	}

	for retries := config.GetInt("storage.retries"); retries > 0 && !conf.StopFlag; retries-- {
		createDb, err := sql.Open("postgres", "postgres://"+dbCreateCreds+dbURL+dbURLOptions)
		// Attempt to create the database if it doesn't exist
		if err == nil {
			defer createDb.Close()
			// See if it exists
			var one sql.NullInt64
			err = createDb.QueryRow(`SELECT 1 from pg_database WHERE datname=$1`, dbName).Scan(&one)
			if err == nil {
				break // already exists
			} else if err != sql.ErrNoRows && !strings.Contains(err.Error(), "does not exist") {
				// Some other error besides does not exist
				return nil, fmt.Errorf("Could not check for database: %s", err)
			}
		} else if strings.Contains(err.Error(), "permission denied") {
			return nil, fmt.Errorf("Could not connect to database: %s", err)
		} else if strings.Contains(err.Error(), "connection refused") {
			logger.Warnw("Connection to database timed out. Sleeping and retry.",
				"storage.host", config.GetString("storage.host"),
				"storage.username", config.GetString("storage.username"),
				"storage.password", "****",
				"storage.port", config.GetInt("storage.port"),
			)
			time.Sleep(config.GetDuration("storage.sleep_between_retries"))
			continue
		}
		logger.Infow("Creating database", "database", dbName)
		_, err = createDb.Exec(`CREATE DATABASE ` + dbName)
		if err != nil {
			return nil, fmt.Errorf("Could not create database: %s", err)
		}
		break
	}

	// Build the full DB URL
	fullDbURL := "postgres://" + dbCreds + dbURL + "/" + dbName + dbURLOptions

	// If we caught the stop flag while sleeping
	if conf.StopFlag {
		return nil, fmt.Errorf("Database connection aborted")
	}

	// Still not connected?
	if err != nil {
		return nil, fmt.Errorf("Could not connect to database: %s", err)
	}

	// Make the connection using the sqlx connector now that we know the database exists
	db, err := sqlx.Connect("postgres", fullDbURL)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to database: %s", err)
	}

	// Ping the database
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Could not ping database %s", err)
	}

	db.SetMaxOpenConns(config.GetInt("storage.max_connections"))

	logger.Debugw("Connected to database server",
		"storage.host", config.GetString("storage.host"),
		"storage.username", config.GetString("storage.username"),
		"storage.port", config.GetInt("storage.port"),
		"storage.database", config.GetString("storage.database"),
	)

	c := &Client{
		logger: logger,
		dbName: dbName,
		db:     db,
		newID: func() string {
			return xid.New().String()
		},
	}

	// wrap assets into Resource
	s := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})
	d, err := bindata.WithInstance(s)
	if err != nil {
		return nil, fmt.Errorf("Could not create migrations reader: %v", err)
	}
	m, err := migrate.NewWithSourceInstance("go-bindata", d, fullDbURL)
	if err != nil {
		logger.Errorw("Database migration error",
			"error", err,
		)
		return nil, fmt.Errorf("Migrate Error:%s", err)
	}

	// Do we wipe the database
	if config.GetBool("storage.wipe_confirm") {
		err = m.Down()
		if err == migrate.ErrNoChange {
			// Okay
		} else if err != nil {
			logger.Errorw("Migrate Database Down Error",
				"Error", err,
			)
			return nil, fmt.Errorf("Migrate Error:%s", err)
		} else {
			logger.Warn("Database wipe complete...")
		}
	}

	// Perform the migration up
	err = m.Up()
	if err == migrate.ErrNoChange {
		logger.Info("Database schmea current")
	} else if err != nil {
		logger.Errorw("Migrate Error",
			"error", err,
		)
		return nil, fmt.Errorf("Migrate Error:%s", err)
	} else {
		logger.Info("Database migration completed")
	}

	return c, nil

}
