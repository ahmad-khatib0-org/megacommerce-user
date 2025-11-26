// Package controller contains the grpc handlers for this service
package controller

import (
	"net"
	"net/http"

	common "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/logger"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/utils"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/worker"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Controller struct {
	pb.UnimplementedUsersServiceServer
	store          store.UsersStore
	objStorage     *minio.Client
	config         func() *common.Config
	tracerProvider *sdktrace.TracerProvider
	log            *logger.Logger
	tasker         worker.TaskDistributor
	httpClient     *http.Client
}

type ControllerArgs struct {
	Config         func() *common.Config
	Store          store.UsersStore
	ObjStorage     *minio.Client
	TracerProvider *sdktrace.TracerProvider
	Log            *logger.Logger
	Tasker         worker.TaskDistributor
}

func NewController(ca *ControllerArgs) (*Controller, *models.InternalError) {
	c := &Controller{
		config:         ca.Config,
		store:          ca.Store,
		objStorage:     ca.ObjStorage,
		tracerProvider: ca.TracerProvider,
		log:            ca.Log,
		tasker:         ca.Tasker,
	}

	c.httpClient = utils.GetHTTPClient()

	defaultLang := c.config().Localization.GetDefaultClientLocale()
	availableLangs := c.config().GetLocalization().GetAvailableLocales()

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(int(c.config().Services.GetUsersServiceMaxReceiveMessageSizeBytes())),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			models.ResponseInterceptor(defaultLang, availableLangs),
			models.UnaryMetadataInterceptor(defaultLang, availableLangs),
			// c.metrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(traceID)),
			// selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authMiddleware), selector.MatchFunc(authMatcher)),
		),
		grpc.ChainStreamInterceptor(
			models.StreamMetadataInterceptor(defaultLang, availableLangs),
			// c.metrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(traceID)),
			// selector.StreamServerInterceptor(auth.StreamServerInterceptor(authMiddleware), selector.MatchFunc(authMatcher)),
		),
	)

	addr := c.config().GetServices().GetUserServiceGrpcUrl()
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, &models.InternalError{Path: "user.controller.NewController", Err: err, Msg: "failed to initiate an http listener"}
	}

	reflection.Register(s)
	pb.RegisterUsersServiceServer(s, c)
	// c.metrics.InitializeMetrics(s)

	go func() {
		c.log.Infof("grpc user service is running on %s", addr)
		if err := s.Serve(listener); err != nil {
			s.GracefulStop()
			s.Stop()
		}
	}()

	return c, nil
}
