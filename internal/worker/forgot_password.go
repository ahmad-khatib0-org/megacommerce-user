package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	intModels "github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/hibiken/asynq"
	"google.golang.org/grpc/codes"
)

// SendPasswordResetEmail implements SendPasswordResetEmail.
func (atp *AsynqTaksDistributor) SendPasswordResetEmail(context context.Context, payload *intModels.TaskSendPasswordResetEmailPayload, opts ...asynq.Option) *models.AppError {
	path := "user.worker.SendPasswordResetEmail"
	ctx, Err := models.ContextGet(context)
	if Err != nil {
		return Err
	}

	pay, err := json.Marshal(payload)
	if err != nil {
		return models.NewAppError(ctx, path, models.ErrMsgInternal, nil, fmt.Sprintf("failed to marshal json payload, err: %v", err), int(codes.Internal), &models.AppErrorErrorsArgs{Err: err})
	}

	task := asynq.NewTask(string(intModels.TaskNameSendPasswordResetEmail), pay, opts...)
	info, err := atp.cli.EnqueueContext(context, task)
	if err != nil {
		return models.NewAppError(ctx, path, models.ErrMsgInternal, nil, fmt.Sprintf("failed to enqueue a task , err: %v", err), int(codes.Internal), &models.AppErrorErrorsArgs{Err: err})
	}

	if atp.config().Main.GetEnv() == "dev" {
		atp.log.Infof("enqueued task: %v", info)
	}

	return nil
}

// ProcessSendPasswordResetEmail implements ProcessSendPasswordResetEmail.
func (atp *AsynqTaksProcessor) ProcessSendPasswordResetEmail(context context.Context, task *asynq.Task) error {
	path := "user.worker.ProcessSendPasswordResetEmail"
	var pay intModels.TaskSendPasswordResetEmailPayload
	if err := json.Unmarshal(task.Payload(), &pay); err != nil {
		return models.NewAppError(pay.Ctx, path, models.ErrMsgInternal, nil, fmt.Sprintf("failed to unmarshal json payload, err: %v", err), int(codes.Internal), &models.AppErrorErrorsArgs{Err: err})
	}

	if err := atp.mailer.SendPasswordResetEmail(pay.Ctx.GetAcceptLanguage(), pay.Email, pay.Token, pay.TokenID, pay.Hours); err != nil {
		return models.NewAppError(pay.Ctx, path, models.ErrMsgInternal, nil, fmt.Sprintf("failed to send an email, err: %v", err), int(codes.Internal), &models.AppErrorErrorsArgs{Err: err})
	}

	if atp.config().Main.GetEnv() == "dev" {
		atp.log.Infof("processed: %s task successfully", intModels.TaskNameSendVerifyEmail)
	}

	return nil
}
