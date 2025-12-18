package errors

import "errors"

var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrSQLError            = errors.New("database server failed to execute query")
	ErrTooManyRequests     = errors.New("too many requests")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidUploadFile   = errors.New("invalid upload file")
	ErrSizeTooBig          = errors.New("size too big")
	ErrForbidden           = errors.New("forbidden")
	ErrSendEmail           = errors.New("error send mail")
)

var GeneralErrors = []error{
	ErrInternalServerError,
	ErrSQLError,
	ErrTooManyRequests,
	ErrUnauthorized,
	ErrInvalidToken,
	ErrInvalidUploadFile,
	ErrSizeTooBig,
	ErrForbidden,
	ErrSendEmail,
}
