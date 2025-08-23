package controller

import (
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"google.golang.org/grpc/codes"
)

func getGrpcErrMsg(lang string, code codes.Code) string {
	switch code {
	case codes.Canceled:
		return models.Tr(lang, "error.canceled", nil)
	case codes.Unknown:
		return models.Tr(lang, "error.unknown", nil)
	case codes.InvalidArgument:
		return models.Tr(lang, "error.invalid_argument", nil)
	case codes.DeadlineExceeded:
		return models.Tr(lang, "error.deadline_exceeded", nil)
	case codes.NotFound:
		return models.Tr(lang, "error.not_found", nil)
	case codes.AlreadyExists:
		return models.Tr(lang, "error.already_exists", nil)
	case codes.PermissionDenied:
		return models.Tr(lang, "error.permission_denied", nil)
	case codes.ResourceExhausted:
		return models.Tr(lang, "error.resource_exhausted", nil)
	case codes.FailedPrecondition:
		return models.Tr(lang, "error.failed_precondition", nil)
	case codes.Aborted:
		return models.Tr(lang, "error.aborted", nil)
	case codes.OutOfRange:
		return models.Tr(lang, "error.out_of_range", nil)
	case codes.Unimplemented:
		return models.Tr(lang, "error.unimplemented", nil)
	case codes.Internal:
		return models.Tr(lang, "error.internal", nil)
	case codes.Unavailable:
		return models.Tr(lang, "error.unavailable", nil)
	case codes.DataLoss:
		return models.Tr(lang, "error.data_loss", nil)
	case codes.Unauthenticated:
		return models.Tr(lang, "error.unauthenticated", nil)
	}

	return models.Tr(lang, "error.unknown", nil)
}
