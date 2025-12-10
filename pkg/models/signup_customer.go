package models

import (
	"fmt"
	"strings"
	"unicode/utf8"

	common "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	user "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/utils"
	"google.golang.org/grpc/codes"
)

func SignupCustomerRequestSanitize(c *user.CustomerCreateRequest) *user.CustomerCreateRequest {
	return &user.CustomerCreateRequest{
		Username:  utils.SanitizeUnicode(c.GetUsername()),
		Email:     strings.ToLower(c.GetEmail()),
		FirstName: utils.SanitizeUnicode(c.GetFirstName()),
		LastName:  utils.SanitizeUnicode(c.GetLastName()),
		Password:  c.GetPassword(),
		Image:     c.Image,
	}
}

func SignupCustomerRequestIsValid(ctx *models.Context, c *user.CustomerCreateRequest, passCfg *common.ConfigPassword) *models.AppError {
	un := c.GetUsername()
	email := c.GetEmail()
	fn := c.GetFirstName()
	ln := c.GetLastName()
	pass := c.GetPassword()

	if un == "" || utf8.RuneCountInString(un) > UserNameMaxLength || utf8.RuneCountInString(un) < UserNameMinLength {
		return signupCustomerRequestErrorBuilder(ctx, "username", un, map[string]any{"Min": UserNameMinLength, "Max": UserNameMaxLength})
	}

	if !utils.IsValidUsernameChars(un) {
		return signupCustomerRequestErrorBuilder(ctx, "username.valid", un, nil)
	}

	if email == "" || !utils.IsValidEmail(email) {
		return signupCustomerRequestErrorBuilder(ctx, "email", email, nil)
	}

	if utf8.RuneCountInString(fn) > UserFirstNameMaxRunes || utf8.RuneCountInString(fn) < UserFirstNameMinRunes {
		return signupCustomerRequestErrorBuilder(ctx, "first_name", fn, map[string]any{"Min": UserFirstNameMinRunes, "Max": UserFirstNameMaxRunes})
	}

	if utf8.RuneCountInString(ln) > UserLastNameMaxRunes || utf8.RuneCountInString(ln) < UserLastNameMinRunes {
		return signupCustomerRequestErrorBuilder(ctx, "last_name", ln, map[string]any{"Min": UserLastNameMinRunes, "Max": UserLastNameMaxRunes})
	}

	if err := utils.IsValidPassword(pass, passCfg, ""); err != nil {
		errors := &models.AppErrorErrorsArgs{Err: err, ErrorsInternal: map[string]*models.AppErrorError{"password": {ID: err.ID, Params: err.Params}}}
		e := models.NewAppError(ctx, "user.models.CustomerCreateRequest.SignupCustomerRequestIsValid", err.ID, err.Params, fmt.Sprintf("invalid password %s ", pass), int(codes.InvalidArgument), errors)
		return e
	}

	return nil
}

func signupCustomerRequestErrorBuilder(ctx *models.Context, fieldName string, fieldValue any, params map[string]any) *models.AppError {
	where := "user.models.CustomerCreateRequest.SignupCustomerRequestIsValid"
	id := fmt.Sprintf("user.create.%s.error", fieldName)
	details := fmt.Sprintf(" %s=%v ", fieldName, fieldValue)
	errors := &models.AppErrorErrorsArgs{ErrorsInternal: map[string]*models.AppErrorError{fieldName: {ID: id, Params: params}}}
	err := models.NewAppError(ctx, where, id, params, details, int(codes.InvalidArgument), errors)
	return err
}

func SignupCustomerRequestAuditable(c *user.CustomerCreateRequest) map[string]string {
	return map[string]string{
		"username":   c.GetUsername(),
		"email":      c.GetEmail(),
		"first_name": c.GetFirstName(),
		"last_name":  c.GetLastName(),
	}
}

// SignupCustomerRequestPreSave convert CustomerCreateRequest to User
// and populate the necessary fields with values to be stored in db
func SignupCustomerRequestPreSave(ctx *models.Context, c *user.User) (*user.User, *models.AppError) {
	pass, err := utils.PasswordHash(c.GetPassword())
	if err != nil {
		return nil, models.NewAppError(ctx,
			"user.models.SignupCustomerRequestPreSave", models.ErrMsgInternal, nil,
			fmt.Sprintf("failed to generate password %v", err), int(codes.Internal),
			&models.AppErrorErrorsArgs{Err: err},
		)
	}

	u := &user.User{
		Id:                 utils.NewIDPointer(),
		Username:           utils.NewPointer(c.GetUsername()),
		FirstName:          utils.NewPointer(c.GetFirstName()),
		LastName:           utils.NewPointer(c.GetLastName()),
		Email:              utils.NewPointer(c.GetEmail()),
		UserType:           utils.NewPointer(string(UserTypeCustomer)),
		Membership:         utils.NewPointer("free"),
		IsEmailVerified:    utils.NewPointer(false),
		Password:           utils.NewPointer(pass),
		AuthData:           utils.NewPointer(c.GetAuthData()),
		AuthService:        utils.NewPointer(c.GetAuthService()),
		Roles:              []string{string(models.RoleIDCustomer)},
		Props:              c.GetProps(),
		NotifyProps:        c.GetNotifyProps(),
		Locale:             utils.NewPointer(c.GetLocale()),
		MfaActive:          utils.NewPointer(false),
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

func SignupCustomerRequestResultState(c *user.User) map[string]any {
	return map[string]any{
		"id":                   c.GetId(),
		"username":             c.GetUsername(),
		"first_name":           c.GetFirstName(),
		"last_name":            c.GetLastName(),
		"email":                c.GetEmail(),
		"user_type":            c.GetUserType(),
		"membership":           c.GetMembership(),
		"is_email_verified":    c.GetIsEmailVerified(),
		"password":             c.GetPassword(),
		"auth_data":            c.GetAuthData(),
		"auth_service":         c.GetAuthService(),
		"roles":                c.GetRoles(),
		"props":                c.GetProps(),
		"notify_props":         c.GetNotifyProps(),
		"locale":               c.GetLocale(),
		"mfa_active":           c.GetMfaActive(),
		"last_password_update": nil,
		"last_picture_update":  nil,
		"failed_attempts":      nil,
		"mfa_secret":           nil,
		"last_activity_at":     nil,
		"last_login":           nil,
		"updated_at":           nil,
		"deleted_at":           nil,
		"created_at":           c.GetCreatedAt(),
	}
}
