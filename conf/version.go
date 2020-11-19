package conf

import (
	"net/http"

	"github.com/go-chi/render"
)

var (
	// Executable is overridden by Makefile with executable name
	Executable = "NoExecutable"
	// GitVersion is overridden by Makefile with git information
	GitVersion = "NoGitVersion"
)

// GetVersion returns version
func GetVersion() http.HandlerFunc {

	// Simple version struct
	type version struct {
		Version string `json:"version"`
	}
	var v = &version{Version: GitVersion}

	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, v)
	}
}
