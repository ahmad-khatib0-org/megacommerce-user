package controller

import (
	"context"
	"strconv"

	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TODO: complete it
func authMiddleware(ctx context.Context) (context.Context, error) {
	_, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func authMatcher(ctx context.Context, callMeta interceptors.CallMeta) bool {
	_, ok := protectedMethods[callMeta.Method]
	return ok
}

func traceID(ctx context.Context) prometheus.Labels {
	method, ok := ctx.Value(models.ContextKeyMethodName).(string)
	if !ok || !traceIDForMethods[method] {
		return nil
	}

	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() || !spanCtx.IsSampled() {
		return nil
	}

	return prometheus.Labels{"trace_id": spanCtx.SpanID().String()}
}

func (c *Controller) unaryMetadataInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = context.WithValue(ctx, models.ContextKeyMethodName, info.FullMethod)

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ctx = c.extractMetadataToContext(ctx, md)
		}

		return handler(ctx, req)
	}
}

func (c *Controller) streamMetadataInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := context.WithValue(ss.Context(), models.ContextKeyMethodName, info.FullMethod)

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ctx = c.extractMetadataToContext(ctx, md)
		}

		wrapped := middleware.WrapServerStream(ss)
		wrapped.WrappedContext = ctx

		return handler(srv, wrapped)
	}
}

func (c *Controller) extractMetadataToContext(ctx context.Context, md metadata.MD) context.Context {
	mc := &models.Context{}
	mc.Session = &models.Session{}
	defaultAcceptLang := c.cfg.Localization.GetDefaultClientLocale()
	availableLocales := c.cfg.Localization.GetAvailableLocales()

	if vals := md.Get(string(models.HeaderUserAgent)); len(vals) > 0 {
		mc.UserAgent = vals[0]
	}
	if vals := md.Get(string(models.HeaderXRequestID)); len(vals) > 0 {
		mc.RequestID = vals[0]
	}
	if vals := md.Get(models.HeaderAuthorization); len(vals) > 0 {
		mc.Session.Token = vals[0]
	}
	if vals := md.Get(string(models.HeaderXIPAddress)); len(vals) > 0 {
		mc.IPAddress = vals[0]
	}
	if vals := md.Get(string(models.HeaderXForwardedFor)); len(vals) > 0 {
		mc.XForwardedFor = vals[0]
	}
	if vals := md.Get(models.HeaderAcceptLanguage); len(vals) > 0 {
		mc.AcceptLanguage = utils.ProcessAcceptedLanguage(vals[0], availableLocales, defaultAcceptLang)
	} else {
		mc.AcceptLanguage = defaultAcceptLang
	}
	if vals := md.Get(models.HeaderSessionID); len(vals) > 0 {
		mc.Session.ID = vals[0]
	}
	if vals := md.Get(models.HeaderToken); len(vals) > 0 {
		mc.Session.Token = vals[0]
	}
	if vals := md.Get(models.HeaderCreatedAt); len(vals) > 0 {
		if val, err := strconv.Atoi(vals[0]); err == nil {
			mc.Session.CreatedAt = int64(val)
		}
	}
	if vals := md.Get(models.HeaderLastActivityAt); len(vals) > 0 {
		if val, err := strconv.Atoi(vals[0]); err == nil {
			mc.Session.LastActivityAt = int64(val)
		}
	}
	if vals := md.Get(models.HeaderUserID); len(vals) > 0 {
		mc.Session.UserID = vals[0]
	}
	if vals := md.Get(models.HeaderDeviceID); len(vals) > 0 {
		mc.Session.DeviceID = vals[0]
	}
	if vals := md.Get(models.HeaderRoles); len(vals) > 0 {
		mc.Session.Roles = vals[0]
	}
	if vals := md.Get(models.HeaderProps); len(vals) > 0 {
		mc.Session.Props = utils.GetMetadataValue(vals)
	}
	if vals := md.Get(models.HeaderServerName); len(vals) > 0 {
		mc.ServerName = vals[0]
	}

	mc.Context = ctx
	return context.WithValue(ctx, models.ContextKeyMetadata, mc)
}
