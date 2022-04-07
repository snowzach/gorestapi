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
	"github.com/snowzach/gorestapi/store/driver/postgres"
)

const (
	ThingSchema = ``
	ThingTable  = `thing`
	ThingJoins  = ``
	ThingFields = `COALESCE(thing.id, '') as "thing.id",
	COALESCE(thing.created, '0001-01-01 00:00:00 UTC') as "thing.created",
	COALESCE(thing.updated, '0001-01-01 00:00:00 UTC') as "thing.updated",
	COALESCE(thing.name, '') as "thing.name",
	COALESCE(thing.description, '') as "thing.description"
	`
)

var (
	ThingSelect = "SELECT " + strings.Join([]string{
		"thing.*",
	}, ",")
)

// ThingSave saves the record
func (c *Client) ThingSave(ctx context.Context, record *gorestapi.Thing) error {

	if record.ID == "" {
		record.ID = c.newID()
	}

	fields, values, updates, args := postgres.ComposeUpsert([]postgres.Field{
		{Name: "id", Insert: "$#", Update: "", Arg: record.ID},
		{Name: "created", Insert: "NOW()", Update: ""},
		{Name: "updated", Insert: "", Update: "NOW()"},
		{Name: "name", Insert: "$#", Update: "$#", Arg: record.Name},
		{Name: "description", Insert: "$#", Update: "$#", Arg: record.Description},
	})

	err := c.db.GetContext(ctx, record, `
	WITH `+ThingTable+` AS (
        INSERT INTO `+ThingSchema+ThingTable+` (`+fields+`)
        VALUES(`+values+`) ON CONFLICT (id) DO UPDATE
        SET `+updates+` RETURNING *
	) `+ThingSelect+" FROM "+ThingTable+ThingJoins, args...)
	if err != nil {
		return postgres.WrapError(err)
	}
	return nil

}

// ThingGetByID returns the the record by id
func (c *Client) ThingGetByID(ctx context.Context, id string) (*gorestapi.Thing, error) {

	thing := new(gorestapi.Thing)
	err := c.db.GetContext(ctx, thing, ThingSelect+` FROM `+ThingSchema+ThingTable+ThingJoins+` WHERE `+ThingTable+`.id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, postgres.WrapError(err)
	}
	return thing, nil

}

// ThingDeleteByID deletes a record by id
func (c *Client) ThingDeleteByID(ctx context.Context, id string) error {

	_, err := c.db.ExecContext(ctx, `DELETE FROM `+ThingSchema+ThingTable+` WHERE `+ThingTable+`.id = $1`, id)
	if err != nil {
		return postgres.WrapError(err)
	}
	return nil

}

// ThingsFind fetches records with filter and pagination
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
		qp.Sort.Append("thing.id", false)
	}

	if len(qp.Filter) > 0 {
		queryClause.WriteString(" WHERE ")
	}

	if err := qppg.FilterQuery(filterFields, qp.Filter, &queryClause, &queryParams); err != nil {
		return nil, 0, &store.Error{Type: store.ErrorTypeQuery, Err: err}
	}
	var count int64
	if err := c.db.GetContext(ctx, &count, `SELECT COUNT(*) AS count FROM `+ThingSchema+ThingTable+ThingJoins+queryClause.String(), queryParams...); err != nil {
		return nil, 0, postgres.WrapError(err)
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

	var records = make([]*gorestapi.Thing, 0)
	err := c.db.SelectContext(ctx, &records, ThingSelect+` FROM `+ThingSchema+ThingTable+ThingJoins+queryClause.String(), queryParams...)
	if err != nil {
		return records, 0, postgres.WrapError(err)
	}

	return records, count, nil
}
