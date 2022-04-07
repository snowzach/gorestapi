package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"

	"github.com/snowzach/gorestapi/store/driver/postgres"
)

type Config struct {
	postgres.Config `conf:",squash"`
}

type Client struct {
	db    *sqlx.DB
	newID func() string
}

// New returns a new database client
func New(cfg *Config) (*Client, error) {

	db, err := postgres.New(&cfg.Config)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return &Client{
		db: db,
		newID: func() string {
			return xid.New().String()
		},
	}, nil

}
