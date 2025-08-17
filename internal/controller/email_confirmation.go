package controller

import (
	"context"

	pbSh "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/shared/v1"
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
)

func (c *Controller) EmailConfirmation(context context.Context, req *pb.EmailConfirmationRequest) (*pb.EmailConfirmationResponse, error) {
	path := "users.controller.EmailConfirmation"
	errBuilder := func(e *models.AppError) (*pb.EmailConfirmationResponse, error) {
		return &pb.EmailConfirmationResponse{Response: &pb.EmailConfirmationResponse_Error{Error: models.AppErrorToProto(e)}}, nil
	}
	sucBuilder := func(data *pbSh.SuccessResponseData) (*pb.EmailConfirmationResponse, error) {
		return &pb.EmailConfirmationResponse{Response: &pb.EmailConfirmationResponse_Data{Data: data}}, nil
	}

	ctx, err := models.ContextGet(context)
	if err != nil {
		return errBuilder(err)
	}

	if err = models.EmailConfirmationIsValid(ctx, req); err != nil {
		return errBuilder(err)
	}

	ar := models.AuditRecordNew(ctx, models.EventNameEmailConfirmation, models.EventStatusFail)
	defer c.ProcessAudit(ar)

	token, errDB := c.store.TokensGet(ctx, req.TokenId)
	if errDB != nil {
		if errDB.ErrType == store.DBErrorTypeNoRows {
			return errBuilder(models.NewAppError(ctx, errDB.Path, "email_confirm.token.not_found", nil, "", int(codes.NotFound), err))
		} else {
			return errBuilder(models.AppErrorInternal(ctx, errDB, errDB.Path, ""))
		}
	}

	if token.Used {
		msg, _ := models.Tr(ctx.AcceptLanguage, "email_confirm.already_confirmed", nil)
		return sucBuilder(&pbSh.SuccessResponseData{Message: &msg})
	}

	if token.ExpiresAt < utils.TimeGetMillis() {
		return errBuilder(models.NewAppError(ctx, path, "email_confirm.token.expired", nil, "", int(codes.InvalidArgument), nil))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(token.Token), []byte(req.Token)); err != nil {
		msg, _ := models.Tr(ctx.AcceptLanguage, "email_confirm.token.error", nil)
		return sucBuilder(&pbSh.SuccessResponseData{Message: &msg})
	}

	if err := c.store.MarkEmailAsConfirmed(ctx, req.TokenId); err != nil {
		return errBuilder(models.AppErrorInternal(ctx, err, err.Path, err.Details))
	}

	ar.Success()
	msg, _ := models.Tr(ctx.AcceptLanguage, "email_confirm.confirmed_successfully", nil)
	return sucBuilder(&pbSh.SuccessResponseData{Message: &msg})
}
