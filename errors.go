package pkg

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
