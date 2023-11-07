package session_fx

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/utils"
	"github.com/mileusna/useragent"
	"github.com/valyala/fasthttp"
	"github.com/zsmartex/pkg/v2/infrastructure/redis_fx"
	"go.uber.org/fx"
)

type SessData struct {
	ID              string    `json:"id,omitempty"`
	UID             string    `json:"uid,omitempty"`
	UserIP          string    `json:"user_ip,omitempty"`
	UserIPCountry   string    `json:"user_ip_country,omitempty"`
	UserAgent       string    `json:"user_agent,omitempty"`
	CrfsToken       string    `json:"crfs_token,omitempty"`
	AuthenticatedAt time.Time `json:"authenticated_at,omitempty"`
	CurrentSession  bool      `json:"-"`
}

type Session struct {
	fresh bool          // if new session
	id    string        // session id
	ctx   *fiber.Ctx    // fiber context
	exp   time.Duration // expiration of this session

	store *Store

	// data of session
	Data *SessData
}

func acquireSession() *Session {
	s := new(Session)

	if s.Data == nil {
		s.Data = new(SessData)
	}

	return s
}

type sessionStoreParams struct {
	fx.In

	RedisClient *redis_fx.Client
}

func NewStore(params sessionStoreParams) *Store {
	return &Store{
		CookieKey:      "zsmartex_id",
		CookiePath:     "/",
		CookieSecure:   false,
		CookieHTTPOnly: true,
		Expiration:     1 * time.Hour,
		RedisClient:    params.RedisClient,
	}
}

func (s *Session) Save(ctx context.Context) error {
	// Check if session has your own expiration, otherwise use default value
	if s.exp <= 0 {
		s.exp = s.store.Expiration
	}

	ua := useragent.Parse(s.Data.UserAgent)
	if ua.Mobile || ua.Tablet || ua.IsIOS() || ua.IsAndroid() {
		s.exp = time.Hour * 24 * 7
	}

	dataBytes, err := json.Marshal(s.Data)
	if err != nil {
		return err
	}

	s.setSession()

	if err := s.store.RedisClient.Set(ctx, fmt.Sprintf("_barong_session:%s", s.id), dataBytes, s.exp); err != nil {
		return err
	}

	if err := s.store.RedisClient.Set(ctx, fmt.Sprintf("_barong_user_session:%s:%s", s.Data.UID, s.id), s.id, s.exp); err != nil {
		return err
	}

	return nil
}

// Destroy will delete the session from Storage and expire session cookie
func (s *Session) Destroy(ctx context.Context) error {
	// Better safe than sorry
	if s.Data == nil {
		return nil
	}

	// Use external Storage if exist
	if err := s.store.RedisClient.Delete(ctx, fmt.Sprintf("_barong_session:%s", s.id)); err != nil {
		return err
	}

	if err := s.store.RedisClient.Delete(ctx, fmt.Sprintf("_barong_user_session:%s:%s", s.Data.UID, s.id)); err != nil {
		return err
	}

	// Expire session
	s.delSession()
	return nil
}

func (s *Session) SetExpire(exp time.Duration) {
	s.exp = exp
}

func (s *Session) SetUID(uid string) {
	s.Data.UID = uid
}

func (s *Session) SetUserIP(ip string) {
	s.Data.UserIP = ip
}

func (s *Session) SetUserIPCountry(country string) {
	s.Data.UserIPCountry = country
}

func (s *Session) SetUserAgent(ua string) {
	s.Data.UserAgent = ua
}

func (s *Session) SetAuthenticatedAt(at time.Time) {
	s.Data.AuthenticatedAt = at
}

func (s *Session) setSession() {
	fcookie := fasthttp.AcquireCookie()
	fcookie.SetKey(s.store.CookieKey)
	fcookie.SetValue(s.id)
	fcookie.SetPath(s.store.CookiePath)
	fcookie.SetDomain(s.store.CookieDomain)
	fcookie.SetMaxAge(int(s.exp.Seconds()))
	fcookie.SetExpire(time.Now().Add(s.exp))
	fcookie.SetSecure(s.store.CookieSecure)
	fcookie.SetHTTPOnly(s.store.CookieHTTPOnly)

	switch utils.ToLower(s.store.CookieSameSite) {
	case "strict":
		fcookie.SetSameSite(fasthttp.CookieSameSiteStrictMode)
	case "none":
		fcookie.SetSameSite(fasthttp.CookieSameSiteNoneMode)
	default:
		fcookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	}
	s.ctx.Response().Header.SetCookie(fcookie)
	fasthttp.ReleaseCookie(fcookie)
}

func (s *Session) delSession() {
	s.ctx.Request().Header.DelCookie(s.store.CookieKey)
	s.ctx.Response().Header.DelCookie(s.store.CookieKey)

	fcookie := fasthttp.AcquireCookie()
	fcookie.SetKey(s.store.CookieKey)
	fcookie.SetPath(s.store.CookiePath)
	fcookie.SetDomain(s.store.CookieDomain)
	fcookie.SetMaxAge(-1)
	fcookie.SetExpire(time.Now().Add(-1 * time.Minute))
	fcookie.SetSecure(s.store.CookieSecure)
	fcookie.SetHTTPOnly(s.store.CookieHTTPOnly)

	switch utils.ToLower(s.store.CookieSameSite) {
	case "strict":
		fcookie.SetSameSite(fasthttp.CookieSameSiteStrictMode)
	case "none":
		fcookie.SetSameSite(fasthttp.CookieSameSiteNoneMode)
	default:
		fcookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	}

	s.ctx.Response().Header.SetCookie(fcookie)
	fasthttp.ReleaseCookie(fcookie)
}
