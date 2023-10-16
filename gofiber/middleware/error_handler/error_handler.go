package error_handler

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/zsmartex/pkg/v2"
	"github.com/zsmartex/pkg/v2/log"
	"github.com/zsmartex/pkg/v2/utils"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	log.Errorf("%+v", err)

	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		// Override status code if fiber.Error type
		code = e.Code

		switch code {
		case fiber.StatusNotFound:
			return c.Status(code).JSON("404 Not Found")
		}
	} else {
		_err := errors.UnwrapAll(err)
		if _err != nil {
			err = _err
		}

		if e, ok := err.(*pkg.Error); ok {
			code = e.Code

			returnedMessages := make([]string, 0)
			for _, msg := range e.Errors {
				if !strings.Contains(msg, ".") {
					validateError, ok := c.Locals("validate_error_prefix").(*ValidateError)
					if ok {
						returnedMessages = append(returnedMessages, fmt.Sprintf("%s.%s.%s", validateError.Prefix, validateError.Method, msg))
					} else {
						if len(e.Errors) == 1 {
							return c.Status(code).JSON(pkg.ErrServerInternal)
						}
					}
				} else {
					returnedMessages = append(returnedMessages, msg)
				}
			}

			return c.Status(code).JSON(pkg.Error{
				Errors: returnedMessages,
			})
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(pkg.ErrRecordNotFound)
		} else if utils.IsDuplicateKeyError(err) {
			errMsg := err.Error()
			errMsg = utils.TrimStringBetween(errMsg, "index_", "(")
			errMsg = strings.TrimSuffix(errMsg, "\" ")

			columnsStr := utils.TrimStringAfter(errMsg, "on_")
			log.Info(columnsStr)

			columns := strings.Split(columnsStr, "_and_")
			log.Info(columns)

			for _, column := range columns {
				return pkg.NewError(fiber.StatusUnprocessableEntity, fmt.Sprintf("%s.taken", column), fmt.Sprintf("data in column %s already exists", column))
			}
		}
	}

	return c.Status(code).JSON(pkg.Error{
		Errors: pkg.ErrServerInternal.Errors,
	})
}
