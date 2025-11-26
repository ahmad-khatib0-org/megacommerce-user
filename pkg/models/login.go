package models

import (
	"fmt"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/utils"
	"google.golang.org/grpc/codes"
)

type OAuthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorHint        string `json:"error_hint,omitempty"`
	ErrorDebug       string `json:"error_debug,omitempty"`
}

func LoginRequestIsValid(ctx *models.Context, req *pb.LoginRequest) *models.AppError {
	email := req.GetEmail()
	password := req.GetPassword()
	challenge := req.GetLoginChallenge()

	path := "users.models.LoginRequestIsValid"
	if !utils.IsValidEmail(email) {
		return models.NewAppError(ctx, path, "email.invalid", nil, fmt.Sprintf("invalid email=%s", email), int(codes.InvalidArgument), &models.AppErrorErrorsArgs{ErrorsInternal: map[string]*models.AppErrorError{"email": {ID: "email.invalid"}}})
	}

	if len(password) < UserPasswordMinLength {
		params := map[string]any{"Min": UserPasswordMinLength}
		return models.NewAppError(ctx, path, "password.min_length", params, "", int(codes.InvalidArgument), &models.AppErrorErrorsArgs{ErrorsInternal: map[string]*models.AppErrorError{"password": {ID: "password.min_length", Params: params}}})
	}

	if len(password) > UserPasswordMaxLength {
		params := map[string]any{"Max": UserPasswordMaxLength}
		return models.NewAppError(ctx, path, "password.max_length", params, "", int(codes.InvalidArgument), &models.AppErrorErrorsArgs{ErrorsInternal: map[string]*models.AppErrorError{"password": {ID: "password.min_length", Params: params}}})
	}

	if len(challenge) == 0 {
		return models.NewAppError(ctx, path, "oauth.login_challenge.missing", nil, "", int(codes.InvalidArgument), nil)
	}

	return nil
}

func GetOAuthRequestErrMsg(lang, code, desc string) string {
	tr := func(id string) string {
		return models.Tr(lang, id, nil)
	}
	switch code {
	case "invalid_request":
		if contains(desc, "redirect_uri") {
			return tr("oauth.invalid_request.redirect_uri")
		}
		return tr("oauth.invalid_request.general")
	case "access_denied":
		return tr("oauth.access_denied.user")
	case "unauthorized_client":
		return tr("oauth.unauthorized_client")
	case "unsupported_response_type":
		return tr("oauth.unsupported_response_type")
	case "invalid_scope":
		return tr("oauth.invalid_scope")
	case "server_error":
		return tr("oauth.server_error.internal")
	case "temporarily_unavailable":
		return tr("oauth.temporarily_unavailable")
	default:
		return tr("oauth.unknown_error")
	}
}

func GetOAuthRequestErrMsgID(lang, code, desc string) string {
	switch code {
	case "invalid_request":
		if contains(desc, "redirect_uri") {
			return "oauth.invalid_request.redirect_uri"
		}
		return "oauth.invalid_request.general"
	case "access_denied":
		return "oauth.access_denied.user"
	case "unauthorized_client":
		return "oauth.unauthorized_client"
	case "unsupported_response_type":
		return "oauth.unsupported_response_type"
	case "invalid_scope":
		return "oauth.invalid_scope"
	case "server_error":
		return "oauth.server_error.internal"
	case "temporarily_unavailable":
		return "oauth.temporarily_unavailable"
	default:
		return "oauth.unknown_error"
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
