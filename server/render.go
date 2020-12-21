package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rs/xid"
)

// RenderJSON writes an http response using the value passed in v as JSON.
// If it cannot convert the value to JSON, it returns an error
func RenderJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(b, `{"render_error":"%s"}`, errString(err))
	} else {
		w.WriteHeader(code)
	}
	_, _ = w.Write(b.Bytes())
}

func RenderNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

type ErrResponse struct {
	Status  string `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
	ErrorID string `json:"error_id,omitempty"`
}

func RenderErrNotFound(w http.ResponseWriter) {
	RenderJSON(w, http.StatusNotFound, ErrResponse{Status: "not found", Error: "not found"})
}

func RenderErrResourceNotFound(w http.ResponseWriter, resource string) {
	RenderJSON(w, http.StatusNotFound, ErrResponse{Status: resource + " not found", Error: resource + " not found"})
}

func RenderErrUnauthorized(w http.ResponseWriter) {
	RenderJSON(w, http.StatusUnauthorized, ErrResponse{Status: "not authorized", Error: "not authorized"})
}

func RenderErrInvalidRequest(w http.ResponseWriter, err error) {
	RenderJSON(w, http.StatusBadRequest, ErrResponse{Status: "invalid request", Error: errString(err)})
}

func RenderErrInternal(w http.ResponseWriter, err error) {
	RenderJSON(w, http.StatusInternalServerError, ErrResponse{Status: "internal error", Error: errString(err)})
}

func RenderErrInternalWithID(w http.ResponseWriter, err error) string {
	errID := xid.New().String()
	RenderJSON(w, http.StatusInternalServerError, ErrResponse{Status: "internal error", Error: errString(err), ErrorID: errID})
	return errID
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
