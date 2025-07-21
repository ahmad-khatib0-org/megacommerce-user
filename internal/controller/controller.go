package controller

import (
	"net"

	common "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/user/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/worker"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/logger"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var protectedMethods = map[string]bool{
	"/user.v1.UserService/CreateSupplier": true,
}

var traceIdForMethods = map[string]bool{
	"/user.v1.UserService/CreateSupplier": true,
}

type Controller struct {
	pb.UnimplementedUserServiceServer
	store          store.UsersStore
	cfg            *common.Config
	tracerProvider *sdktrace.TracerProvider
	metrics        *grpcprom.ServerMetrics
	log            *logger.Logger
	tasker         worker.TaskDistributor
}

type ControllerArgs struct {
	Cfg            *common.Config
	TracerProvider *sdktrace.TracerProvider
	Metrics        *grpcprom.ServerMetrics
	Log            *logger.Logger
	Store          store.UsersStore
	Tasker         worker.TaskDistributor
}

func NewController(ca *ControllerArgs) (*Controller, *models.InternalError) {
	c := &Controller{
		cfg:            ca.Cfg,
		tracerProvider: ca.TracerProvider,
		metrics:        ca.Metrics,
		log:            ca.Log,
		store:          ca.Store,
		tasker:         ca.Tasker,
	}

	defaultAcceptLang := c.cfg.Localization.GetDefaultClientLocale()
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			unaryMetadataInterceptor(defaultAcceptLang),
			c.metrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(traceID)),
			selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authMiddleware), selector.MatchFunc(authMatcher)),
		),
		grpc.ChainStreamInterceptor(
			streamMetadataInterceptor(defaultAcceptLang),
			c.metrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(traceID)),
			selector.StreamServerInterceptor(auth.StreamServerInterceptor(authMiddleware), selector.MatchFunc(authMatcher)),
		),
	)

	addr := c.cfg.GetServices().GetUserServiceGrpcUrl()
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, &models.InternalError{Path: "user.controller.NewController", Err: err, Msg: "failed to initiate an http listener"}
	}

	reflection.Register(s)
	pb.RegisterUserServiceServer(s, c)
	c.metrics.InitializeMetrics(s)

	go func() {
		c.log.Infof("grpc user service is running on %s", addr)
		if err := s.Serve(listener); err != nil {
			s.GracefulStop()
			s.Stop()
		}
	}()

	return c, nil
}

// func InternalError(ctx *models.Context, err error) *models.AppError {
// 	return models.NewAppError(ctx, "user.controller.SignupSupplier", models.ErrMsgInternal, nil, "", int(codes.Internal), err)
// }
