package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/snowzach/gorestapi/gorestapi"
	"github.com/snowzach/gorestapi/mocks"
)

func TestServerThingPost(t *testing.T) {

	// Mock Store and server
	ts := new(mocks.ThingStore)
	s, err := New(ts)
	assert.Nil(t, err)

	// Create Item
	i := &gorestapi.Thing{
		ID:   "id",
		Name: "name",
	}

	// Response
	idResponse := map[string]string{"id": i.ID}

	// Mock call to item store
	ts.On("ThingSave", mock.AnythingOfType("*context.valueCtx"), i).Once().Return(i.ID, nil)

	// Create test server
	server := httptest.NewServer(s.router)
	defer server.Close()

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.POST("/things").WithJSON(i).Expect().Status(http.StatusOK).JSON().Object().Equal(&idResponse)

	// Check remaining expectations
	ts.AssertExpectations(t)

}

func TestServerThingGetAll(t *testing.T) {

	// Mock Store and server
	ts := new(mocks.ThingStore)
	s, err := New(ts)
	assert.Nil(t, err)

	// Create Item
	i := []*gorestapi.Thing{
		&gorestapi.Thing{
			ID:   "id1",
			Name: "name1",
		},
		&gorestapi.Thing{
			ID:   "id2",
			Name: "name2",
		},
	}

	// Mock call to item store
	ts.On("ThingFind", mock.AnythingOfType("*context.valueCtx")).Once().Return(i, nil)

	// Create test server
	server := httptest.NewServer(s.router)
	defer server.Close()

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.GET("/things").Expect().Status(http.StatusOK).JSON().Array().Equal(&i)

	// Check remaining expectations
	ts.AssertExpectations(t)

}

func TestServerThingGet(t *testing.T) {

	// Mock Store and server
	ts := new(mocks.ThingStore)
	s, err := New(ts)
	assert.Nil(t, err)

	// Create Item
	i := &gorestapi.Thing{
		ID:   "id",
		Name: "name",
	}

	// Mock call to item store
	ts.On("ThingGetByID", mock.AnythingOfType("*context.valueCtx"), "1234").Once().Return(i, nil)

	// Create test server
	server := httptest.NewServer(s.router)
	defer server.Close()

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.GET("/things/1234").Expect().Status(http.StatusOK).JSON().Object().Equal(&i)

	// Check remaining expectations
	ts.AssertExpectations(t)

}

func TestServerThingDelete(t *testing.T) {

	// Mock Store and server
	ts := new(mocks.ThingStore)
	s, err := New(ts)
	assert.Nil(t, err)

	// Mock call to item store
	ts.On("ThingDeleteByID", mock.AnythingOfType("*context.valueCtx"), "1234").Once().Return(nil)

	// Create test server
	server := httptest.NewServer(s.router)
	defer server.Close()

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.DELETE("/things/1234").Expect().Status(http.StatusNoContent)

	// Check remaining expectations
	ts.AssertExpectations(t)

}
