// Package server binds everything required together for this service,
// E,g grpc, init metrics, oauth server, listen to errors, init clients....
package server

import (
	"context"
	"sync"

	com "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/common"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/controller"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/mailer"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/worker"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/logger"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Server struct {
	commonClient   *common.CommonClient
	configMux      sync.RWMutex
	configFn       func() *com.Config
	config         *com.Config
	errors         chan *models.InternalError
	objectStorage  *minio.Client
	tracerProvider *sdktrace.TracerProvider
	metrics        *grpcprom.ServerMetrics
	log            *logger.Logger
	dbConn         *pgxpool.Pool
	dbStore        store.UsersStore
	mailer         mailer.MailerService
	tasker         worker.TaskDistributor
}

type ServerArgs struct {
	Log *logger.Logger
	Cfg *models.Config
}

func RunServer(s *ServerArgs) error {
	com, err := common.NewCommonClient(&common.CommonArgs{Config: s.Cfg, Log: s.Log})
	app := &Server{
		commonClient: com,
		errors:       make(chan *models.InternalError, 1),
		log:          s.Log,
	}

	if err != nil {
		app.errors <- err
	}

	ctx := context.Background()

	app.initSharedConfig()
	app.initTrans()
	app.initTracerProvider(ctx)
	app.initObjectStorage()
	app.initMetrics()

	app.initDB()
	defer app.dbConn.Close()
	app.initStore()
	app.initMailer()
	app.initWorker()
	app.initOauthServer()

	_, err = controller.NewController(&controller.ControllerArgs{
		Cfg:            app.config,
		Store:          app.dbStore,
		ObjStorage:     app.objectStorage,
		TracerProvider: app.tracerProvider,
		Metrics:        app.metrics,
		Log:            app.log,
		Tasker:         app.tasker,
	})
	if err != nil {
		app.errors <- err
	}

	err = <-app.errors
	// TODO: cleanup things
	if err != nil {
		s.Log.Infof("an error occurred %v ", err)
		if err := app.tracerProvider.Shutdown(ctx); err != nil {
			s.Log.Errorf("failed to shutdown tracer provider %v", err)
		}
	}

	return err
}
