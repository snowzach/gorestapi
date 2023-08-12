package postgres

import (
	"context"
	"database/sql/driver"

	"github.com/rs/xid"
	"github.com/snowzach/golib/store/driver/postgres"
	"github.com/snowzach/queryp"

	"github.com/snowzach/gorestapi/gorestapi"
)

var (
	ThingTable = postgres.Generate(postgres.Table[gorestapi.Thing]{
		Table: `"thing"`,
		Fields: []*postgres.Field[gorestapi.Thing]{
			{Name: "id", ID: true, Insert: "$#", Value: func(rec *gorestapi.Thing) (driver.Value, error) { return rec.ID, nil }},
			{Name: "created", Insert: "NOW()", NullVal: "0001-01-01 00:00:00 UTC"},
			{Name: "updated", Insert: "NOW()", Update: "NOW()", NullVal: "0001-01-01 00:00:00 UTC"},
			{Name: "name", Insert: "$#", Update: "$#", Value: func(rec *gorestapi.Thing) (driver.Value, error) { return rec.Name, nil }},
			{Name: "description", Insert: "$#", Update: "$#", Value: func(rec *gorestapi.Thing) (driver.Value, error) { return rec.Description, nil }},
		},
		Selector: postgres.Selector[gorestapi.Thing]{
			FilterFieldTypes: queryp.FilterFieldTypes{
				"thing.id":          queryp.FilterTypeSimple,
				"thing.created":     queryp.FilterTypeTime,
				"thing.updated":     queryp.FilterTypeTime,
				"thing.name":        queryp.FilterTypeString,
				"thing.description": queryp.FilterTypeString,
			},
			SortFields: queryp.SortFields{
				"thing.id":          "",
				"thing.created":     "",
				"thing.updated":     "",
				"thing.name":        "",
				"thing.description": "",
			},
			DefaultSort: queryp.Sort{
				&queryp.SortTerm{Field: "thing.name", Desc: false},
			},
		},
	})
)

// ThingSave saves the record
func (c *Client) ThingSave(ctx context.Context, record *gorestapi.Thing) error {
	if record.ID == "" {
		record.ID = xid.New().String()
	}
	return ThingTable.Upsert(ctx, c.db, record)
}

// ThingGetByID returns the the record by id
func (c *Client) ThingGetByID(ctx context.Context, id string) (*gorestapi.Thing, error) {
	return ThingTable.GetByID(ctx, c.db, id)
}

// ThingDeleteByID deletes a record by id
func (c *Client) ThingDeleteByID(ctx context.Context, id string) error {
	return ThingTable.DeleteByID(ctx, c.db, id)
}

// ThingsFind fetches records with filter and pagination
func (c *Client) ThingsFind(ctx context.Context, qp *queryp.QueryParameters) ([]*gorestapi.Thing, *int64, error) {
	return ThingTable.Selector.Select(ctx, c.db, qp)
}
