package gorestapi

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWidgetString(t *testing.T) {

	thingID := "tid1"

	widget := &Widget{
		ID:          "id1",
		Created:     time.Now(),
		Updated:     time.Now().Add(time.Minute),
		Name:        "name1",
		Description: "description1",
		ThingID:     &thingID,
	}

	b, err := json.Marshal(widget)
	assert.Nil(t, err)

	assert.Equal(t, string(b), widget.String())
}
