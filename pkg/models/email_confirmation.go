package models

import (
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/utils"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
)

func EmailConfirmationIsValid(ctx *models.Context, req *pb.EmailConfirmationRequest) *models.AppError {
	if !utils.IsValidEmail(req.GetEmail()) {
		return errorBuilder(ctx, "email_confirm.email.error", nil)
	}

	if req.GetToken() == "" {
		return errorBuilder(ctx, "email_confirm.token.error", nil)
	}

	if _, err := ulid.ParseStrict(req.GetTokenId()); err != nil {
		return errorBuilder(ctx, "email_confirm.token_id.error", err)
	}

	return nil
}

func errorBuilder(ctx *models.Context, id string, err error) *models.AppError {
	return models.NewAppError(ctx, "users.models.", id, nil, "", int(codes.InvalidArgument), &models.AppErrorErrorsArgs{Err: err})
}
