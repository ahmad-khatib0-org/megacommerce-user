package server

import (
	"context"
	"sync"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/common"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/controller"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/logger"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Server struct {
	commonClient   *common.CommonClient
	configMux      sync.RWMutex
	config         *pb.Config
	done           chan *models.InternalError
	tracerProvider *sdktrace.TracerProvider
	metrics        *grpcprom.ServerMetrics
	log            *logger.Logger
}

type ServerArgs struct {
	Log *logger.Logger
	Cfg *models.Config
}

func RunServer(s *ServerArgs) error {
	com, err := common.NewCommonClient(s.Cfg)
	app := &Server{
		commonClient: com,
		done:         make(chan *models.InternalError, 1),
		log:          s.Log,
	}

	if err != nil {
		app.done <- err
	}

	ctx := context.Background()

	app.initSharedConfig()
	app.initTrans()
	app.initTracerProvider(ctx)
	app.initMetrics()

	_, err = controller.NewController(&controller.ControllerArgs{
		Cfg:            app.config,
		TracerProvider: app.tracerProvider,
		Metrics:        app.metrics,
		Log:            app.log,
	})
	if err != nil {
		app.done <- err
	}

	err = <-app.done
	if err != nil {
		if err := app.tracerProvider.Shutdown(ctx); err != nil {
			s.Log.Errorf("failed to shutdown tracer provider %v", err)
		}
		// TODO: cleanup things
	}

	return err
}
