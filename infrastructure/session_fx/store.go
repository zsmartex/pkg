package session_fx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"github.com/zsmartex/pkg/v2/infrastructure/redis_fx"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type Store struct {
	RedisClient *redis_fx.Client

	CookieKey      string
	CookieSameSite string
	CookiePath     string
	CookieDomain   string
	CookieSecure   bool
	CookieHTTPOnly bool
	Expiration     time.Duration
}

func (s *Store) GetSessionsByUID(ctx context.Context, uid string) ([]*SessData, error) {
	keys, err := s.RedisClient.Keys(ctx, fmt.Sprintf("_barong_user_session:%s:*", uid))
	if err != nil {
		return nil, err
	}

	sessionsData := make([]*SessData, 0)

	for _, key := range keys {
		result, err := s.RedisClient.Get(ctx, key)
		if err != nil {
			return nil, err
		}

		sessionID := result.Val()

		result, err = s.RedisClient.Get(ctx, fmt.Sprintf("_barong_session:%s", sessionID))
		if err != nil {
			return nil, err
		}

		bytes, err := result.Bytes()
		if err != nil {
			return nil, err
		}

		var sessionData *SessData
		if err := json.Unmarshal(bytes, &sessionData); err != nil {
			return nil, err
		}

		sessionsData = append(sessionsData, sessionData)
	}

	return sessionsData, nil
}

func (s *Store) DeleteSession(uid, sessionID string) error {
	exist, err := s.RedisClient.Exist(context.Background(), fmt.Sprintf("_barong_user_session:%s:%s", uid, sessionID))
	if err != nil {
		return err
	}

	if !exist {
		return ErrSessionNotFound
	}

	if err := s.RedisClient.Delete(context.Background(), fmt.Sprintf("_barong_user_session:%s:%s", uid, sessionID)); err != nil {
		return err
	}

	if err := s.RedisClient.Delete(context.Background(), fmt.Sprintf("_barong_session:%s", sessionID)); err != nil {
		return err
	}

	return nil
}

// Get will get/create a session
func (s *Store) Get(c *fiber.Ctx) (*Session, error) {
	var err error
	var fresh bool

	ctx := c.Context()

	id := s.getSessionID(c)
	if len(id) == 0 {
		fresh = true

		if id, err = s.responseCookies(c); err != nil {
			return nil, err
		}
	}

	if id == "" {
		id = uuid.New().String()
	}

	sess := acquireSession()
	sess.ctx = c
	sess.store = s
	sess.id = id
	sess.fresh = fresh
	sess.Data.ID = id

	if !fresh {
		key := fmt.Sprintf("_barong_session:%s", id)
		exist, err := s.RedisClient.Exist(ctx, key)
		if err != nil {
			return nil, err
		}

		if exist {
			result, err := s.RedisClient.Get(ctx, key)
			if err != nil {
				return nil, err
			}

			bytes, err := result.Bytes()
			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal(bytes, &sess.Data); err != nil {
				return nil, err
			}
		}
	}

	return sess, nil
}

// getSessionID will return the session id from cookie
func (s *Store) getSessionID(c *fiber.Ctx) string {
	return c.Cookies(s.CookieKey)
}

func (s *Store) responseCookies(c *fiber.Ctx) (string, error) {
	// Get key from response cookie
	cookieValue := c.Response().Header.PeekCookie(s.CookieKey)
	if len(cookieValue) == 0 {
		return "", nil
	}

	cookie := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(cookie)
	err := cookie.ParseBytes(cookieValue)
	if err != nil {
		return "", err
	}

	value := make([]byte, len(cookie.Value()))
	copy(value, cookie.Value())

	return string(value), nil
}
