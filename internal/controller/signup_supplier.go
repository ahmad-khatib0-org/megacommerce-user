package controller

import (
	"context"
	"fmt"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/user/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"google.golang.org/grpc/codes"
)

func (c *Controller) CreateSupplier(context context.Context, req *pb.SupplierCreateRequest) (*pb.SupplierCreateResponse, error) {
	ctx, err := getContext(context)
	errBuilder := func(e *models.AppError) (*pb.SupplierCreateResponse, error) {
		return &pb.SupplierCreateResponse{Response: &pb.SupplierCreateResponse_Error{Error: models.AppErrorToProto(e)}}, nil
	}

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
			Password:   utils.NewPointer(sanitized.GetPassword()),
		},
	)
	if err != nil {
		return errBuilder(err)
	}

	if err := c.store.SignupSupplier(ctx, dbPay); err != nil {
		if err.ErrType == store.DBErrorTypeUniqueViolation {
			return errBuilder(
				models.NewAppError(
					ctx, "user.controller.SignupSupplier",
					"user.create.email.not_unique", nil,
					fmt.Sprintf("the email: %s is already in use", dbPay.GetEmail()),
					int(codes.AlreadyExists), err,
				))
		} else {
			return errBuilder(
				models.NewAppError(ctx, "user.controller.SignupSupplier", "server.internal.error", nil, "", int(codes.Internal), err),
			)
		}
	}

	ar.Success()

	return &pb.SupplierCreateResponse{Response: &pb.SupplierCreateResponse_Data{}}, nil
}
