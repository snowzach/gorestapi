package gorestapi

import (
	"encoding/json"
	"time"
)

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
	// Description
	Description string `json:"description"`
}

// ThingExample
// swagger:model gorestapi_ThingExample
type ThingExample struct {
	// Name
	Name string `json:"name"`
	// Description
	Description string `json:"description"`
}

// String is the stringer method
func (t *Thing) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}
