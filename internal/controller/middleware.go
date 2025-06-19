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
	method, ok := ctx.Value(ContextKeyMethodName).(string)
	if !ok || !traceIdForMethods[method] {
		return nil
	}

	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() || !spanCtx.IsSampled() {
		return nil
	}

	return prometheus.Labels{"trace_id": spanCtx.SpanID().String()}
}

func unaryMetadataInterceptor(defaultAcceptLang string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = context.WithValue(ctx, ContextKeyMethodName, info.FullMethod)

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ctx = extractMetadataToContext(ctx, md, defaultAcceptLang)
		}

		return handler(ctx, req)
	}
}

func streamMetadataInterceptor(defaultAcceptLang string) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := context.WithValue(ss.Context(), ContextKeyMethodName, info.FullMethod)

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ctx = extractMetadataToContext(ctx, md, defaultAcceptLang)
		}

		wrapped := middleware.WrapServerStream(ss)
		wrapped.WrappedContext = ctx

		return handler(srv, wrapped)
	}
}

func extractMetadataToContext(ctx context.Context, md metadata.MD, defAcceptLang string) context.Context {
	c := &models.Context{}
	c.Session = &models.Session{}

	if vals := md.Get(string(models.HeaderUserAgent)); len(vals) > 0 {
		c.UserAgent = vals[0]
	}
	if vals := md.Get(string(models.HeaderXRequestID)); len(vals) > 0 {
		c.RequestId = vals[0]
	}
	if vals := md.Get(models.HeaderAuthorization); len(vals) > 0 {
		c.Session.Token = vals[0]
	}
	if vals := md.Get(string(models.HeaderXIPAddress)); len(vals) > 0 {
		c.IPAddress = vals[0]
	}
	if vals := md.Get(string(models.HeaderXForwardedFor)); len(vals) > 0 {
		c.XForwardedFor = vals[0]
	}
	if vals := md.Get(models.HeaderAcceptLanguage); len(vals) > 0 {
		c.AcceptLanguage = vals[0]
	} else {
		c.AcceptLanguage = defAcceptLang
	}
	if vals := md.Get(models.HeaderSessionID); len(vals) > 0 {
		c.Session.Id = vals[0]
	}
	if vals := md.Get(models.HeaderToken); len(vals) > 0 {
		c.Session.Token = vals[0]
	}
	if vals := md.Get(models.HeaderCreatedAt); len(vals) > 0 {
		if val, err := strconv.Atoi(vals[0]); err == nil {
			c.Session.CreatedAt = int64(val)
		}
	}
	if vals := md.Get(models.HeaderLastActivityAt); len(vals) > 0 {
		if val, err := strconv.Atoi(vals[0]); err == nil {
			c.Session.LastActivityAt = int64(val)
		}
	}
	if vals := md.Get(models.HeaderUserID); len(vals) > 0 {
		c.Session.UserId = vals[0]
	}
	if vals := md.Get(models.HeaderDeviceID); len(vals) > 0 {
		c.Session.DeviceId = vals[0]
	}
	if vals := md.Get(models.HeaderRoles); len(vals) > 0 {
		c.Session.Roles = vals[0]
	}
	if vals := md.Get(models.HeaderProps); len(vals) > 0 {
		c.Session.Props = utils.GetMetadataValue(vals)
	}

	c.Context = ctx
	return context.WithValue(ctx, ContextKeyMetadata, c)
}
