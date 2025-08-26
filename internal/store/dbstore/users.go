package dbstore

import (
	"encoding/json"
	"fmt"

	usersPb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/jackc/pgx/v5"
)

const SelectUserStatment = `
	SELECT
		id,
		username,
		first_name,
		last_name,
		email,
		user_type,
		image,
		image_metadata,
		membership,
		is_email_verified,
		password,
		auth_data,
		auth_service,
		roles,
		props,
		notify_props,
		last_password_update,
		last_picture_update,
		failed_attempts,
		locale,
		is_mfa_active,
		mfa_secret,
		last_activity_at,
		last_login,
		created_at,
		updated_at,
		deleted_at
	FROM users
`

func (ds *DBStore) UsersGetByEmail(ctx *models.Context, email string) (*usersPb.User, *store.DBError) {
	path := "users.store.UserGetByEmail"
	row := ds.db.QueryRow(ctx.Context, fmt.Sprintf("%s %s", SelectUserStatment, "WHERE email = $1"), email)

	return ds.scanUser(ctx, row, path)
}

// scanUser scans the whole user row given by a pgx.Row result
func (ds *DBStore) scanUser(ctx *models.Context, row pgx.Row, path string) (*usersPb.User, *store.DBError) {
	user := &usersPb.User{}

	var (
		imageMetadataBytes []byte
		propsBytes         []byte
		notifyPropsBytes   []byte
		roles              []string
	)

	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.UserType,
		&user.Image,
		&imageMetadataBytes,
		&user.Membership,
		&user.IsEmailVerified,
		&user.Password,
		&user.AuthData,
		&user.AuthService,
		&roles,
		&propsBytes,
		&notifyPropsBytes,
		&user.LastPasswordUpdate,
		&user.LastPictureUpdate,
		&user.FailedAttempts,
		&user.Locale,
		&user.MfaActive,
		&user.MfaSecret,
		&user.LastActivityAt,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return nil, store.HandleDBError(ctx, err, path, nil)
	}

	user.Roles = roles

	if len(imageMetadataBytes) > 0 {
		var meta usersPb.UserImageMetadata
		if err := json.Unmarshal(imageMetadataBytes, &meta); err != nil {
			return nil, store.JSONUnmarshalError(err, path, "an error occurred while trying to unmarshal User.image_metadat")
		}
		user.ImageMetadata = &meta
	}

	if len(propsBytes) > 0 {
		if err := json.Unmarshal(propsBytes, &user.Props); err != nil {
			return nil, store.JSONUnmarshalError(err, path, "an error occurred while trying to unmarshal User.props")
		}
	}

	// unmarshal notify_props
	if len(notifyPropsBytes) > 0 {
		if err := json.Unmarshal(notifyPropsBytes, &user.NotifyProps); err != nil {
			return nil, store.JSONUnmarshalError(err, path, "an error occurred while trying to unmarshal User.notify_props")
		}
	}

	return user, nil
}
