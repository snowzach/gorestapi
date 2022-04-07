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
	WidgetSchema = ``
	WidgetTable  = `widget`
	WidgetJoins  = `
	LEFT JOIN thing ON widget.thing_id = thing.id
	`
	WidgetFields = `
	COALESCE(widget.id, '') as "widget.id",
	COALESCE(widget.created, '0001-01-01 00:00:00 UTC') as "widget.created",
	COALESCE(widget.updated, '0001-01-01 00:00:00 UTC') as "widget.updated",
	COALESCE(widget.name, '') as "widget.name",
	COALESCE(widget.description, '') as "widget.description",
	COALESCE(widget.thing_id, '') as "widget.thing_id"
	`
)

var (
	WidgetSelect = "SELECT " + strings.Join([]string{
		"widget.*",
		ThingFields,
	}, ",")
)

// WidgetSave saves a record
func (c *Client) WidgetSave(ctx context.Context, record *gorestapi.Widget) error {

	// Generate an ID if needed
	if record.ID == "" {
		record.ID = c.newID()
	}

	fields, values, updates, args := postgres.ComposeUpsert([]postgres.Field{
		{Name: "id", Insert: "$#", Update: "", Arg: record.ID},
		{Name: "created", Insert: "NOW()", Update: ""},
		{Name: "updated", Insert: "", Update: "NOW()"},
		{Name: "name", Insert: "$#", Update: "$#", Arg: record.Name},
		{Name: "description", Insert: "$#", Update: "$#", Arg: record.Description},
		{Name: "thing_id", Insert: "$#", Update: "$#", Arg: record.ThingID},
	})

	err := c.db.GetContext(ctx, record, `
	WITH `+WidgetTable+` AS (
        INSERT INTO `+ThingSchema+WidgetTable+` (`+fields+`)
        VALUES(`+values+`) ON CONFLICT (id) DO UPDATE
        SET `+updates+` RETURNING *
	) `+WidgetSelect+" FROM "+WidgetTable+WidgetJoins, args...)
	if err != nil {
		return postgres.WrapError(err)
	}

	record.SyncDB() // Clean empty structs

	return nil

}

// WidgetGetByID returns the record by id
func (c *Client) WidgetGetByID(ctx context.Context, id string) (*gorestapi.Widget, error) {

	record := new(gorestapi.Widget)
	err := c.db.GetContext(ctx, record, WidgetSelect+` FROM `+WidgetSchema+WidgetTable+WidgetJoins+` WHERE `+WidgetTable+`.id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, postgres.WrapError(err)
	}

	record.SyncDB() // Clean empty structs

	return record, nil

}

// WidgetDeleteByID deletes a record by id
func (c *Client) WidgetDeleteByID(ctx context.Context, id string) error {

	_, err := c.db.ExecContext(ctx, `DELETE FROM `+WidgetSchema+WidgetTable+` WHERE `+WidgetTable+`.id = $1`, id)
	if err != nil {
		return postgres.WrapError(err)
	}
	return nil

}

// WidgetsFind fetches records with filter and pagination
func (c *Client) WidgetsFind(ctx context.Context, qp *queryp.QueryParameters) ([]*gorestapi.Widget, int64, error) {

	var queryClause strings.Builder
	var queryParams = []interface{}{}

	filterFields := queryp.FilterFieldTypes{
		"widget.id":          queryp.FilterTypeSimple,
		"widget.name":        queryp.FilterTypeString,
		"widget.description": queryp.FilterTypeString,
		"widget.thing_id":    queryp.FilterTypeSimple,
		"thing.name":         queryp.FilterTypeString,
	}

	sortFields := queryp.SortFields{
		"widget.id":      "",
		"widget.created": "",
		"widget.updated": "",
		"widget.name":    "",
		"thing.name":     "",
	}
	// Default sort
	if len(qp.Sort) == 0 {
		qp.Sort.Append("widget.id", false)
	}

	if len(qp.Filter) > 0 {
		queryClause.WriteString(" WHERE ")
	}

	if err := qppg.FilterQuery(filterFields, qp.Filter, &queryClause, &queryParams); err != nil {
		return nil, 0, &store.Error{Type: store.ErrorTypeQuery, Err: err}
	}

	var count int64
	if err := c.db.GetContext(ctx, &count, `SELECT COUNT(*) AS count FROM `+WidgetSchema+WidgetTable+WidgetJoins+queryClause.String(), queryParams...); err != nil {
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

	var records = make([]*gorestapi.Widget, 0)
	err := c.db.SelectContext(ctx, &records, WidgetSelect+` FROM `+WidgetSchema+WidgetTable+WidgetJoins+queryClause.String(), queryParams...)
	if err != nil {
		return records, 0, postgres.WrapError(err)
	}

	for _, record := range records {
		record.SyncDB()
	}

	return records, count, nil
}
