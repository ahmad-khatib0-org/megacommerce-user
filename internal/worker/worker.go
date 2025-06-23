package worker

import (
	"context"

	com "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/logger"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	SendVerifyEmail(ctx context.Context, pay *models.TaskSendVerifyEmailPayload, opts ...asynq.Option) *models.AppError
}

type TaskDistributorArgs struct {
	Log     *logger.Logger
	Options *asynq.RedisClientOpt
	Config  func() *com.Config
}

type AsynqTaksDistributor struct {
	cli     *asynq.Client
	log     *logger.Logger
	options *asynq.RedisClientOpt
	config  func() *com.Config
}

func NewAsynqTaksDistributor(tda *TaskDistributorArgs) TaskDistributor {
	cli := asynq.NewClient(tda.Options)
	return &AsynqTaksDistributor{cli: cli, log: tda.Log, options: tda.Options, config: tda.Config}
}
