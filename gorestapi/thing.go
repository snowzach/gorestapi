package gorestapi

import (
	"context"
	"fmt"
)

// ThingStore is the persistent store of things
type ThingStore interface {
	ThingGetByID(context.Context, string) (*Thing, error)
	ThingSave(context.Context, *Thing) (string, error)
	ThingDeleteByID(context.Context, string) error
	ThingFind(context.Context) ([]*Thing, error)
}

// Thing is an example struct
type Thing struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// String is the stringer method
func (t *Thing) String() string {
	return fmt.Sprintf(`{"id":"%s","name":"%s"}`, t.ID, t.Name)
}
