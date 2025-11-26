package controller

import (
	"context"
	"fmt"
	"time"

	sharedPb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/shared/v1"
	usersPb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/utils"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/worker"
	intModels "github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/hibiken/asynq"
	"google.golang.org/grpc/codes"
)

func (c *Controller) PasswordForgot(context context.Context, req *usersPb.PasswordForgotRequest) (*usersPb.PasswordForgotResponse, error) {
	path := "users.controller.PasswordForgot"
	errBuilder := func(e *models.AppError) (*usersPb.PasswordForgotResponse, error) {
		return &usersPb.PasswordForgotResponse{Response: &usersPb.PasswordForgotResponse_Error{Error: models.AppErrorToProto(e)}}, nil
	}
	sucBuilder := func(data *sharedPb.SuccessResponseData) (*usersPb.PasswordForgotResponse, error) {
		return &usersPb.PasswordForgotResponse{Response: &usersPb.PasswordForgotResponse_Data{Data: data}}, nil
	}
	internalErr := func(ctx *models.Context, err error) *models.AppError {
		return models.NewAppError(ctx, path, models.ErrMsgInternal, nil, "", int(codes.Internal), &models.AppErrorErrorsArgs{Err: err})
	}

	ctx, err := models.ContextGet(context)
	if err != nil {
		return errBuilder(err)
	}

	email := req.GetEmail()
	if !utils.IsValidEmail(email) {
		details := fmt.Sprintf("invalid email: %s ", email)
		return errBuilder(models.NewAppError(ctx, path, "email.invalid", nil, details, int(codes.InvalidArgument), nil))
	}

	ar := models.AuditRecordNew(ctx, intModels.EventNamePasswordForgot, models.EventStatusFail)
	defer c.ProcessAudit(ar)
	models.AuditEventDataParameter(ar, "email", email)

	user, dbErr := c.store.UsersGetByEmail(ctx, email)
	if dbErr != nil {
		if dbErr.ErrType == models.DBErrorTypeNoRows {
			return errBuilder(models.NewAppError(ctx, path, "email.not_found", nil, dbErr.Details, int(codes.NotFound), &models.AppErrorErrorsArgs{Err: dbErr}))
		} else {
			return errBuilder(internalErr(ctx, err))
		}
	}

	// SSO account has no password
	if user.GetAuthData() != "" {
		return errBuilder(models.NewAppError(ctx, path, "forgot.password.sso.error", nil, "", int(codes.InvalidArgument), nil))
	}

	token := &utils.Token{}
	tokenData, errTok := token.GenerateToken(time.Duration(time.Hour * time.Duration(c.config().Security.GetTokenPasswordResetExpiryInHours())))
	if errTok != nil {
		return errBuilder(internalErr(ctx, errTok))
	}

	_, dbErr = c.store.TokensDeleteAllPasswordResetByUserID(ctx, user.GetId())
	if dbErr != nil {
		return errBuilder(internalErr(ctx, dbErr))
	}

	dbErr = c.store.TokensAdd(ctx, user.GetId(), tokenData, intModels.TokenTypePasswordReset, path)
	if dbErr != nil {
		return errBuilder(internalErr(ctx, dbErr))
	}

	optoins := []asynq.Option{asynq.MaxRetry(10), asynq.ProcessIn(time.Second * 10), asynq.Queue(worker.QueuePriorityCritical)}
	taskPayload := &intModels.TaskSendPasswordResetEmailPayload{
		Ctx:     ctx,
		Email:   email,
		Token:   tokenData.Token,
		TokenID: tokenData.ID,
		Hours:   int(c.config().Security.GetTokenPasswordResetExpiryInHours()),
	}
	if err := c.tasker.SendPasswordResetEmail(context, taskPayload, optoins...); err != nil {
		return errBuilder(internalErr(ctx, err))
	}

	ar.Success()
	msg := models.Tr(ctx.AcceptLanguage, "forgot.password.success_message", nil)
	metadata := map[string]string{"description": models.Tr(ctx.AcceptLanguage, "forgot.password.success_message.description", map[string]any{"Email": req.GetEmail()})}
	return sucBuilder(&sharedPb.SuccessResponseData{Message: &msg, Metadata: metadata})
}
