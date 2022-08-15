package jwt

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/volatiletech/null/v9"
)

func TestAuth_appendClaims(t *testing.T) {
	claims := jwt.MapClaims{
		"iat":          time.Now().Unix(),
		"jti":          strconv.FormatInt(time.Now().Unix(), 10),
		"exp":          time.Now().UTC().Add(time.Hour).Unix(),
		"sub":          "session",
		"iss":          "barong",
		"aud":          [3]string{"peatio", "barong", "finex"},
		"state":        "active",
		"referral_uid": "UID132132165",
	}

	t.Run("merges claims with nil", func(t *testing.T) {
		res := appendClaims(claims, nil)

		if !reflect.DeepEqual(claims, res) {
			t.Errorf("expected: %v actual: %v", claims, res)
		}
	})

	t.Run("merges nil with claims", func(t *testing.T) {
		res := appendClaims(nil, claims)
		if !reflect.DeepEqual(claims, res) {
			t.Errorf("expected: %v actual: %v", claims, res)
		}
	})

	t.Run("adds claim", func(t *testing.T) {
		res := appendClaims(claims, jwt.MapClaims{"custom": "claim"})

		if claims["custom"] != "claim" {
			t.Errorf("expected: %v actual: %v", claims, res)
		}
	})

	t.Run("rewrites claim", func(t *testing.T) {
		res := appendClaims(claims, jwt.MapClaims{"state": "banned"})

		if claims["state"] != "banned" {
			t.Errorf("expected: %v actual: %v", claims, res)
		}
	})
}

func TestAuth_JWT(t *testing.T) {
	ks, err := LoadOrGenerateKeys("./testdata/rsa-key", "./testdata/rsa-key.pub")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should validate jwt", func(t *testing.T) {
		token, err := ForgeToken("uid", "email", "role", null.StringFrom("UID123165658"), 3, false, ks.PrivateKey, nil)
		if err != nil {
			t.Fatal(err)
		}

		_, err = ParseAndValidate(token, ks.PublicKey)
		if err != nil {
			t.Fatal(err)
		}
	})
}
