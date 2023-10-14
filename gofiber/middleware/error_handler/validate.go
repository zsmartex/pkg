package error_handler

import "github.com/gofiber/fiber/v2"

type ValidateError struct {
	Prefix string
	Method string
}

func ValidateHander(c *fiber.Ctx, method, prefix string) {
	c.Locals("validate_error_prefix", &ValidateError{
		Prefix: prefix,
		Method: method,
	})
}
