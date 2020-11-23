package gorestapi

import (
	"context"

	"github.com/snowzach/gorestapi/store"
)

// GRStore is the persistent store of things
type GRStore interface {
	ThingGetByID(ctx context.Context, id string) (*Thing, error)
	ThingSave(ctx context.Context, thing *Thing) error
	ThingDeleteByID(ctx context.Context, id string) error
	ThingsFind(ctx context.Context, fqp *store.FindQueryParameters) ([]*Thing, int64, error)

	WidgetGetByID(ctx context.Context, id string) (*Widget, error)
	WidgetSave(ctx context.Context, thing *Widget) error
	WidgetDeleteByID(ctx context.Context, id string) error
	WidgetsFind(ctx context.Context, fqp *store.FindQueryParameters) ([]*Widget, int64, error)
}
