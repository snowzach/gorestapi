package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/snowzach/gorestapi/gorestapi"
	"github.com/snowzach/gorestapi/mocks"
)

func TestServerThingPost(t *testing.T) {

	// Mock Store
	ts := new(mocks.ThingStore)

	// Server Setup
	s := &Server{
		logger:     zap.S().With("package", "api"),
		router:     chi.NewRouter(),
		thingStore: ts,
	}
	s.SetupRoutes()

	// Create Mock Item
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
