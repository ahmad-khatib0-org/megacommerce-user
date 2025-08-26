package worker

import (
	"context"
	"time"

	com "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/mailer"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/logger"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	Start() error
	ProcessSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	ProcessSendPasswordResetEmail(ctx context.Context, task *asynq.Task) error
}

const (
	QueuePriorityCritical = "critical"
	QueuePriorityDefault  = "default"
	QueuePriorityLow      = "low"
)

type TaskProcessorArgs struct {
	Store   store.UsersStore
	Config  func() *com.Config
	Mailer  mailer.MailerService
	Log     *logger.Logger
	Options *asynq.RedisClientOpt
}

type AsynqTaksProcessor struct {
	server  *asynq.Server
	store   store.UsersStore
	config  func() *com.Config
	mailer  mailer.MailerService
	options *asynq.RedisClientOpt
	log     *logger.Logger
}

func NewAsynqTaskProcessor(tpa *TaskProcessorArgs) TaskProcessor {
	server := asynq.NewServer(tpa.Options, asynq.Config{
		TaskCheckInterval:   time.Second * 1,  // default
		JanitorInterval:     time.Second * 8,  // default
		ShutdownTimeout:     time.Second * 8,  // default
		HealthCheckInterval: time.Second * 15, // default
		Queues: map[string]int{
			QueuePriorityCritical: 6,
			QueuePriorityDefault:  3,
			QueuePriorityLow:      1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			tpa.Log.ErrorStruct("user.worker.NewAsynqTaksProcessor", err)
		}),
	})

	return &AsynqTaksProcessor{server: server, store: tpa.Store, config: tpa.Config, mailer: tpa.Mailer, options: tpa.Options, log: tpa.Log}
}

// Start implements TaskProcessor.
func (atp *AsynqTaksProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(string(models.TaskNameSendVerifyEmail), atp.ProcessSendVerifyEmail)
	mux.HandleFunc(string(models.TaskNameSendPasswordResetEmail), atp.ProcessSendPasswordResetEmail)
	return atp.server.Start(mux)
}
