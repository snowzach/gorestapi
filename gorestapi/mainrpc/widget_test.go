package mainrpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/snowzach/gorestapi/gorestapi"
	"github.com/snowzach/gorestapi/mocks"
	"github.com/snowzach/gorestapi/store"
)

func TestWidgetPost(t *testing.T) {

	// Create test server
	r := chi.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	// Mock Store and server
	grs := new(mocks.GRStore)
	err := Setup(r, grs)
	assert.Nil(t, err)

	// Create Item
	i := &gorestapi.Widget{
		ID:   "id",
		Name: "name",
	}

	// Mock call to item store
	grs.On("WidgetSave", mock.Anything, i).Once().Return(nil)

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.POST("/api/widgets").WithJSON(i).Expect().Status(http.StatusOK).JSON().Object().Equal(i)

	// Check remaining expectations
	grs.AssertExpectations(t)

}

func TestWidgetsFind(t *testing.T) {

	// Create test server
	r := chi.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	// Mock Store and server
	grs := new(mocks.GRStore)
	err := Setup(r, grs)
	assert.Nil(t, err)

	// Return Item
	i := []*gorestapi.Widget{
		&gorestapi.Widget{
			ID:   "id1",
			Name: "name1",
		},
		&gorestapi.Widget{
			ID:   "id2",
			Name: "name2",
		},
	}

	// Mock call to item store
	grs.On("WidgetsFind", mock.Anything, mock.AnythingOfType("*queryp.QueryParameters")).Once().Return(i, int64(2), nil)

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.GET("/api/widgets").Expect().Status(http.StatusOK).JSON().Object().Equal(&store.Results{Count: 2, Results: i})

	// Check remaining expectations
	grs.AssertExpectations(t)

}

func TestWidgetGetByID(t *testing.T) {

	// Create test server
	r := chi.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	// Mock Store and server
	grs := new(mocks.GRStore)
	err := Setup(r, grs)
	assert.Nil(t, err)

	// Create Item
	i := &gorestapi.Widget{
		ID:   "id",
		Name: "name",
	}

	// Mock call to item store
	grs.On("WidgetGetByID", mock.Anything, "1234").Once().Return(i, nil)

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.GET("/api/widgets/1234").Expect().Status(http.StatusOK).JSON().Object().Equal(&i)

	// Check remaining expectations
	grs.AssertExpectations(t)

}

func TestWidgetDeleteByID(t *testing.T) {

	// Create test server
	r := chi.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	// Mock Store and server
	grs := new(mocks.GRStore)
	err := Setup(r, grs)
	assert.Nil(t, err)

	// Mock call to item store
	grs.On("WidgetDeleteByID", mock.Anything, "1234").Once().Return(nil)

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.DELETE("/api/widgets/1234").Expect().Status(http.StatusNoContent)

	// Check remaining expectations
	grs.AssertExpectations(t)

}
