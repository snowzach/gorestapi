package postgres

import (
	"context"
	"database/sql/driver"

	"github.com/rs/xid"
	"github.com/snowzach/queryp"

	"github.com/snowzach/golib/store/driver/postgres"
	"github.com/snowzach/gorestapi/gorestapi"
)

var (
	WidgetTable = postgres.Generate(postgres.Table[gorestapi.Widget]{
		Table: `"widget"`,
		Fields: []*postgres.Field[gorestapi.Widget]{
			{Name: "id", ID: true, Insert: "$#", Value: func(rec *gorestapi.Widget) (driver.Value, error) { return rec.ID, nil }},
			{Name: "created", Insert: "NOW()", NullVal: "0001-01-01 00:00:00 UTC"},
			{Name: "updated", Insert: "NOW()", Update: "NOW()", NullVal: "0001-01-01 00:00:00 UTC"},
			{Name: "name", Insert: "$#", Update: "$#", Value: func(rec *gorestapi.Widget) (driver.Value, error) { return rec.Name, nil }},
			{Name: "description", Insert: "$#", Update: "$#", Value: func(rec *gorestapi.Widget) (driver.Value, error) { return rec.Description, nil }},
			{Name: "thing_id", Insert: "$#", Update: "$#", Value: func(rec *gorestapi.Widget) (driver.Value, error) { return rec.ThingID, nil }},
		},
		Joins: `
		LEFT JOIN thing ON widget.thing_id = thing.id
		`,
		Selector: postgres.Selector[gorestapi.Widget]{
			FilterFieldTypes: queryp.FilterFieldTypes{
				"widget.id":          queryp.FilterTypeSimple,
				"widget.created":     queryp.FilterTypeTime,
				"widget.updated":     queryp.FilterTypeTime,
				"widget.name":        queryp.FilterTypeString,
				"widget.description": queryp.FilterTypeString,
				"thing.name":         queryp.FilterTypeString,
				"thing.description":  queryp.FilterTypeString,
			},
			SortFields: queryp.SortFields{
				"widget.id":          "",
				"widget.created":     "",
				"widget.updated":     "",
				"widget.name":        "",
				"widget.description": "",
				"thing.name":         "",
				"thing.description":  "",
			},
			DefaultSort: queryp.Sort{
				&queryp.SortTerm{Field: "thing.name", Desc: false},
			},
			PostProcessRecord: func(rec *gorestapi.Widget) error {
				if rec.ThingID == nil || *rec.ThingID == "" {
					rec.Thing = nil
				}
				return nil
			},
		},
		PostProcessRecord: func(rec *gorestapi.Widget) error {
			if rec.ThingID == nil || *rec.ThingID == "" {
				rec.Thing = nil
			}
			return nil
		},
		SelectAdditionalFields: ThingTable.GenerateAdditionalFields(true),
	})
)

// WidgetSave saves the record
func (c *Client) WidgetSave(ctx context.Context, record *gorestapi.Widget) error {
	if record.ID == "" {
		record.ID = xid.New().String()
	}
	return WidgetTable.Upsert(ctx, c.db, record)
}

// WidgetGetByID returns the the record by id
func (c *Client) WidgetGetByID(ctx context.Context, id string) (*gorestapi.Widget, error) {
	return WidgetTable.GetByID(ctx, c.db, id)
}

// WidgetDeleteByID deletes a record by id
func (c *Client) WidgetDeleteByID(ctx context.Context, id string) error {
	return WidgetTable.DeleteByID(ctx, c.db, id)
}

// WidgetsFind fetches records with filter and pagination
func (c *Client) WidgetsFind(ctx context.Context, qp *queryp.QueryParameters) ([]*gorestapi.Widget, *int64, error) {
	return WidgetTable.Selector.Select(ctx, c.db, qp)
}
