package server

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/snowzach/gorestapi/conf"
)

// GetVersion returns version
func GetVersion() http.HandlerFunc {

	// Simple version struct
	type version struct {
		Version string `json:"version"`
	}
	var v = &version{Version: conf.GitVersion}

	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, v)
	}
}
