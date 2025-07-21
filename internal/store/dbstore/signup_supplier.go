package dbstore

import (
	"encoding/json"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"github.com/jackc/pgx/v5"
)

func (ds *DBStore) SignupSupplier(ctx *models.Context, u *pb.User, token *utils.Token) *store.DBError {
	path := "user.store.SignupSupplier"
	tr, err := ds.db.BeginTx(ctx.Context, pgx.TxOptions{})
	if err != nil {
		return store.StartTransactionError(err, path)
	}

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

	var props any
	var nontifyProps any

	if len(u.GetProps()) > 0 {
		props, err = json.Marshal(u.GetProps())
		if err != nil {
			return store.JsonMarshalError(err, path, "an error occurred while trying to encode User.props")
		}
	}

	if len(u.GetNotifyProps()) > 0 {
		nontifyProps, err = json.Marshal(u.GetNotifyProps())
		if err != nil {
			return store.JsonMarshalError(err, path, "an error occurred while trying to encode User.notify_props")
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

	_, err = tr.Exec(ctx.Context, stmt, args...)
	if err != nil {
		return store.HandleDBError(ctx, err, path, tr)
	}

	stmt = `
	  INSERT INTO tokens(id, token, type, created_at, expires_at) VALUES($1, $2, $3, $4, $5)
	`

	args = []any{
		token.Id,
		string(token.Hash),
		string(models.TokenTypeEmailConfirmation),
		utils.TimeGetMillis(),
		utils.TimeGetMillisFromTime(token.Expiry),
	}

	_, err = tr.Exec(ctx.Context, stmt, args...)
	if err != nil {
		return store.HandleDBError(ctx, err, path, tr)
	}

	if err := tr.Commit(ctx.Context); err != nil {
		return store.CommitTransactionError(err, path)
	}
	return nil
}
