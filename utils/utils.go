package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"hash/crc32"
	mathRand "math/rand"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

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

// Generate random string with length
func RandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-"

	ll := len(chars)
	b := make([]byte, length)
	rand.Read(b) // generates len(b) random bytes
	for i := 0; i < length; i++ {
		b[i] = chars[int(b[i])%ll]
	}
	return string(b)
}

// Generate random number with length
func RandomNumber(length int) string {
	chars := "0123456789"

	ll := len(chars)
	b := make([]byte, length)
	rand.Read(b) // generates len(b) random bytes
	for i := 0; i < length; i++ {
		b[i] = chars[int(b[i])%ll]
	}
	return string(b)
}

// Hash string using CRC32 algorithm and return int64
func HashStringCRC32(s string) int64 {
	return int64(crc32.ChecksumIEEE([]byte(s)))
}

func GenerateUID() string {
	mathRand.Seed(time.Now().UnixNano())
	arr := mathRand.Perm(10)

	uid := "UID"

	for _, v := range arr {
		uid += strconv.Itoa(v)
	}

	return uid
}

// Generate slice int from range start to end
func SliceIntRange(start, end int) []int {
	if start > end {
		return []int{}
	}

	slice := make([]int, end-start+1)
	for i := 0; i < len(slice); i++ {
		slice[i] = start + i
	}

	return slice
}

// Generate random hex with (n) length
func RandomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return strings.ToUpper(hex.EncodeToString(bytes))
}

// IsImageFile check file is image file.
func IsImageFile(file multipart.File) bool {
	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		return false
	}
	file.Seek(0, 0)

	return strings.HasPrefix(http.DetectContentType(fileHeader), "image/")
}

func ValidateImageFile(file multipart.File) bool {
	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		return false
	}
	file.Seek(0, 0)

	return strings.HasPrefix(http.DetectContentType(fileHeader), "image/")
}
