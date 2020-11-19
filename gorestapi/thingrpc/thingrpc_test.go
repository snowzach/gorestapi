package thingrpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/snowzach/gorestapi/gorestapi"
	"github.com/snowzach/gorestapi/mocks"
	"github.com/snowzach/gorestapi/store"
)

func TestServerThingPost(t *testing.T) {

	// Create test server
	r := chi.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	// Mock Store and server
	ts := new(mocks.ThingStore)
	err := Setup(r, ts)
	assert.Nil(t, err)

	// Create Item
	i := &gorestapi.Thing{
		ID:   "id",
		Name: "name",
	}

	// Mock call to item store
	ts.On("ThingSave", mock.AnythingOfType("*context.valueCtx"), i).Once().Return(nil)

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.POST("/things").WithJSON(i).Expect().Status(http.StatusOK).JSON().Object().Equal(i)

	// Check remaining expectations
	ts.AssertExpectations(t)

}

func TestServerThingsFind(t *testing.T) {

	// Create test server
	r := chi.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	// Mock Store and server
	ts := new(mocks.ThingStore)
	err := Setup(r, ts)
	assert.Nil(t, err)

	// Return Item
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
	ts.On("ThingsFind", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*store.FindQueryParameters")).Once().Return(i, int64(2), nil)

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.GET("/things").Expect().Status(http.StatusOK).JSON().Object().Equal(&store.Results{Count: 2, Results: i})

	// Check remaining expectations
	ts.AssertExpectations(t)

}

func TestServerThingGetByID(t *testing.T) {

	// Create test server
	r := chi.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	// Mock Store and server
	ts := new(mocks.ThingStore)
	err := Setup(r, ts)
	assert.Nil(t, err)

	// Create Item
	i := &gorestapi.Thing{
		ID:   "id",
		Name: "name",
	}

	// Mock call to item store
	ts.On("ThingGetByID", mock.AnythingOfType("*context.valueCtx"), "1234").Once().Return(i, nil)

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.GET("/things/1234").Expect().Status(http.StatusOK).JSON().Object().Equal(&i)

	// Check remaining expectations
	ts.AssertExpectations(t)

}

func TestServerThingDeleteByID(t *testing.T) {

	// Create test server
	r := chi.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	// Mock Store and server
	ts := new(mocks.ThingStore)
	err := Setup(r, ts)
	assert.Nil(t, err)

	// Mock call to item store
	ts.On("ThingDeleteByID", mock.AnythingOfType("*context.valueCtx"), "1234").Once().Return(nil)

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.DELETE("/things/1234").Expect().Status(http.StatusNoContent)

	// Check remaining expectations
	ts.AssertExpectations(t)

}
