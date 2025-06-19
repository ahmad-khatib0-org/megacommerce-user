package dbstore

import (
	"encoding/json"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/user/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

func (ds *DBStore) SignupSupplier(ctx *models.Context, u *pb.User) *store.DBError {
	stmt := `
	  INSERT INTO users(
			id,
			username,
			first_name,
			last_name,
			email,
			user_type,
			membership,
			is_email_verified,
	    password,
			auth_data,
			auth_service,
			roles,
			props,
			notify_props,
			locale,
			is_mfa_active,
			created_at
	  ) VALUES (
	     $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
	     $11, $12, $13, $14, $15, $16, $17
	  )
	`

	var err error
	var props any
	var nontifyProps any

	if len(u.GetProps()) > 0 {
		props, err = json.Marshal(u.GetProps())
		if err != nil {
			return &store.DBError{
				ErrType: store.DBErrorTypeJsonMarshal,
				Err:     err,
				Msg:     "failed to marshal props into json",
				Path:    "user.store.SignupSupplier",
				Details: "an error occurred while trying to encode User.props",
			}
		}
	}

	if len(u.GetNotifyProps()) > 0 {
		nontifyProps, err = json.Marshal(u.GetNotifyProps())
		if err != nil {
			return &store.DBError{
				ErrType: store.DBErrorTypeJsonMarshal,
				Err:     err,
				Msg:     "failed to marshal notify_props",
				Path:    "user.store.SignupSupplier",
				Details: "an error occurred while trying to encode User.notify_props",
			}
		}
	}

	args := []any{
		u.GetId(),
		u.GetUsername(),
		u.GetFirstName(),
		u.GetLastName(),
		u.GetEmail(),
		u.GetUserType(),
		u.GetMembership(),
		u.GetIsEmailVerified(),
		u.GetPassword(),
		u.GetAuthData(),
		u.GetAuthService(),
		u.GetRoles(),
		props,
		nontifyProps,
		u.GetLocale(),
		u.GetMfaActive(),
		u.GetCreatedAt(),
	}

	_, err = ds.db.Exec(ctx.Ctx(), stmt, args...)
	if err != nil {
		return store.HandleDBError(err, "user.store.SignupSupplier")
	}

	return nil
}
