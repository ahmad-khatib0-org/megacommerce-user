package models

import (
	"context"
	"time"

	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
)

type ContextKey string

const (
	ContextKeyMetadata   ContextKey = "metadata"
	ContextKeyMethodName ContextKey = "method_name"
)

type StringMap map[string]string

type Context struct {
	Context        context.Context `json:"-"`
	Session        *Session        `json:"session"`
	RequestId      string          `json:"request_id"`
	IPAddress      string          `json:"ip_address"`
	XForwardedFor  string          `json:"x_forwarded_for"`
	Path           string          `json:"path"`
	UserAgent      string          `json:"user_agent"`
	AcceptLanguage string          `json:"accept_language"`
}

func NewContext(ctx context.Context, session *Session, requestId, ipAddress, xForwardedFor, path, userAgent, acceptLanguage string) *Context {
	return &Context{
		Context:        ctx,
		Session:        session,
		RequestId:      requestId,
		IPAddress:      ipAddress,
		XForwardedFor:  xForwardedFor,
		Path:           path,
		UserAgent:      userAgent,
		AcceptLanguage: acceptLanguage,
	}
}

// clone creates a shallow copy of Context, allowing clones to apply per-request changes.
func (c *Context) clone() *Context {
	cp := *c
	return &cp
}

func (c *Context) Ctx() context.Context {
	return c.Context
}

func (c *Context) GetSession() *Session {
	return c.Session
}

func (c *Context) GetRequestId() string {
	return c.RequestId
}

func (c *Context) GetIPAddress() string {
	return c.IPAddress
}

func (c *Context) GetXForwardedFor() string {
	return c.XForwardedFor
}

func (c *Context) GetPath() string {
	return c.Path
}

func (c *Context) GetUserAgent() string {
	return c.UserAgent
}

func (c *Context) GetAcceptLanguage() string {
	return c.AcceptLanguage
}

func ContextGet(ctx context.Context) (*Context, *AppError) {
	c, ok := ctx.Value(ContextKeyMetadata).(*Context)
	if !ok {
		return nil, &AppError{
			Ctx:           c,
			Id:            ErrMsgInternal,
			DetailedError: "failed to get the context from the incoming request",
			Where:         "user.models.ContextGet",
			StatusCode:    int(codes.Internal),
		}
	}

	return c, nil
}

func ContextWith(ctx context.Context, appCtx *Context) context.Context {
	return context.WithValue(ctx, ContextKeyMetadata, appCtx)
}

// ContextForTesting get a context with dummy filled data for testing
func ContextForTesting() *Context {
	ctx := &Context{
		RequestId:      utils.NewID(),
		IPAddress:      gofakeit.IPv4Address(),
		XForwardedFor:  gofakeit.IPv4Address(),
		UserAgent:      gofakeit.UserAgent(),
		AcceptLanguage: "en",
		Session: &Session{
			Id:        utils.NewID(),
			Token:     utils.NewID(),
			CreatedAt: utils.TimeGetMillis(),
			ExpiresAt: utils.TimeGetMillis() + time.Duration(time.Hour).Milliseconds(),
			UserId:    utils.NewID(),
			DeviceId:  utils.NewID(),
			IsOAuth:   gofakeit.Bool(),
		},
	}

	return ctx
}

type Session struct {
	Id             string    `json:"id"`
	Token          string    `json:"token"`
	CreatedAt      int64     `json:"created_at"`
	ExpiresAt      int64     `json:"expires_at"`
	LastActivityAt int64     `json:"last_activity_at"`
	UserId         string    `json:"user_id"`
	DeviceId       string    `json:"device_id"`
	Roles          string    `json:"roles"`
	IsOAuth        bool      `json:"is_oauth"`
	Props          StringMap `json:"props"`
}

func (s *Session) GetId() string {
	return s.Id
}

func (s *Session) GetToken() string {
	return s.Token
}

func (s *Session) GetCreatedAt() float64 {
	return float64(s.CreatedAt)
}

func (s *Session) GetExpiresAt() float64 {
	return float64(s.ExpiresAt)
}

func (s *Session) GetLastActivityAt() float64 {
	return float64(s.LastActivityAt)
}

func (s *Session) GetUserID() string {
	return s.UserId
}

func (s *Session) GetDeviceID() string {
	return s.DeviceId
}

func (s *Session) GetRoles() string {
	return s.Roles
}

func (s *Session) GetIsAuth() bool {
	return s.IsOAuth
}

func (s *Session) GetProps() StringMap {
	return s.Props
}
