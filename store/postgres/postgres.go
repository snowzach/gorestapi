package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	// Import Database Migrate Postgres suppose
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/knadh/koanf"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/snowzach/gorestapi/store"
)

// Client is the database client
type Client struct {
	logger *zap.SugaredLogger
	db     *sqlx.DB
	newID  func() string
}

// New returns a new database client
func New(cfg *koanf.Koanf, migrationSource source.Driver) (*Client, error) {

	logger := zap.S().With("package", "store.postgres")

	var err error

	// Credentials
	var dbCreds string
	if username := cfg.String("database.username"); username != "" {
		dbCreds += fmt.Sprintf("user=%s ", username)
	}
	if password := cfg.String("database.password"); password != "" {
		dbCreds += fmt.Sprintf("password=%s ", password)
	}

	// Host + Port
	var connStr strings.Builder // Regular credentials
	if hostname := cfg.String("database.host"); hostname != "" {
		connStr.WriteString(fmt.Sprintf("host=%s ", hostname))
	} else {
		return nil, fmt.Errorf("No hostname specified")
	}
	if port := cfg.String("database.port"); port != "" {
		connStr.WriteString(fmt.Sprintf("port=%s ", port))
	}

	// SSL Mode
	connStr.WriteString(fmt.Sprintf("sslmode=%s ", cfg.String("database.sslmode")))
	if sslCert := cfg.String("database.sslcert"); sslCert != "" {
		connStr.WriteString(fmt.Sprintf("sslcert=%s ", sslCert))
	}
	if sslKey := cfg.String("database.sslkey"); sslKey != "" {
		connStr.WriteString(fmt.Sprintf("sslkey=%s ", sslKey))
	}
	if sslRootCert := cfg.String("database.sslrootcert"); sslRootCert != "" {
		connStr.WriteString(fmt.Sprintf("sslrootcert=%s ", sslRootCert))
	}

	// Search Path
	if searchPath := cfg.String("database.search_path"); searchPath != "" {
		connStr.WriteString(fmt.Sprintf("search_path=%s ", searchPath))
	}

	// Database Name
	dbName := cfg.String("database.database")

	var db *sqlx.DB

	// Auto-create the database if requested/needed/possible
	if cfg.Bool("database.auto_create") {

		// Check to see if we have an admin password specified (that we will use to create the database if it does not exist)
		var dbCreateCreds string
		if pgPassword := os.Getenv("POSTGRES_PASSWORD"); pgPassword != "" {
			pgUser := os.Getenv("POSTGRES_USER")
			if pgUser == "" {
				pgUser = "postgres"
			}
			dbCreateCreds = fmt.Sprintf("user=%s password=%s ", pgUser, pgPassword)
		} else {
			dbCreateCreds = dbCreds // otherwise use the default credentials
		}

		// Connect using create credentials
		createConnConfig, err := pgx.ParseConfig(dbCreateCreds + connStr.String())
		if err != nil {
			return nil, fmt.Errorf("could not parse pgx create config: %w", err)
		}

		for retries := cfg.Int("database.retries"); retries > 0; retries-- {
			// Attempt connecting to the database
			db, err = sqlx.Connect("pgx", stdlib.RegisterConnConfig(createConnConfig))
			if err == nil {
				// Ping the database
				if err = db.Ping(); err != nil {
					return nil, fmt.Errorf("could not ping database %w", err)
				}
				break // connected

			} else if strings.Contains(err.Error(), "connection refused") {
				logger.Warn("failed to connect to database, sleeping and retry")
				time.Sleep(cfg.Duration("database.sleep_between_retries"))
				continue
			}

			// Some other error
			return nil, err
		}
		if err != nil {
			return nil, fmt.Errorf("retries exausted, last error: %v", err)
		}

		logger.Infow("Checking for database", "database", dbName)
		var one int
		if err := db.Get(&one, `SELECT 1 from pg_database WHERE datname=$1`, dbName); err == sql.ErrNoRows {

			logger.Infow("Creating database", "database", dbName)
			_, err = db.Exec(`CREATE DATABASE ` + dbName)
			if err != nil {
				return nil, fmt.Errorf("could not create database: %w", err)
			}

		} else if err != nil {
			// Some other error besides does not exist
			return nil, fmt.Errorf("could not check for database: %w", err)
		}

		_ = db.Close()
		db = nil

	}

	// Connect to database
	connStr.WriteString(fmt.Sprintf("dbname=%s ", dbName))
	connConfig, err := pgx.ParseConfig(dbCreds + connStr.String())
	if err != nil {
		return nil, fmt.Errorf("could not parse pgx config: %s", err)
	}
	if cfg.Bool("database.log_queries") {
		connConfig.Logger = &queryLogger{logger: logger}
	}

	for retries := cfg.Int("database.retries"); retries > 0; retries-- {

		// Attempt connecting to the database
		db, err = sqlx.Connect("pgx", stdlib.RegisterConnConfig(connConfig))
		if err == nil {
			// Ping the database
			if err = db.Ping(); err != nil {
				return nil, fmt.Errorf("could not ping database %w", err)
			}
			break // connected
		} else if strings.Contains(err.Error(), "connection refused") {
			logger.Warn("failed to connect to database, sleeping and retry")
			time.Sleep(cfg.Duration("database.sleep_between_retries"))
			continue
		}

		// Some other error
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("retries exausted, last error: %v", err)
	}

	db.SetMaxOpenConns(cfg.Int("database.max_connections"))

	logger.Debugw("Connected to database server",
		"database.host", cfg.String("database.host"),
		"database.username", cfg.String("database.username"),
		"database.port", cfg.Int("database.port"),
		"database.database", cfg.String("database.database"),
	)

	c := &Client{
		logger: logger,
		db:     db,
		newID: func() string {
			return xid.New().String()
		},
	}

	// If we have migrations we can apply, use them
	if migrationSource != nil {

		databaseDriver, err := postgres.WithInstance(db.DB, &postgres.Config{})
		if err != nil {
			return nil, fmt.Errorf("could not create migrations database instance: %w", err)
		}
		migrateInstance, err := migrate.NewWithInstance("source", migrationSource, "pgx", databaseDriver)
		if err != nil {
			return nil, fmt.Errorf("could not create migrations instance error: %w", err)
		}

		// Do we wipe the database
		if cfg.Bool("database.wipe_confirm") {
			err = migrateInstance.Down()
			if err == migrate.ErrNoChange {
				// Okay
			} else if err != nil {
				return nil, fmt.Errorf("migrate down error: %w", err)
			} else {
				logger.Warn("Database wipe complete")
			}
		}

		// Perform the migration up
		err = migrateInstance.Up()
		if err == migrate.ErrNoChange {
			logger.Info("Database schema current")
		} else if err != nil {
			return nil, fmt.Errorf("migrate up error: %w", err)
		} else {
			logger.Info("Database migration completed")
		}
	}

	return c, nil

}

type queryLogger struct {
	logger *zap.SugaredLogger
}

func (ql *queryLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	ql.logger.Debugw(msg, "level", level, zap.Any("data", data))
}

// Lookup of postgres error codes to basic errors we can return to a user
var pgErrorCodeToStoreErrorType = map[string]store.ErrorType{
	"23502": store.ErrorTypeIncomplete,
	"23503": store.ErrorTypeForeignKey,
	"23505": store.ErrorTypeDuplicate,
	"23514": store.ErrorTypeInvalid,
}

func wrapError(err error) error {
	switch e := err.(type) {
	case *pgconn.PgError:
		if et, found := pgErrorCodeToStoreErrorType[e.Code]; found {
			return &store.Error{
				Type: et,
				Err:  err,
			}
		}
	}
	return err
}

type field struct {
	name   string
	insert string
	update string
	arg    interface{}
}

// Builds the values needed to compose an upsert statement
func composeUpsert(fields []field) (string, string, string, []interface{}) {

	names := make([]string, 0)
	inserts := make([]string, 0)
	updates := make([]string, 0)
	args := make([]interface{}, 0)

	for _, field := range fields {
		index := "$#"
		if field.arg != nil {
			args = append(args, field.arg)
			index = "$" + strconv.Itoa(len(args))
		}
		if field.insert != "" {
			names = append(names, field.name)
			inserts = append(inserts, strings.ReplaceAll(field.insert, "$#", index))
		}
		if field.update != "" {
			updates = append(updates, field.name+" = "+strings.ReplaceAll(field.update, "$#", index))
		}
	}

	return strings.Join(names, ","), strings.Join(inserts, ","), strings.Join(updates, ","), args

}
