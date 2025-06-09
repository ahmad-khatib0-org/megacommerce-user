package controller

import (
	"context"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// TODO: complete it
func authMiddleware(ctx context.Context) (context.Context, error) {
	_, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	return ctx, nil
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

func unaryMethodNameInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx = context.WithValue(ctx, ContextKeyMethodName, info.FullMethod)
		return handler(ctx, req)
	}
}

func streamMethodNameInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapped := middleware.WrapServerStream(ss)
		wrapped.WrappedContext = context.WithValue(ss.Context(), ContextKeyMethodName, info.FullMethod)
		return handler(srv, wrapped)
	}
}
