package errors

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (

	// Validation

	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT_ERROR"

	// External Services
	ErrCodeCurrencyAPIRequest  ErrorCode = "CURRENCY_API_REQUEST_ERROR"
	ErrCodeCurrencyAPIResponse ErrorCode = "CURRENCY_API_RESPONSE_ERROR"
	ErrCodeCurrencyAPIParsing  ErrorCode = "CURRENCY_API_PARSING_ERROR"
	ErrCodeCurrencyNotFound    ErrorCode = "CURRENCY_NOT_FOUND"

	ErrCodeTelegramBotInit ErrorCode = "TELEGRAM_BOT_INIT_ERROR"
	ErrCodeTelegramBotSend ErrorCode = "TELEGRAM_BOT_SEND_ERROR"

	//Configuration

	ErrCodeConfig ErrorCode = "CONFIG_ERROR"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func IsErrorCode(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}
