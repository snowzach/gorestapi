package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"

	"github.com/snowzach/gorestapi/conf"
	"github.com/snowzach/gorestapi/mocks"
)

func TestVersionGet(t *testing.T) {

	// Mock Store and server
	ts := new(mocks.ThingStore)
	s, err := New(ts)
	assert.Nil(t, err)

	// Create test server
	server := httptest.NewServer(s.router)
	defer server.Close()

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.GET("/version").Expect().Status(http.StatusOK).JSON().Object().Value("version").Equal(conf.GitVersion)

	// Check remaining expectations
	ts.AssertExpectations(t)

}
