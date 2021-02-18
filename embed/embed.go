package embed

import (
	"embed"
	"io/fs"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/johejo/golang-migrate-extra/source/iofs"
)

//go:embed postgres_migrations
var postgresMigrations embed.FS

func MigrationSource() (source.Driver, error) {
	return iofs.New(postgresMigrations, "postgres_migrations")
}

//go:embed public_html
var publicHTML embed.FS

func PublicHTMLFS() fs.FS {
	publicHTMLfs, _ := fs.Sub(publicHTML, "public_html")
	return publicHTMLfs
}
