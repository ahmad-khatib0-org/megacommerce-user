package controller

import (
	"context"
	"time"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	intModels "github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"google.golang.org/grpc/codes"
)

func (c *Controller) GetCustomerProfile(context context.Context, req *pb.CustomerProfileRequest) (*pb.CustomerProfileResponse, error) {
	start := time.Now()
	path := "user.controller.GetCustomerProfile"
	errBuilder := func(e *models.AppError) (*pb.CustomerProfileResponse, error) {
		return &pb.CustomerProfileResponse{Response: &pb.CustomerProfileResponse_Error{Error: models.AppErrorToProto(e)}}, nil
	}

	ctx, err := models.ContextGet(context)
	if err != nil {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordCustomerProfileGetRequest(false, duration)
		return errBuilder(err)
	}

	ar := models.AuditRecordNew(ctx, intModels.EventNameCustomerProfileGet, models.EventStatusFail)
	defer c.ProcessAudit(ar)

	// Get user ID from context session
	userID := ctx.Session.UserID
	if userID == "" {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordCustomerProfileGetRequest(false, duration)
		return errBuilder(models.NewAppError(ctx, path, "error.unauthenticated", nil, "user not authenticated", int(codes.Unauthenticated), nil))
	}

	user, dbErr := c.store.UsersGetByID(ctx, userID)
	if dbErr != nil {
		duration := time.Since(start).Seconds()
		c.metricsCollector.RecordCustomerProfileGetRequest(false, duration)
		if dbErr.ErrType == models.DBErrorTypeNoRows {
			return errBuilder(models.NewAppError(ctx, path, "error.not_found", nil, "user not found", int(codes.NotFound), nil))
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

	profile := &pb.CustomerProfile{
		Id:              user.GetId(),
		FullName:        fullName,
		Email:           user.GetEmail(),
		Username:        user.GetUsername(),
		Image:           user.GetImage(),
		UserType:        user.GetUserType(),
		IsEmailVerified: user.GetIsEmailVerified(),
		CreatedAt:       user.GetCreatedAt(),
		UpdatedAt:       user.GetUpdatedAt(),
	}

	duration := time.Since(start).Seconds()
	c.metricsCollector.RecordCustomerProfileGetRequest(true, duration)

	return &pb.CustomerProfileResponse{Response: &pb.CustomerProfileResponse_Data{Data: profile}}, nil
}
