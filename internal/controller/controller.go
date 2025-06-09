package controller

import (
	"context"
	"fmt"
	"net"

	v1 "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/user/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

type ContextKey string

const (
	ContextKeyMetadata   ContextKey = "metadata"
	ContextKeyMethodName ContextKey = "method_name"
)

var protectedMethods = map[string]bool{
	"/user.v1.UserService/CreateSupplier": true,
}

var traceIdForMethods = map[string]bool{
	"/user.v1.UserService/CreateSupplier": true,
}

type Controller struct {
	pb.UnimplementedUserServiceServer
	cfg            *v1.Config
	tracerProvider *sdktrace.TracerProvider
	metrics        *grpcprom.ServerMetrics
}

type ControllerArgs struct {
	Cfg            *v1.Config
	TracerProvider *sdktrace.TracerProvider
	Metrics        *grpcprom.ServerMetrics
}

func NewController(ca *Controller) (*Controller, *models.InternalError) {
	c := &Controller{cfg: ca.cfg, tracerProvider: ca.tracerProvider, metrics: ca.metrics}

	authMatcher := func(ctx context.Context, callMeta interceptors.CallMeta) bool {
		_, ok := protectedMethods[callMeta.Method]
		return ok
	}

	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			unaryMethodNameInterceptor(),
			c.metrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(traceID)),
			selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authMiddleware), selector.MatchFunc(authMatcher)),
		),
		grpc.ChainStreamInterceptor(
			streamMethodNameInterceptor(),
			c.metrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(traceID)),
			selector.StreamServerInterceptor(auth.StreamServerInterceptor(authMiddleware), selector.MatchFunc(authMatcher)),
		),
	)

	addr := fmt.Sprintf("%s:%d", c.cfg.GetServices().GetUserServiceGrpcHost(), c.cfg.GetServices().GetUserServiceGrpcPort())
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, &models.InternalError{Path: "user.controller.NewController", Err: err, Msg: "failed to initiate an http listener"}
	}

	go func() {
		if err := s.Serve(listener); err != nil {
			s.GracefulStop()
			s.Stop()
		}
	}()

	return c, nil
}
