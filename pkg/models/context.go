package models

import (
	"context"

	"google.golang.org/grpc/codes"
)

type ContextKey string

const (
	ContextKeyMetadata   ContextKey = "metadata"
	ContextKeyMethodName ContextKey = "method_name"
)

type StringMap map[string]string

type Context struct {
	Context        context.Context
	Session        *Session
	RequestId      string
	IPAddress      string
	XForwardedFor  string
	Path           string
	UserAgent      string
	AcceptLanguage string
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
			Id:            "server.internal.error",
			DetailedError: "failed to get the context from the incoming request",
			Where:         "user.models.ContextGet",
			StatusCode:    int(codes.Internal),
		}
	}

	return c, nil
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
