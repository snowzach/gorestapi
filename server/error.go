package server

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse is a generic struct for returning a standard error document
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// ErrNotFound is a pre-built not-found error
var ErrNotFound = &ErrResponse{HTTPStatusCode: http.StatusNotFound, StatusText: "Resource not found."}
var ErrUnauthorized = &ErrResponse{HTTPStatusCode: http.StatusUnauthorized, StatusText: "Not Authorized."}

// Render is the Renderer for ErrResponse struct
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrResourceNotFound(resource string) render.Renderer {
	return &ErrResponse{HTTPStatusCode: http.StatusNotFound, StatusText: resource + " not found."}
}

// ErrInvalidRequest is used to indicate an error on user input (with wrapped error)
func ErrInvalidRequest(err error) render.Renderer {
	var errorText string
	if err != nil {
		errorText = err.Error()
	}
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      errorText,
	}
}

// ErrInternal returns a generic server error to the user
func ErrInternal(err error) render.Renderer {
	var errorText string
	if err != nil {
		errorText = err.Error()
	}
	return &ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Server Error.",
		ErrorText:      errorText,
	}
}
