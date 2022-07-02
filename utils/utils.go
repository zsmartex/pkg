package utils

import (
	"errors"
	"reflect"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"gorm.io/gorm"
)

func Contains[T any](arr []T, v T) bool {
	for _, a := range arr {
		if reflect.DeepEqual(a, v) {
			return true
		}
	}
	return false
}

// remove uniq elements
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, gorm.ErrRecordNotFound)
}

func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	if pqErr, ok := err.(*pgconn.PgError); ok {
		return pqErr.Code == pgerrcode.UniqueViolation
	}

	return false
}

func RemoveDuplicate[T string | int | int64 | float64](slice []T) []T {
	keys := make(map[T]bool)
	list := []T{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

func TrimStringBetween(str string, a, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(str, a)
	if posFirst == -1 {
		return ""
	}

	posLast := strings.Index(str, b)
	if posLast == -1 {
		return ""
	}

	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}

	return str[posFirstAdjusted:posLast]
}

func TrimStringAfter(str string, a string) string {
	// Get substring after string.
	pos := strings.Index(str, a)
	if pos == -1 {
		return ""
	}

	posAdjusted := pos + len(a)
	if posAdjusted >= len(str) {
		return ""
	}

	return str[posAdjusted:]
}
