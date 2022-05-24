package app_errors

// ErrorCode ...
type ErrorCode string

const (
	ErrInvalidBody    ErrorCode = "INVALID_BODY"
	ErrInternalError  ErrorCode = "INTERNAL_ERROR"
	ErrNotImplemented ErrorCode = "NOT_IMPLEMENTED"
	ErrNotFound       ErrorCode = "NOT_FOUND"

	ErrInvalidTotalWagerValue   ErrorCode = "INVALID_TOTAL_WAGER_VALUE"
	ErrInvalidOdds              ErrorCode = "INVALID_ODDS"
	ErrInvalidSellingPercentage ErrorCode = "INVALID_SELLING_PERCENTAGE"
	ErrInvalidSellingPrice      ErrorCode = "INVALID_SELLING_PRICE"

	ErrInvalidWagerID     ErrorCode = "INVALID_WAGER_ID"
	ErrInvalidBuyingPrice ErrorCode = "INVALID_BUYING_PRICE"
	ErrWagerSoldOut       ErrorCode = "WAGER_SOLD_OUT"
)

// ErrorResponse is response object for errors
type ErrorResponse struct {
	Status int       `json:"-"`
	Code   ErrorCode `json:"error"`
}

// Error returns error string message
func (e *ErrorResponse) Error() string {
	return string(e.Code)
}
