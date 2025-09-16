// Package controller contains the grpc handlers for this service
package controller

import (
	"net"
	"net/http"

	common "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/worker"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/logger"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var protectedMethods = map[string]bool{
	"/users.v1.UsersService/CreateSupplier": true,
}

var traceIDForMethods = map[string]bool{
	"/users.v1.UsersService/CreateSupplier": true,
}

type Controller struct {
	pb.UnimplementedUsersServiceServer
	store          store.UsersStore
	objStorage     *minio.Client
	config         func() *common.Config
	tracerProvider *sdktrace.TracerProvider
	metrics        *grpcprom.ServerMetrics
	log            *logger.Logger
	tasker         worker.TaskDistributor
	httpClient     *http.Client
}

type ControllerArgs struct {
	Config         func() *common.Config
	Store          store.UsersStore
	ObjStorage     *minio.Client
	TracerProvider *sdktrace.TracerProvider
	Metrics        *grpcprom.ServerMetrics
	Log            *logger.Logger
	Tasker         worker.TaskDistributor
}

func NewController(ca *ControllerArgs) (*Controller, *models.InternalError) {
	c := &Controller{
		config:         ca.Config,
		store:          ca.Store,
		objStorage:     ca.ObjStorage,
		tracerProvider: ca.TracerProvider,
		metrics:        ca.Metrics,
		log:            ca.Log,
		tasker:         ca.Tasker,
	}

	c.httpClient = utils.GetHTTPClient()

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(int(c.config().Services.GetUsersServiceMaxReceiveMessageSizeBytes())),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			c.responseInterceptor,
			c.unaryMetadataInterceptor(),
			c.metrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(traceID)),
			selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authMiddleware), selector.MatchFunc(authMatcher)),
		),
		grpc.ChainStreamInterceptor(
			c.streamMetadataInterceptor(),
			c.metrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(traceID)),
			selector.StreamServerInterceptor(auth.StreamServerInterceptor(authMiddleware), selector.MatchFunc(authMatcher)),
		),
	)

	addr := c.config().GetServices().GetUserServiceGrpcUrl()
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, &models.InternalError{Path: "user.controller.NewController", Err: err, Msg: "failed to initiate an http listener"}
	}

	reflection.Register(s)
	pb.RegisterUsersServiceServer(s, c)
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
