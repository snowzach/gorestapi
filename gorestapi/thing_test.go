package gorestapi

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThingString(t *testing.T) {
	thing := &Thing{
		ID:          "id1",
		Created:     time.Now(),
		Updated:     time.Now().Add(time.Minute),
		Name:        "name1",
		Description: "description1",
	}

	b, err := json.Marshal(thing)
	assert.Nil(t, err)

	assert.Equal(t, string(b), thing.String())
}
