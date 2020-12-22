package postgres

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/snowzach/queryp"
	"github.com/snowzach/queryp/qppg"

	"github.com/snowzach/gorestapi/gorestapi"
	"github.com/snowzach/gorestapi/store"
)

var (
	ThingSelect = "SELECT " + strings.Join([]string{
		"thing.*",
	}, ",")
)

const (
	ThingFrom = ` FROM thing`

	ThingFields = `COALESCE(thing.id, '') as "thing.id",
	COALESCE(thing.created, '0001-01-01 00:00:00 UTC') as "thing.created",
	COALESCE(thing.updated, '0001-01-01 00:00:00 UTC') as "thing.updated",
	COALESCE(thing.name, '') as "thing.name",
	COALESCE(thing.description, '') as "thing.description"
	`
)

// ThingSave saves the thing
func (c *Client) ThingSave(ctx context.Context, thing *gorestapi.Thing) error {

	// Generate an ID if needed
	if thing.ID == "" {
		thing.ID = c.newID()
	}

	err := c.db.GetContext(ctx, thing, `
	WITH thing AS (
		INSERT INTO thing (id, created, updated, name, description)
		VALUES($1, NOW(), NOW(), $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET 
		updated = NOW(),
		name = $2,
		description = $3
		RETURNING *
	) `+ThingSelect+ThingFrom, thing.ID, thing.Name, thing.Description)
	if err != nil {
		return wrapError(err)
	}
	return nil

}

// ThingGetByID returns the the thing by id
func (c *Client) ThingGetByID(ctx context.Context, id string) (*gorestapi.Thing, error) {

	thing := new(gorestapi.Thing)
	err := c.db.GetContext(ctx, thing, ThingSelect+ThingFrom+` WHERE thing.id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, wrapError(err)
	}
	return thing, nil

}

// ThingDeleteByID an thing
func (c *Client) ThingDeleteByID(ctx context.Context, id string) error {

	_, err := c.db.ExecContext(ctx, `DELETE FROM thing WHERE thing.id = $1`, id)
	if err != nil {
		return wrapError(err)
	}
	return nil

}

// ThingsFind fetches a things with filter and pagination
func (c *Client) ThingsFind(ctx context.Context, qp *queryp.QueryParameters) ([]*gorestapi.Thing, int64, error) {

	var queryClause strings.Builder
	var queryParams = []interface{}{}

	filterFields := queryp.FilterFieldTypes{
		"thing.id":          queryp.FilterTypeSimple,
		"thing.name":        queryp.FilterTypeString,
		"thing.description": queryp.FilterTypeString,
	}

	sortFields := queryp.SortFields{
		"thing.id":      "",
		"thing.created": "",
		"thing.updated": "",
		"thing.name":    "",
	}
	// Default sort
	if len(qp.Sort) == 0 {
		qp.Sort = queryp.Sort{queryp.SortTerm{Field: "thing.id", Desc: false}}
	}

	if len(qp.Filter) > 0 {
		queryClause.WriteString(" WHERE ")
	}

	if err := qppg.FilterQuery(filterFields, qp.Filter, &queryClause, &queryParams); err != nil {
		return nil, 0, &store.Error{Type: store.ErrorTypeQuery, Err: err}
	}
	var count int64
	if err := c.db.GetContext(ctx, &count, `SELECT COUNT(*) AS count`+ThingFrom+queryClause.String(), queryParams...); err != nil {
		return nil, 0, wrapError(err)
	}
	if err := qppg.SortQuery(sortFields, qp.Sort, &queryClause, &queryParams); err != nil {
		return nil, 0, &store.Error{Type: store.ErrorTypeQuery, Err: err}
	}
	if qp.Limit > 0 {
		queryClause.WriteString(" LIMIT " + strconv.FormatInt(qp.Limit, 10))
	}
	if qp.Offset > 0 {
		queryClause.WriteString(" OFFSET " + strconv.FormatInt(qp.Offset, 10))
	}

	var things = make([]*gorestapi.Thing, 0)
	err := c.db.SelectContext(ctx, &things, ThingSelect+ThingFrom+queryClause.String(), queryParams...)
	if err != nil {
		return things, 0, wrapError(err)
	}

	return things, count, nil
}
