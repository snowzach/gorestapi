package embed

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4/source"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

func MigrationSource() (source.Driver, error) {

	// wrap assets into Resource
	assets, err := AssetDir("postgres_migrations")
	if err != nil {
		return nil, fmt.Errorf("could not get migrations assets: %w", err)
	}
	assetSource := bindata.Resource(assets,
		func(name string) ([]byte, error) {
			return Asset("postgres_migrations/" + name)
		})
	sourceDriver, err := bindata.WithInstance(assetSource)
	if err != nil {
		return nil, fmt.Errorf("could not create migrations source driver: %w", err)
	}

	return sourceDriver, nil

}
