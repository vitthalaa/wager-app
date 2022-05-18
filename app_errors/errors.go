package app_errors

// ErrorCode ...
type ErrorCode string

const (
	ErrInvalidBody    ErrorCode = "INVALID_BODY"
	ErrInternalError  ErrorCode = "INTERNAL_ERROR"
	ErrNotImplemented ErrorCode = "NOT_IMPLEMENTED"
)
