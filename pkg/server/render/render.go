package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// JSON writes an http response using the value passed in v as JSON.
// If it cannot convert the value to JSON, it returns an error
func JSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(b, `{"render_error":"%s"}`, errString(err))
	} else {
		w.WriteHeader(status)
	}
	_, _ = w.Write(b.Bytes())
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

type ErrResponse struct {
	Status  string `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
	ErrorID string `json:"error_id,omitempty"`
}

type ErrOption func(er *ErrResponse)

func Err(w http.ResponseWriter, status int, opts ...ErrOption) {
	var err ErrResponse
	for _, opt := range opts {
		opt(&err)
	}
	JSON(w, status, err)
}

func WithStatus(status string) ErrOption {
	return func(er *ErrResponse) {
		er.Status = status
	}
}

func WithError(err error) ErrOption {
	return func(er *ErrResponse) {
		er.Error = errString(err)
	}
}

func WithErrorID(id string) ErrOption {
	return func(er *ErrResponse) {
		er.ErrorID = id
	}
}

func ErrNotFound(w http.ResponseWriter) {
	JSON(w, http.StatusNotFound, ErrResponse{Status: "not found", Error: "not found"})
}

func ErrResourceNotFound(w http.ResponseWriter, resource string) {
	JSON(w, http.StatusNotFound, ErrResponse{Status: resource + " not found", Error: resource + " not found"})
}

func ErrUnauthorizedWithID(w http.ResponseWriter, id string) {
	JSON(w, http.StatusUnauthorized, ErrResponse{Status: "not authorized", Error: "not authorized", ErrorID: id})
}

func ErrUnauthorized(w http.ResponseWriter) {
	ErrUnauthorizedWithID(w, "")
}

func ErrInvalidRequestWithID(w http.ResponseWriter, id string, err error) {
	JSON(w, http.StatusBadRequest, ErrResponse{Status: "invalid request", Error: errString(err), ErrorID: id})
}

func ErrInvalidRequest(w http.ResponseWriter, err error) {
	ErrInvalidRequestWithID(w, "", err)
}

func ErrInternalWithID(w http.ResponseWriter, id string, err error) {
	JSON(w, http.StatusInternalServerError, ErrResponse{Status: "internal error", Error: errString(err), ErrorID: id})
}

func ErrInternal(w http.ResponseWriter, err error) {
	ErrInternalWithID(w, "", err)
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func DecodeJSON(r io.Reader, v interface{}) error {
	defer io.Copy(ioutil.Discard, r)
	return json.NewDecoder(r).Decode(v)
}
