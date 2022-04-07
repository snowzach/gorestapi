package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	// Import Database Migrate Postgres suppose

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Username            string        `conf:"username" default:"postgres"`
	Password            string        `conf:"password" default:"password"`
	Host                string        `conf:"host" default:"postgres"`
	Port                string        `conf:"port" default:"5432"`
	Database            string        `conf:"database" default:"postgres"`
	AutoCreate          bool          `conf:"auto_create" default:"false"`
	SearchPath          string        `conf:"search_path" default:""`
	SSLMode             string        `conf:"sslmode" default:"false"`
	SSLCert             string        `conf:"sslcert" default:""`
	SSLKey              string        `conf:"sslkey" default:""`
	SSLRootCert         string        `conf:"sslrootcert" default:""`
	Retries             int           `conf:"retries" default:"5"`
	SleepBetweenRetries time.Duration `conf:"sleep_between_retries" default:"7s"`
	MaxConnections      int           `conf:"max_connections" default:"40"`
	WipeConfirm         bool          `conf:"wipe_confirm" default:"false"`

	Logger          Logger
	QueryLogger     Logger
	MigrationSource source.Driver
}

// New returns a new database client
func New(c *Config) (*sqlx.DB, error) {

	if c.Logger == nil {
		c.Logger = nopLogger{}
	}
	if c.QueryLogger == nil {
		c.QueryLogger = nopLogger{}
	}

	var err error

	// Credentials
	var dbCreds string
	if c.Username != "" {
		dbCreds += fmt.Sprintf("user=%s ", c.Username)
	}
	if c.Password != "" {
		dbCreds += fmt.Sprintf("password=%s ", c.Password)
	}

	// Host + Port
	var connStr strings.Builder // Regular credentials
	if c.Host != "" {
		connStr.WriteString(fmt.Sprintf("host=%s ", c.Host))
	} else {
		return nil, fmt.Errorf("No hostname specified")
	}
	if c.Port != "" {
		connStr.WriteString(fmt.Sprintf("port=%s ", c.Port))
	}

	// SSL Mode
	connStr.WriteString(fmt.Sprintf("sslmode=%s ", c.SSLMode))
	if c.SSLCert != "" {
		connStr.WriteString(fmt.Sprintf("sslcert=%s ", c.SSLCert))
	}
	if c.SSLKey != "" {
		connStr.WriteString(fmt.Sprintf("sslkey=%s ", c.SSLKey))
	}
	if c.SSLRootCert != "" {
		connStr.WriteString(fmt.Sprintf("sslrootcert=%s ", c.SSLRootCert))
	}

	// Search Path
	if c.SearchPath != "" {
		connStr.WriteString(fmt.Sprintf("search_path=%s ", c.SearchPath))
	}

	// Database Name
	dbName := c.Database

	var db *sqlx.DB

	// Auto-create the database if requested/needed/possible
	if c.AutoCreate {

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

		for retries := c.Retries; retries > 0; retries-- {
			// Attempt connecting to the database
			db, err = sqlx.Connect("pgx", stdlib.RegisterConnConfig(createConnConfig))
			if err == nil {
				// Ping the database
				if err = db.Ping(); err != nil {
					return nil, fmt.Errorf("could not ping database %w", err)
				}
				break // connected

			} else if strings.Contains(err.Error(), "connection refused") {
				c.Logger.Printf("Failed to connect to database, sleeping %d seconds before retry", c.SleepBetweenRetries)
				time.Sleep(c.SleepBetweenRetries)
				continue
			}

			// Some other error
			return nil, err
		}
		if err != nil {
			return nil, fmt.Errorf("retries exausted, last error: %v", err)
		}

		c.Logger.Printf("Checking for database: %s", dbName)
		var one int
		if err := db.Get(&one, `SELECT 1 from pg_database WHERE datname=$1`, dbName); err == sql.ErrNoRows {
			c.Logger.Printf("Creating database: %s", dbName)
			_, err = db.Exec(`CREATE DATABASE ` + dbName)
			if err != nil {
				return nil, fmt.Errorf("could not create database: %w", err)
			}

		} else if err != nil {
			// Some other error besides does not exist
			return nil, fmt.Errorf("could not check for database: %w", err)
		}

		// Close the temporary database.
		_ = db.Close()
		db = nil

	}

	// Connect to database
	connStr.WriteString(fmt.Sprintf("dbname=%s ", dbName))
	connConfig, err := pgx.ParseConfig(dbCreds + connStr.String())
	if err != nil {
		return nil, fmt.Errorf("could not parse pgx config: %s", err)
	}

	if c.QueryLogger != nil {
		connConfig.Logger = &queryLogger{Logger: c.QueryLogger}
	}

	for retries := c.Retries; retries > 0; retries-- {
		// Attempt connecting to the database
		db, err = sqlx.Connect("pgx", stdlib.RegisterConnConfig(connConfig))
		if err == nil {
			// Ping the database
			if err = db.Ping(); err != nil {
				return nil, fmt.Errorf("could not ping database %w", err)
			}
			break // connected
		} else if strings.Contains(err.Error(), "connection refused") {
			c.Logger.Printf("Failed to connect to database, sleeping %d seconds before retry", c.SleepBetweenRetries)
			time.Sleep(c.SleepBetweenRetries)
			continue
		}

		// Some other error
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("retries exhausted, last error: %v", err)
	}

	db.SetMaxOpenConns(c.MaxConnections)
	c.Logger.Printf("Connected to database %s:xxx@%s:%s/%s", c.Username, c.Host, c.Port, c.Database)

	// If we have migrations we can apply, use them
	if c.MigrationSource != nil {

		databaseDriver, err := postgres.WithInstance(db.DB, &postgres.Config{})
		if err != nil {
			return nil, fmt.Errorf("could not create migrations database instance: %w", err)
		}
		migrateInstance, err := migrate.NewWithInstance("source", c.MigrationSource, "pgx", databaseDriver)
		if err != nil {
			return nil, fmt.Errorf("could not create migrations instance error: %w", err)
		}

		// Do we wipe the database
		if c.WipeConfirm {
			err = migrateInstance.Down()
			if err == migrate.ErrNoChange {
				// Okay
			} else if err != nil {
				return nil, fmt.Errorf("migrate down error: %w", err)
			} else {
				c.Logger.Printf("Database wipe complete")
			}
		}

		// Perform the migration up
		err = migrateInstance.Up()
		if err == migrate.ErrNoChange {
			c.Logger.Printf("Database schema current")
		} else if err != nil {
			return nil, fmt.Errorf("migrate up error: %w", err)
		} else {
			c.Logger.Printf("Database migration completed")
		}
	}

	return db, nil

}
