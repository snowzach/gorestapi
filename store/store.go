package store

import (
	"errors"
)

// ErrNotFound is a standard not found error
var ErrNotFound = errors.New("not found")

type Results struct {
	Count   int64       `json:"count"`
	Results interface{} `json:"results"`
}
