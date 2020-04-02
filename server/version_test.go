package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/go-chi/chi"

	"github.com/snowzach/gorestapi/conf"
)

func TestVersionGet(t *testing.T) {

	// Create test server
	r := chi.NewRouter()
	server := httptest.NewServer(r)
	defer server.Close()

	r.Get("/version", GetVersion())

	// Make request and validate we get back proper response
	e := httpexpect.New(t, server.URL)
	e.GET("/version").Expect().Status(http.StatusOK).JSON().Object().Value("version").Equal(conf.GitVersion)

}
