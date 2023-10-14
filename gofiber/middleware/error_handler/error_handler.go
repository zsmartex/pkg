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
		unwrapError := errors.Unwrap(err)
		if unwrapError == nil {
			return c.Status(code).JSON(pkg.Error{
				Errors: pkg.ErrServerInternal.Errors,
			})
		}

		if e, ok := err.(*pkg.Error); ok {
			code = e.Code

			return c.Status(code).JSON(pkg.Error{
				Errors: e.Errors,
			})
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(pkg.ErrRecordNotFound)
		} else if utils.IsDuplicateKeyError(err) {
			errMsg := err.Error()
			errMsg = utils.TrimStringBetween(errMsg, "index_", "(")
			errMsg = strings.TrimSuffix(errMsg, "\" ")

			columnsStr := utils.TrimStringAfter(errMsg, "on_")

			columns := strings.Split(columnsStr, "_and_")

			for _, column := range columns {
				return pkg.NewError(fiber.StatusUnprocessableEntity, fmt.Sprintf("%s.taken", column), fmt.Sprintf("data in column %s already exists", column))
			}
		}
	}

	return c.Status(code).JSON(pkg.Error{
		Errors: pkg.ErrServerInternal.Errors,
	})
}
