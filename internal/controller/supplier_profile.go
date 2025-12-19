package controller

import (
	"context"
	"time"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	intModels "github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"google.golang.org/grpc/codes"
)

func (c *Controller) GetSupplierProfile(context context.Context, req *pb.SupplierProfileRequest) (*pb.SupplierProfileResponse, error) {
	start := time.Now()
	path := "user.controller.GetSupplierProfile"
	errBuilder := func(e *models.AppError) (*pb.SupplierProfileResponse, error) {
		return &pb.SupplierProfileResponse{Response: &pb.SupplierProfileResponse_Error{Error: models.AppErrorToProto(e)}}, nil
	}

	ctx, err := models.ContextGet(context)
	if err != nil {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordSupplierProfileGetRequest(false, duration)
		return errBuilder(err)
	}

	ar := models.AuditRecordNew(ctx, intModels.EventNameSupplierProfileGet, models.EventStatusFail)
	defer c.ProcessAudit(ar)

	userID := ctx.Session.UserID
	if userID == "" {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordSupplierProfileGetRequest(false, duration)
		return errBuilder(models.NewAppError(ctx, path, "error.unauthenticated", nil, "user not authenticated", int(codes.Unauthenticated), nil))
	}

	user, dbErr := c.store.UsersGetByID(ctx, userID)
	if dbErr != nil {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordSupplierProfileGetRequest(false, duration)
		if dbErr.ErrType == models.DBErrorTypeNoRows {
			return errBuilder(models.NewAppError(ctx, path, "error.not_found", nil, "supplier not found", int(codes.NotFound), nil))
		}
		return errBuilder(models.NewAppError(ctx, path, models.ErrMsgInternal, nil, dbErr.Details, int(codes.Internal), &models.AppErrorErrorsArgs{Err: dbErr}))
	}

	// Combine first_name and last_name for full_name
	fullName := user.GetFirstName()
	if lastName := user.GetLastName(); lastName != "" {
		if fullName != "" {
			fullName += " " + lastName
		} else {
			fullName = lastName
		}
	}

	ar.Success()

	profile := &pb.SupplierProfile{
		Id:              user.GetId(),
		FullName:        fullName,
		Email:           user.GetEmail(),
		Username:        user.GetUsername(),
		Image:           user.GetImage(),
		UserType:        user.GetUserType(),
		Membership:      user.GetMembership(),
		IsEmailVerified: user.GetIsEmailVerified(),
		CreatedAt:       user.GetCreatedAt(),
		UpdatedAt:       user.GetUpdatedAt(),
	}

	duration := time.Since(start).Seconds()
	c.metricsCollector.RecordSupplierProfileGetRequest(true, duration)

	return &pb.SupplierProfileResponse{Response: &pb.SupplierProfileResponse_Data{Data: profile}}, nil
}
