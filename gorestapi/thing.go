package gorestapi

import (
	"context"
	"fmt"
	"time"

	"github.com/snowzach/gorestapi/store"
)

// ThingStore is the persistent store of things
type ThingStore interface {
	ThingGetByID(ctx context.Context, id string) (*Thing, error)
	ThingSave(ctx context.Context, thing *Thing) error
	ThingDeleteByID(ctx context.Context, id string) error
	ThingsFind(ctx context.Context, fqp *store.FindQueryParameters) ([]*Thing, int64, error)
}

// Thing
// swagger:model gorestapi_Thing
type Thing struct {
	// ID (Auto-Generated)
	ID string `json:"id"`
	// Created Timestamp
	Created time.Time `json:"created,omitempty"`
	// Updated Timestamp
	Updated time.Time `json:"updated,omitempty"`
	// Name
	Name string `json:"name"`
}

// ThingExample
// swagger:model gorestapi_ThingExample
type ThingExample struct {
	// Name
	Name string `json:"name"`
}

// String is the stringer method
func (t *Thing) String() string {
	return fmt.Sprintf(`{"id":"%s","name":"%s"}`, t.ID, t.Name)
}
