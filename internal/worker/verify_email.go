package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/hibiken/asynq"
	"google.golang.org/grpc/codes"
)

// SendVerifyEmail implements TaskDistributor.
func (a *AsynqTaksDistributor) SendVerifyEmail(context context.Context, payload *models.TaskSendVerifyEmailPayload, opts ...asynq.Option) *models.AppError {
	ctx, Err := models.ContextGet(context)
	if Err != nil {
		fmt.Println("from SendVerifyEmail ")
		return Err
	}

	pay, err := json.Marshal(payload)
	if err != nil {
		return models.NewAppError(ctx, "user.worker.SendVerifyEmail", "server.internal.error", nil, fmt.Sprintf("failed to marshal json payload, err: %v", err), int(codes.Internal), err)
	}

	task := asynq.NewTask(string(models.TaskNameSendVerifyEmail), pay, opts...)
	info, err := a.cli.EnqueueContext(context, task)
	if err != nil {
		return models.NewAppError(ctx, "user.worker.SendVerifyEmail", "server.internal.error", nil, fmt.Sprintf("failed to enqueue a task , err: %v", err), int(codes.Internal), err)
	}

	if a.config().Main.GetEnv() == "dev" {
		a.log.Infof("enqueued task: %v", info)
	}

	return nil
}

// ProcessSendVerifyEmail implements TaskProcessor.
func (a *AsynqTaksProcessor) ProcessSendVerifyEmail(context context.Context, task *asynq.Task) error {
	var pay models.TaskSendVerifyEmailPayload
	if err := json.Unmarshal(task.Payload(), &pay); err != nil {
		return models.NewAppError(pay.Ctx, "user.worker.ProcessSendVerifyEmail", "server.internal.error", nil, fmt.Sprintf("failed to unmarshal json payload, err: %v", err), int(codes.Internal), err)
	}

	if err := a.mailer.SendVerifyEmail(pay.Ctx.GetAcceptLanguage(), pay.Email, pay.Email, pay.Hours); err != nil {
		return models.NewAppError(pay.Ctx, "user.worker.ProcessSendVerifyEmail", "server.internal.error", nil, fmt.Sprintf("failed to send an email, err: %v", err), int(codes.Internal), err)
	}

	if a.config().Main.GetEnv() == "dev" {
		a.log.Infof("processed: %s task successfully", models.TaskNameSendVerifyEmail)
	}

	return nil
}
