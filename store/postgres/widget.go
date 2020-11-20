package postgres

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/snowzach/gorestapi/gorestapi"
	"github.com/snowzach/gorestapi/store"
)

var (
	WidgetSelect = "SELECT " + strings.Join([]string{
		"widget.*",
		ThingFields,
	}, ",")
)

const (
	WidgetFrom = ` FROM widget
	LEFT JOIN thing ON widget.thing_id = thing.id
	`

	WidgetFields = `COALESCE(widget.id, '') as "widget.id",
	COALESCE(widget.created, '0001-01-01 00:00:00 UTC') as "widget.created",
	COALESCE(widget.updated, '0001-01-01 00:00:00 UTC') as "widget.updated",
	COALESCE(widget.name, '') as "widget.name",
	COALESCE(widget.description, '') as "widget.description"
	widget.thing_id
	`
)

// WidgetSave saves the widget
func (c *Client) WidgetSave(ctx context.Context, widget *gorestapi.Widget) error {

	// Generate an ID if needed
	if widget.ID == "" {
		widget.ID = c.newID()
	}

	err := c.db.GetContext(ctx, widget, `
	WITH widget AS (
		INSERT INTO widget (id, created, updated, name, description, thing_id)
		VALUES($1, NOW(), NOW(), $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET 
		updated = NOW(),
		name = $2,
		description = $3,
		thing_id = $4
		RETURNING *
	) `+WidgetSelect+WidgetFrom, widget.ID, widget.Name, widget.Description, widget.ThingID)
	if err != nil {
		return err
	}

	widget.SyncDB()

	return nil

}

// WidgetGetByID returns the the widget by id
func (c *Client) WidgetGetByID(ctx context.Context, id string) (*gorestapi.Widget, error) {

	widget := new(gorestapi.Widget)
	err := c.db.GetContext(ctx, widget, WidgetSelect+WidgetFrom+`WHERE widget.id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	widget.SyncDB()

	return widget, nil

}

// WidgetDeleteByID an widget
func (c *Client) WidgetDeleteByID(ctx context.Context, id string) error {

	_, err := c.db.ExecContext(ctx, `DELETE FROM widget WHERE widget.id = $1`, id)
	if err != nil {
		return err
	}
	return nil

}

// WidgetsFind fetches a widgets with filter and pagination
func (c *Client) WidgetsFind(ctx context.Context, fqp *store.FindQueryParameters) ([]*gorestapi.Widget, int64, error) {

	var queryClause strings.Builder
	var queryParams = []interface{}{}

	filterFields := store.FilterFieldTypes{
		"widget.id":   store.FilterTypeEquals,
		"widget.name": store.FilterTypeILike,
	}

	sortFields := store.SortFields{
		"widget.id",
		"widget.created",
		"widget.updated",
		"widget.name",
	}
	// Default sort
	if len(fqp.Sort) == 0 {
		fqp.Sort = store.SortValues{&store.Sort{Field: "widget.id", Desc: false}}
	}

	if err := c.FilterQuery(filterFields, fqp.PreFilter, fqp.PreFilterInclusive, &queryClause, &queryParams); err != nil {
		return nil, 0, err
	}
	if err := c.FilterQuery(filterFields, fqp.Filter, fqp.FilterInclusive, &queryClause, &queryParams); err != nil {
		return nil, 0, err
	}
	var count int64
	if err := c.db.GetContext(ctx, &count, `SELECT COUNT(*) AS count`+WidgetFrom+` WHERE 1=1`+queryClause.String(), queryParams...); err != nil {
		return nil, 0, err
	}
	if err := c.SortQuery(sortFields, fqp.Sort, &queryClause, &queryParams); err != nil {
		return nil, 0, err
	}
	if fqp.Limit > 0 {
		queryClause.WriteString(" LIMIT " + strconv.Itoa(fqp.Limit))
	}
	if fqp.Offset > 0 {
		queryClause.WriteString(" OFFSET " + strconv.Itoa(fqp.Offset))
	}

	var widgets = make([]*gorestapi.Widget, 0)
	err := c.db.SelectContext(ctx, &widgets, WidgetSelect+WidgetFrom+` WHERE 1=1`+queryClause.String(), queryParams...)
	if err != nil {
		return widgets, 0, err
	}

	for _, widget := range widgets {
		widget.SyncDB()
	}

	return widgets, count, nil
}
