package controller

import (
	"context"
	"fmt"
	"time"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/user/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/worker"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"github.com/hibiken/asynq"
	"google.golang.org/grpc/codes"
)

func (c *Controller) CreateSupplier(context context.Context, req *pb.SupplierCreateRequest) (*pb.SupplierCreateResponse, error) {
	path := "user.controller.SignupSupplier"
	errBuilder := func(e *models.AppError) (*pb.SupplierCreateResponse, error) {
		return &pb.SupplierCreateResponse{Response: &pb.SupplierCreateResponse_Error{Error: models.AppErrorToProto(e)}}, nil
	}

	ctx, err := models.ContextGet(context)
	if err != nil {
		return errBuilder(err)
	}

	ar := models.AuditRecordNew(ctx, models.EventNameSupplierCreate, models.EventStatusFail)
	models.AuditEventDataParameter(ar, "supplier", models.SignupSupplierRequestAuditable(req))
	defer c.ProcessAudit(ar)

	sanitized := models.SignupSupplierRequestSanitize(req)
	if err = models.SignupSupplierRequestIsValid(ctx, sanitized, c.cfg.Password); err != nil {
		return errBuilder(err)
	}

	dbPay, err := models.SignupSupplierRequestPreSave(
		ctx,
		&pb.User{
			Username:   utils.NewPointer(sanitized.GetUsername()),
			FirstName:  utils.NewPointer(sanitized.GetFirstName()),
			LastName:   utils.NewPointer(sanitized.GetLastName()),
			Email:      utils.NewPointer(sanitized.GetEmail()),
			Membership: utils.NewPointer("free"),
			Password:   utils.NewPointer(req.GetPassword()),
			Roles:      []string{string(models.RoleIdSupplierAdmin)},
		},
	)
	if err != nil {
		return errBuilder(err)
	}

	token := &utils.Token{}
	tokenData, errTok := token.GenerateToken(time.Duration(time.Hour * time.Duration(c.cfg.Security.GetTokenConfirmationExpiryInHours())))
	if errTok != nil {
		return errBuilder(err.ToInternal(err, nil))
	}

	if err := c.store.SignupSupplier(ctx, dbPay, tokenData); err != nil {
		if err.ErrType == store.DBErrorTypeUniqueViolation {
			details := fmt.Sprintf("the email %s is already in use", dbPay.GetEmail())
			return errBuilder(models.NewAppError(ctx, path, "user.create.email.not_unique", nil, details, int(codes.AlreadyExists), err))
		} else {
			return errBuilder(models.AppErrorInternal(err, ctx, err.Path, err.Msg))
		}
	}

	optoins := []asynq.Option{asynq.MaxRetry(10), asynq.ProcessIn(time.Second * 10), asynq.Queue(worker.QueuePriorityCritical)}
	taskPayload := &models.TaskSendVerifyEmailPayload{
		Ctx:   ctx,
		Email: dbPay.GetEmail(),
		Token: tokenData.Token,
		Hours: int(c.cfg.Security.GetTokenConfirmationExpiryInHours()),
	}

	// TODO: handle error
	if err := c.tasker.SendVerifyEmail(context, taskPayload, optoins...); err != nil {
		return errBuilder(models.AppErrorInternal(err, ctx, err.Where, err.Message))
	}

	ar.AuditEventDataResultState(models.SignupSupplierRequestResultState(dbPay))
	ar.Success()

	return &pb.SupplierCreateResponse{Response: &pb.SupplierCreateResponse_Data{}}, nil
}
