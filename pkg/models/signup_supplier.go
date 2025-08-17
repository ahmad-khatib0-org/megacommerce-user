package models

import (
	"fmt"
	"strings"
	"unicode/utf8"

	common "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	user "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"google.golang.org/grpc/codes"
)

func SignupSupplierRequestSanitize(s *user.SupplierCreateRequest) *user.SupplierCreateRequest {
	return &user.SupplierCreateRequest{
		Username:  utils.SanitizeUnicode(s.GetUsername()),
		Email:     strings.ToLower(s.GetEmail()),
		FirstName: utils.SanitizeUnicode(s.GetFirstName()),
		LastName:  utils.SanitizeUnicode(s.GetLastName()),
		Password:  s.GetPassword(),
	}
}

func SignupSupplierRequestIsValid(ctx *Context, s *user.SupplierCreateRequest, passCfg *common.ConfigPassword) *AppError {
	un := s.GetUsername()
	email := s.GetEmail()
	fn := s.GetFirstName()
	ln := s.GetLastName()
	pass := s.GetPassword()

	if un == "" || utf8.RuneCountInString(un) > UserNameMaxLength || utf8.RuneCountInString(un) < UserNameMinLength {
		return signupSupplierRequestErrorBuilder(ctx, "username", un, map[string]any{"Min": UserNameMinLength, "Max": UserNameMaxLength})
	}

	if !utils.IsValidUsernameChars(un) {
		return signupSupplierRequestErrorBuilder(ctx, "username.valid", un, nil)
	}

	if email == "" || !utils.IsValidEmail(email) {
		return signupSupplierRequestErrorBuilder(ctx, "email", email, nil)
	}

	if utf8.RuneCountInString(fn) > UserFirstNameMaxRunes || utf8.RuneCountInString(fn) < UserFirstNameMinRunes {
		return signupSupplierRequestErrorBuilder(ctx, "first_name", fn, map[string]any{"Min": UserFirstNameMinRunes, "Max": UserFirstNameMaxRunes})
	}

	if utf8.RuneCountInString(ln) > UserLastNameMaxRunes || utf8.RuneCountInString(ln) < UserLastNameMinRunes {
		return signupSupplierRequestErrorBuilder(ctx, "last_name", ln, map[string]any{"Min": UserLastNameMinRunes, "Max": UserLastNameMaxRunes})
	}

	if err := utils.IsValidPassword(pass, passCfg, ""); err != nil {
		return NewAppError(ctx, "user.models.SupplierCreateRequest.password", err.ID, err.Params, fmt.Sprintf("invalid password %s ", pass), int(codes.InvalidArgument), err)
	}

	return nil
}

func signupSupplierRequestErrorBuilder(ctx *Context, fieldName string, fieldValue any, params map[string]any) *AppError {
	where := "user.models.SupplierCreateRequest.SignupSupplierRequestIsValid"
	id := fmt.Sprintf("user.create.%s.error", fieldName)
	details := fmt.Sprintf(" %s=%v ", fieldName, fieldValue)
	return NewAppError(ctx, where, id, params, details, int(codes.InvalidArgument), nil)
}

func SignupSupplierRequestAuditable(s *user.SupplierCreateRequest) map[string]string {
	return map[string]string{
		"username":   s.GetUsername(),
		"email":      s.GetEmail(),
		"first_name": s.GetFirstName(),
		"last_name":  s.GetLastName(),
		"membership": s.GetMembership(),
	}
}

// SignupSupplierRequestPreSave convert SupplierCreateRequest to User
// and populate the necessary fields with values to be stored in db
func SignupSupplierRequestPreSave(ctx *Context, s *user.User) (*user.User, *AppError) {
	pass, err := utils.PasswordHash(s.GetPassword())
	if err != nil {
		return nil, NewAppError(ctx, "user.models.SignupSupplierRequestPreSave", ErrMsgInternal, nil, "failed to generate password %v", int(codes.Internal), err)
	}

	u := &user.User{
		Id:                 utils.NewIDPointer(),
		Username:           utils.NewPointer(s.GetUsername()),
		FirstName:          utils.NewPointer(s.GetFirstName()),
		LastName:           utils.NewPointer(s.GetLastName()),
		Email:              utils.NewPointer(s.GetEmail()),
		UserType:           utils.NewPointer(string(UserTypeSupplier)),
		Membership:         utils.NewPointer(s.GetMembership()),
		IsEmailVerified:    utils.NewPointer(s.GetIsEmailVerified()),
		Password:           utils.NewPointer(pass),
		AuthData:           utils.NewPointer(s.GetAuthData()),
		AuthService:        utils.NewPointer(s.GetAuthService()),
		Roles:              s.GetRoles(),
		Props:              s.GetProps(),
		NotifyProps:        s.GetNotifyProps(),
		Locale:             utils.NewPointer(s.GetLocale()),
		MfaActive:          utils.NewPointer(s.GetMfaActive()),
		LastPasswordUpdate: nil,
		LastPictureUpdate:  nil,
		FailedAttempts:     nil,
		MfaSecret:          nil,
		LastActivityAt:     nil,
		LastLogin:          nil,
		UpdatedAt:          nil,
		DeletedAt:          nil,
		CreatedAt:          utils.NewPointer(utils.TimeGetMillis()),
	}

	return u, nil
}

func SignupSupplierRequestResultState(s *user.User) map[string]any {
	return map[string]any{
		"id":                   s.GetId(),
		"username":             s.GetUsername(),
		"first_name":           s.GetFirstName(),
		"last_name":            s.GetLastName(),
		"email":                s.GetEmail(),
		"user_type":            s.GetUserType(),
		"membership":           s.GetMembership(),
		"is_email_verified":    s.GetIsEmailVerified(),
		"password":             s.GetPassword(),
		"auth_data":            s.GetAuthData(),
		"auth_service":         s.GetAuthService(),
		"roles":                s.GetRoles(),
		"props":                s.GetProps(),
		"notify_props":         s.GetNotifyProps(),
		"locale":               s.GetLocale(),
		"mfa_active":           s.GetMfaActive(),
		"last_password_update": nil,
		"last_picture_update":  nil,
		"failed_attempts":      nil,
		"mfa_secret":           nil,
		"last_activity_at":     nil,
		"last_login":           nil,
		"updated_at":           nil,
		"deleted_at":           nil,
		"created_at":           s.GetCreatedAt(),
	}
}
