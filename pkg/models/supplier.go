package models

import (
	"fmt"
	"unicode/utf8"

	common "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/user/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"google.golang.org/grpc/codes"
)

func IsValidSignupSupplierRequest(ctx *Context, s *pb.SupplierCreateRequest, passCfg *common.ConfigPassword) *AppError {
	un := s.GetUsername()
	email := s.GetEmail()
	fn := s.GetFirstName()
	ln := s.GetLastName()
	pass := s.GetPassword()
	// mem := s.GetMembership()

	if un == "" || utf8.RuneCountInString(un) > UserNameMaxLength || utf8.RuneCountInString(un) < UserNameMinLength {
		return InvalidSupplierErrorBuilder(ctx, "username", un)
	}

	if email == "" || !utils.IsValidEmail(email) {
		return InvalidSupplierErrorBuilder(ctx, "email", email)
	}

	if fn != "" {
		if utf8.RuneCountInString(fn) > UserFirstNameMaxRunes || utf8.RuneCountInString(fn) < UserFirstNameMinRunes {
			return InvalidSupplierErrorBuilder(ctx, "first_name", fn)
		}
	}

	if ln != "" {
		if utf8.RuneCountInString(ln) > UserLastNameMaxRunes || utf8.RuneCountInString(fn) < UserLastNameMinRunes {
			return InvalidSupplierErrorBuilder(ctx, "first_name", fn)
		}
	}

	if err := utils.IsValidPassword(pass, passCfg, ""); err != nil {
		return NewAppError(ctx, "user.models.SupplierCreateRequest", err.Id, nil, fmt.Sprintf("invalid password %s ", pass), int(codes.Internal))
	}

	return nil
}

func InvalidSupplierErrorBuilder(ctx *Context, fieldName string, fieldValue any) *AppError {
	where := "user.models.SupplierCreateRequest.IsValidSignupSupplierRequest"
	id := fmt.Sprintf("user.create.%s.error", fieldName)
	details := fmt.Sprintf(" %s=%v ", fieldName, fieldValue)
	return NewAppError(ctx, where, id, nil, details, int(codes.Internal))
}
