package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	// Import Database Migrate Postgres suppose
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/snowzach/gorestapi/conf"
	"github.com/snowzach/gorestapi/embed"
	"github.com/snowzach/gorestapi/store"
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

	logger := zap.S().With("package", "store.postgres")

	var err error

	var dbCreds string       // Regular credentials
	var dbCreateCreds string // Optional admin credentials to create database
	var dbURL string         // Postgres connect url for db migration tool
	var dbURLOptions string  // OPtions for connection

	// Username
	if username := conf.C.String("database.username"); username != "" {
		dbCreds = username + ":" + conf.C.String("database.password")
	} else {
		return nil, fmt.Errorf("No username specified")
	}

	// Check to see if we have an admin password specified (that we will use to create the database if it does not exist)
	pgPassword := conf.C.String("POSTGRES_PASSWORD")
	if pgPassword != "" {
		pgUser := conf.C.String("POSTGRES_USER")
		if pgUser == "" {
			pgUser = "postgres"
		}
		dbCreateCreds = fmt.Sprintf("%s:%s", pgUser, pgPassword)
	} else {
		dbCreateCreds = dbCreds
	}

	// Host + Port
	if hostname := conf.C.String("database.host"); hostname != "" {
		dbURL += "@" + hostname
	} else {
		return nil, fmt.Errorf("No hostname specified")
	}
	if port := conf.C.String("database.port"); port != "" {
		dbURL += ":" + port
	}

	// Database Name
	dbName := conf.C.String("database.database")
	if dbName == "" {
		return nil, fmt.Errorf("No database specified")
	}

	// SSL Mode
	if sslMode := conf.C.String("database.sslmode"); sslMode != "" {
		// dbConnection += fmt.Sprintf("sslmode=%s ", sslMode)
		dbURLOptions += fmt.Sprintf("?sslmode=%s", sslMode)
	}

	for retries := conf.C.Int("database.retries"); retries > 0 && !conf.Stop.Bool(); retries-- {
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
				"database.host", conf.C.String("database.host"),
				"database.username", conf.C.String("database.username"),
				"database.password", "****",
				"database.port", conf.C.Int("database.port"),
			)
			time.Sleep(conf.C.Duration("database.sleep_between_retries"))
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
	if conf.Stop.Bool() {
		return nil, fmt.Errorf("Database connection aborted")
	}

	// Still not connected?
	if err != nil {
		return nil, fmt.Errorf("Could not connect to database: %s", err)
	}

	connConfig, err := pgx.ParseConfig(fullDbURL)
	if err != nil {
		return nil, fmt.Errorf("Could not parse pgx config: %s", err)
	}
	if conf.C.Bool("database.log_queries") {
		connConfig.Logger = &queryLogger{logger: logger}
	}

	// Make the connection using the sqlx connector now that we know the database exists
	db, err := sqlx.Connect("pgx", stdlib.RegisterConnConfig(connConfig))
	if err != nil {
		return nil, fmt.Errorf("Could not connect to database: %s", err)
	}

	// Ping the database
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Could not ping database %s", err)
	}

	db.SetMaxOpenConns(conf.C.Int("database.max_connections"))

	logger.Debugw("Connected to database server",
		"database.host", conf.C.String("database.host"),
		"database.username", conf.C.String("database.username"),
		"database.port", conf.C.Int("database.port"),
		"database.database", conf.C.String("database.database"),
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
	assets, err := embed.AssetDir("postgres_migrations")
	if err != nil {
		return nil, fmt.Errorf("Could not get migrations assets")
	}
	assetSource := bindata.Resource(assets,
		func(name string) ([]byte, error) {
			return embed.Asset("postgres_migrations/" + name)
		})
	sourceDriver, err := bindata.WithInstance(assetSource)
	if err != nil {
		return nil, fmt.Errorf("Could not create migrations source driver: %v", err)
	}
	databaseDriver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("Could not create migrations database driver: %v", err)
	}
	migrateInstance, err := migrate.NewWithInstance("go-bindata", sourceDriver, "pgx", databaseDriver)
	if err != nil {
		logger.Errorw("Database migration error",
			"error", err,
		)
		return nil, fmt.Errorf("Migrate Error:%s", err)
	}

	// Do we wipe the database
	if conf.C.Bool("database.wipe_confirm") {
		err = migrateInstance.Down()
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
	err = migrateInstance.Up()
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

func wrapError(err error) error {

	switch e := err.(type) {
	case *pgconn.PgError:
		switch e.Code {
		case "23502":
			return fmt.Errorf("missing data: %s", e.Message)
		case "23503":
			// foreign key violation
			return fmt.Errorf("foreign key error: %s", e.Detail)
		case "23505":
			// unique constraint violation
			return fmt.Errorf("duplicate data: %s", e.Detail)
		case "23514":
			// check constraint violation
			return fmt.Errorf("invalid data: %s", e.Message)
		default:
			return &store.InternalError{
				Err: err,
			}
		}
	}
	return err
}
