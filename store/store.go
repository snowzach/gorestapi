package store

import (
	"errors"
)

// ErrNotFound is a standard not found error
var (
	ErrNotFound = errors.New("not found")
	ErrInteral  = &InternalError{Err: errors.New("internal error")}
)

type InternalError struct {
	Err error
}

func (e *InternalError) Error() string { return e.Err.Error() }

func (e *InternalError) Unwrap() error { return e.Err }

type Results struct {
	Count   int64       `json:"count"`
	Results interface{} `json:"results"`
}
