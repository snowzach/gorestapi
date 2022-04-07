package postgres

import (
	"github.com/jackc/pgconn"
	"github.com/snowzach/gorestapi/store"
)

// Lookup of postgres error codes to basic errors we can return to a user
var pgErrorCodeToStoreErrorType = map[string]store.ErrorType{
	"23502": store.ErrorTypeIncomplete,
	"23503": store.ErrorTypeForeignKey,
	"23505": store.ErrorTypeDuplicate,
	"23514": store.ErrorTypeInvalid,
}

func WrapError(err error) error {
	switch e := err.(type) {
	case *pgconn.PgError:
		if et, found := pgErrorCodeToStoreErrorType[e.Code]; found {
			return &store.Error{
				Type: et,
				Err:  err,
			}
		}
	}
	return err
}
