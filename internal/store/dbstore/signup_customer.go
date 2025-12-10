package dbstore

import (
	"encoding/json"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/utils"
	intModels "github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/jackc/pgx/v5"
)

func (ds *DBStore) SignupCustomer(ctx *models.Context, c *pb.User, token *utils.Token) *models.DBError {
	path := "user.store.SignupCustomer"
	tr, err := ds.db.BeginTx(ctx.Context, pgx.TxOptions{})
	if err != nil {
		return models.StartTransactionError(err, path)
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
	    image,
	    image_metadata,
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
	     $11, $12, $13, $14, $15, $16, $17, $18, $19
	  )
	`

	var props any
	var notifyProps any
	var imageMetadata any

	if len(c.GetProps()) > 0 {
		props, err = json.Marshal(c.GetProps())
		if err != nil {
			return models.JSONMarshalError(err, path, "an error occurred while trying to encode User.props")
		}
	}

	if len(c.GetNotifyProps()) > 0 {
		notifyProps, err = json.Marshal(c.GetNotifyProps())
		if err != nil {
			return models.JSONMarshalError(err, path, "an error occurred while trying to encode User.notify_props")
		}
	}

	if c.ImageMetadata != nil {
		imageMetadata, err = json.Marshal(c.GetImageMetadata())
		if err != nil {
			return models.JSONMarshalError(err, path, "an error occurred while trying to encode User.image_metadata")
		}
	}

	args := []any{
		c.GetId(),
		c.GetUsername(),
		c.GetFirstName(),
		c.GetLastName(),
		c.GetEmail(),
		c.GetUserType(),
		c.GetMembership(),
		c.GetImage(),
		imageMetadata,
		c.GetIsEmailVerified(),
		c.GetPassword(),
		c.GetAuthData(),
		c.GetAuthService(),
		c.GetRoles(),
		props,
		notifyProps,
		c.GetLocale(),
		c.GetMfaActive(),
		c.GetCreatedAt(),
	}

	_, err = tr.Exec(ctx.Context, stmt, args...)
	if err != nil {
		return models.HandleDBError(ctx, err, path, tr)
	}

	stmt = `
	  INSERT INTO tokens(id, user_id, token, type, created_at, expires_at) VALUES($1, $2, $3, $4, $5, $6)
	`

	args = []any{
		token.ID,
		c.GetId(),
		string(token.Hash),
		string(intModels.TokenTypeEmailConfirmation),
		utils.TimeGetMillis(),
		utils.TimeGetMillisFromTime(token.Expiry),
	}

	_, err = tr.Exec(ctx.Context, stmt, args...)
	if err != nil {
		return models.HandleDBError(ctx, err, path, tr)
	}

	if err := tr.Commit(ctx.Context); err != nil {
		return models.CommitTransactionError(err, path)
	}
	return nil
}
