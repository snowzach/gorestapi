package postgres

import (
	"context"
	"database/sql"

	"github.com/snowzach/gorestapi/gorestapi"
	"github.com/snowzach/gorestapi/store"
)

// ThingGetByID returns the the thing by ID
func (c *Client) ThingGetByID(ctx context.Context, id string) (*gorestapi.Thing, error) {

	b := new(gorestapi.Thing)
	err := c.db.GetContext(ctx, b, `SELECT * FROM thing WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return b, nil

}

// ThingSave saves the thing
func (c *Client) ThingSave(ctx context.Context, i *gorestapi.Thing) (string, error) {

	// Generate an ID if needed
	if i.ID == "" {
		i.ID = c.newID()
	}

	_, err := c.db.ExecContext(ctx, `
		INSERT INTO thing (id, name)
		VALUES($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET name = $2
	`, i.ID, i.Name)
	if err != nil {
		return i.ID, err
	}
	return i.ID, nil

}

// ThingDeleteByID an thing
func (c *Client) ThingDeleteByID(ctx context.Context, id string) error {

	_, err := c.db.ExecContext(ctx, `DELETE FROM thing WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil

}

// ThingFind gets things
func (c *Client) ThingFind(ctx context.Context) ([]*gorestapi.Thing, error) {

	var bs = make([]*gorestapi.Thing, 0)
	err := c.db.SelectContext(ctx, &bs, `SELECT * FROM thing`)
	if err == sql.ErrNoRows {
		// No Error
	} else if err != nil {
		return bs, err
	}
	return bs, nil

}
