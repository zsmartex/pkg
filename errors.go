package pkg

import "fmt"

type Error struct {
	Errors []string `json:"errors"`
	Code   int      `json:"-"`
}

func NewError(code int, message ...string) *Error {
	return &Error{
		Errors: message,
		Code:   code,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Errors)
}
