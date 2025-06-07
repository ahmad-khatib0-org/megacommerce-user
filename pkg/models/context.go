package models

import "context"

type Context struct {
	context        context.Context
	requestId      string
	ipAddress      string
	xForwardedFor  string
	path           string
	userAgent      string
	acceptLanguage string
}

func NewContext(ctx context.Context, requestId, ipAddress, xForwardedFor, path, userAgent, acceptLanguage string) *Context {
	return &Context{
		context:        ctx,
		requestId:      requestId,
		ipAddress:      ipAddress,
		xForwardedFor:  xForwardedFor,
		path:           path,
		userAgent:      userAgent,
		acceptLanguage: acceptLanguage,
	}
}

// clone creates a shallow copy of Context, allowing clones to apply per-request changes.
func (c *Context) clone() *Context {
	cp := *c
	return &cp
}

func (c *Context) Ctx() context.Context {
	return c.context
}

func (c *Context) RequestId() string {
	return c.requestId
}

func (c *Context) IPAddress() string {
	return c.ipAddress
}

func (c *Context) XForwardedFor() string {
	return c.xForwardedFor
}

func (c *Context) Path() string {
	return c.path
}

func (c *Context) UserAgent() string {
	return c.userAgent
}

func (c *Context) AcceptLanguage() string {
	return c.acceptLanguage
}

type RequestContext struct{}
