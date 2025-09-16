package models

import (
	"fmt"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"google.golang.org/grpc/codes"
)

type OAuthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorHint        string `json:"error_hint,omitempty"`
	ErrorDebug       string `json:"error_debug,omitempty"`
}

func LoginRequestIsValid(ctx *Context, req *pb.LoginRequest) *AppError {
	email := req.GetEmail()
	password := req.GetPassword()

	path := "users.models.LoginRequestIsValid"
	if !utils.IsValidEmail(email) {
		return NewAppError(ctx, path, "email.invalid", nil, fmt.Sprintf("invalid email=%s", email), int(codes.InvalidArgument), &AppErrorErrorsArgs{ErrorsInternal: map[string]*AppErrorError{"email": {ID: "email.invalid"}}})
	}

	if len(password) < UserPasswordMinLength {
		params := map[string]any{"Min": UserPasswordMinLength}
		return NewAppError(ctx, path, "password.min_length", params, "", int(codes.InvalidArgument), &AppErrorErrorsArgs{ErrorsInternal: map[string]*AppErrorError{"password": {ID: "password.min_length", Params: params}}})
	}

	if len(password) > UserPasswordMaxLength {
		params := map[string]any{"Max": UserPasswordMaxLength}
		return NewAppError(ctx, path, "password.max_length", params, "", int(codes.InvalidArgument), &AppErrorErrorsArgs{ErrorsInternal: map[string]*AppErrorError{"password": {ID: "password.min_length", Params: params}}})
	}

	return nil
}

func GetOAuthRequestErrMsg(lang, code, desc string) string {
	tr := func(id string) string {
		return Tr(lang, id, nil)
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
