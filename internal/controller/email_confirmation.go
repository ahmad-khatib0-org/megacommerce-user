package controller

import (
	"context"
	"time"

	pbSh "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/shared/v1"
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/utils"
	intModels "github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
)

func (c *Controller) EmailConfirmation(context context.Context, req *pb.EmailConfirmationRequest) (*pb.EmailConfirmationResponse, error) {
	start := time.Now()
	path := "users.controller.EmailConfirmation"
	errBuilder := func(e *models.AppError) (*pb.EmailConfirmationResponse, error) {
		return &pb.EmailConfirmationResponse{Response: &pb.EmailConfirmationResponse_Error{Error: models.AppErrorToProto(e)}}, nil
	}
	sucBuilder := func(data *pbSh.SuccessResponseData) (*pb.EmailConfirmationResponse, error) {
		return &pb.EmailConfirmationResponse{Response: &pb.EmailConfirmationResponse_Data{Data: data}}, nil
	}

	ctx, err := models.ContextGet(context)
	if err != nil {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordEmailConfirmationRequest(false, duration)
		return errBuilder(err)
	}

	if err = intModels.EmailConfirmationIsValid(ctx, req); err != nil {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordEmailConfirmationRequest(false, duration)
		return errBuilder(err)
	}

	ar := models.AuditRecordNew(ctx, intModels.EventNameEmailConfirmation, models.EventStatusFail)
	defer c.ProcessAudit(ar)

	token, errDB := c.store.TokensGet(ctx, req.TokenId)
	if errDB != nil {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordEmailConfirmationRequest(false, duration)
		if errDB.ErrType == models.DBErrorTypeNoRows {
			return errBuilder(models.NewAppError(ctx, path, "email_confirm.token.not_found", nil, "", int(codes.NotFound), &models.AppErrorErrorsArgs{Err: err}))
		} else {
			return errBuilder(models.NewAppError(ctx, path, models.ErrMsgInternal, nil, errDB.Details, int(codes.Internal), &models.AppErrorErrorsArgs{Err: err}))
		}
	}

	if token.Used {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordEmailConfirmationRequest(true, duration)
		msg := models.Tr(ctx.AcceptLanguage, "email_confirm.already_confirmed", nil)
		return sucBuilder(&pbSh.SuccessResponseData{Message: &msg})
	}

	if token.ExpiresAt < utils.TimeGetMillis() {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordEmailConfirmationRequest(false, duration)
		return errBuilder(models.NewAppError(ctx, path, "email_confirm.token.expired", nil, "", int(codes.InvalidArgument), nil))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(token.Token), []byte(req.Token)); err != nil {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordEmailConfirmationRequest(false, duration)
		msg := models.Tr(ctx.AcceptLanguage, "email_confirm.token.error", nil)
		return sucBuilder(&pbSh.SuccessResponseData{Message: &msg})
	}

	if err := c.store.MarkEmailAsConfirmed(ctx, req.TokenId); err != nil {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordEmailConfirmationRequest(false, duration)
		return errBuilder(models.NewAppError(ctx, path, models.ErrMsgInternal, nil, err.Details, int(codes.Internal), &models.AppErrorErrorsArgs{Err: err}))
	}

	ar.Success()
	duration := time.Since(start).Seconds()
	c.metricsCollector.RecordEmailConfirmationRequest(true, duration)
	msg := models.Tr(ctx.AcceptLanguage, "email_confirm.confirmed_successfully", nil)
	return sucBuilder(&pbSh.SuccessResponseData{Message: &msg})
}
