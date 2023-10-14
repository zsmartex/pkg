package pkg

import (
	"github.com/gofiber/fiber/v2"
)

var (
	ErrJWTDecodeAndVerify = NewError(fiber.StatusUnauthorized, "jwt.decode_and_verify", "failed to decode and verify jwt token")

	ErrServerInternal = NewError(fiber.StatusInternalServerError, "server.internal_error", "internal server error")

	ErrRecordNotFound = NewError(fiber.StatusNotFound, "record.not_found", "record not found")

	ErrServerInvalidQuery = NewError(fiber.StatusBadRequest, "server.method.invalid_message_query", "invalid query")

	ErrServerInvalidBody = NewError(fiber.StatusBadRequest, "server.method.invalid_message_body", "invalid body")
)

type Error struct {
	Errors      []string `json:"errors"`
	Code        int      `json:"-"`
	Description string   `json:"-"`
}

func NewError(code int, msg string, description string) *Error {
	return &Error{
		Errors: []string{
			msg,
		},
		Code:        code,
		Description: description,
	}
}

func (e *Error) Error() string {
	return e.Description
}
