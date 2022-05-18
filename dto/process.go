package dto

import "github.com/vitthalaa/wager-app/app_errors"

// ErrorResponse is response object for errors
type ErrorResponse struct {
	Error app_errors.ErrorCode `json:"error"`
}
