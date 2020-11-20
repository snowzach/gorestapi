package gorestapi

import (
	"encoding/json"
	"time"
)

// Widget
// swagger:model gorestapi_Widget
type Widget struct {
	// ID (Auto-Generated)
	ID string `json:"id"`
	// Created Timestamp
	Created time.Time `json:"created,omitempty"`
	// Updated Timestamp
	Updated time.Time `json:"updated,omitempty"`
	// Name
	Name string `json:"name"`
	// Description
	Description string `json:"description"`
	// ThingID
	ThingID *string `json:"thing_id,omitempty" db:"thing_id"`

	// Loaded Structs
	Thing *Thing `json:"thing,omitempty" db:"thing"`
}

// WidgetExample
// swagger:model gorestapi_WidgetExample
type WidgetExample struct {
	// Name
	Name string `json:"name"`
	// Description
	Description string `json:"description"`
	// ThingID
	ThingID *string `json:"thing_id,omitempty" db:"thing_id"`
}

// String is the stringer method
func (w *Widget) String() string {
	b, _ := json.Marshal(w)
	return string(b)
}

// SyncDB will fix Loaded Structs
func (w *Widget) SyncDB() {
	if w.ThingID == nil {
		w.Thing = nil
	}
}
