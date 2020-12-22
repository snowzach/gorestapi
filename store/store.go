package store

import (
	"errors"
	"fmt"
)

// ErrNotFound is a standard not found error
var (
	ErrNotFound = errors.New("not found")
)

type ErrorType int
type ErrorOp int

const (
	ErrorTypeNone ErrorType = iota
	ErrorTypeIncomplete
	ErrorTypeForeignKey
	ErrorTypeDuplicate
	ErrorTypeInvalid
	ErrorTypeQuery

	ErrorOpSave ErrorOp = iota
	ErrorOpGet
	ErrorOpDelete
	ErrorOpFind
)

type Error struct {
	Type ErrorType
	Err  error
}

func (e *Error) Error() string { return e.Err.Error() }

func (e *Error) Unwrap() error { return e.Err }

func (e *Error) ErrorForOp(op ErrorOp) error {
	switch e.Type {
	case ErrorTypeNone:
		return nil
	case ErrorTypeIncomplete:
		return fmt.Errorf("missing data: %w", e.Err)
	case ErrorTypeForeignKey:
		return fmt.Errorf("foreign key: %w", e.Err)
	case ErrorTypeDuplicate:
		return fmt.Errorf("duplicate: %w", e.Err)
	case ErrorTypeInvalid:
		return fmt.Errorf("invalid data: %w", e.Err)
	}
	return e.Err
}

type Results struct {
	Count   int64       `json:"count"`
	Results interface{} `json:"results"`
}
