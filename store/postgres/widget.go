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
		return wrapError(err)
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
		return nil, wrapError(err)
	}

	widget.SyncDB()

	return widget, nil

}

// WidgetDeleteByID an widget
func (c *Client) WidgetDeleteByID(ctx context.Context, id string) error {

	_, err := c.db.ExecContext(ctx, `DELETE FROM widget WHERE widget.id = $1`, id)
	if err != nil {
		return wrapError(err)
	}
	return nil

}

// WidgetsFind fetches a widgets with filter and pagination
func (c *Client) WidgetsFind(ctx context.Context, qp *queryp.QueryParameters) ([]*gorestapi.Widget, int64, error) {

	var queryClause strings.Builder
	var queryParams = []interface{}{}

	filterFields := queryp.FilterFieldTypes{
		"widget.id":          queryp.FilterTypeSimple,
		"widget.name":        queryp.FilterTypeString,
		"widget.description": queryp.FilterTypeString,
	}

	sortFields := queryp.SortFields{
		"widget.id":      "",
		"widget.created": "",
		"widget.updated": "",
		"widget.name":    "",
	}
	// Default sort
	if len(qp.Sort) == 0 {
		qp.Sort = queryp.Sort{queryp.SortTerm{Field: "widget.id", Desc: false}}
	}

	if len(qp.Filter) > 0 {
		queryClause.WriteString(" WHERE ")
	}

	if err := qppg.FilterQuery(filterFields, qp.Filter, &queryClause, &queryParams); err != nil {
		return nil, 0, &store.Error{Type: store.ErrorTypeQuery, Err: err}
	}

	var count int64
	if err := c.db.GetContext(ctx, &count, `SELECT COUNT(*) AS count`+WidgetFrom+queryClause.String(), queryParams...); err != nil {
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

	var widgets = make([]*gorestapi.Widget, 0)
	err := c.db.SelectContext(ctx, &widgets, WidgetSelect+WidgetFrom+queryClause.String(), queryParams...)
	if err != nil {
		return widgets, 0, wrapError(err)
	}

	for _, widget := range widgets {
		widget.SyncDB()
	}

	return widgets, count, nil
}
