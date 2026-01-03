package errors

import (
	"net/http"

	"github.com/joomcode/errorx"
)

type ErrorType struct {
	StatusCode int
	Type       *errorx.Type
}

var Error = []ErrorType{
	{
		StatusCode: http.StatusBadRequest,
		Type:       ErrInvalidUserInput,
	},
	{
		StatusCode: http.StatusInternalServerError,
		Type:       ErrInternalServerError,
	},
	{
		StatusCode: http.StatusInternalServerError,
		Type:       ErrUnableToGet,
	},
	{
		StatusCode: http.StatusInternalServerError,
		Type:       ErrUnableToCreate,
	},
	{
		StatusCode: http.StatusNotFound,
		Type:       ErrResourceNotFound,
	},
}

// list of error namespaces
var (
	databaseError    = errorx.NewNamespace("database error").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	invalidInput     = errorx.NewNamespace("validation error").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	resourceNotFound = errorx.NewNamespace("not found").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	serverError      = errorx.NewNamespace("server error")
	dbError          = errorx.NewNamespace("db error")
	requestFailed    = errorx.NewNamespace("request binding Failed").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	bodyreadFailed   = errorx.NewNamespace("reading response body Failed").
				ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	pgtypeJsonbParseError = errorx.NewNamespace("failed to parse message data")
)

var (
	ErrUnableToCreate      = errorx.NewType(databaseError, "unable to create")
	ErrUnableToGet         = errorx.NewType(databaseError, "unable to get")
	ErrInvalidUserInput    = errorx.NewType(invalidInput, "invalid user input")
	ErrResourceNotFound    = errorx.NewType(resourceNotFound, "resource not found")
	ErrInternalServerError = errorx.NewType(serverError, "internal server error")
	ErrUnExpectedError     = errorx.NewType(serverError, "unexpected error occurred")
	ErrUnableToUpdate      = errorx.NewType(databaseError, "unable to update")
	ErrHTTPRequestBinding  = errorx.NewType(requestFailed, "binding failure")
	ErrReadingResponseBody = errorx.NewType(bodyreadFailed, "reading body failure")
	UnexpectedError        = errorx.NewType(serverError, "invalid value")
)
