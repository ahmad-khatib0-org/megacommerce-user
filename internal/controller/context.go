package controller

import (
	"context"

	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"google.golang.org/grpc/codes"
)

// TODO: enhance tracking by E,g adding and audit
func getContext(ctx context.Context) (*models.Context, *models.AppError) {
	c, ok := ctx.Value(ContextKeyMetadata).(*models.Context)
	if !ok {
		return nil, &models.AppError{
			Ctx:           c,
			Id:            "server.internal.error",
			DetailedError: "failed to get the context from the incoming request",
			Where:         "user.controller.getContext",
			StatusCode:    int(codes.Internal),
		}
	}

	return c, nil
}
