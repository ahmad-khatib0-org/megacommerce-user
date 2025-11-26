package dbstore

import (
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/utils"
	intModels "github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

func (ds *DBStore) MarkEmailAsConfirmed(ctx *models.Context, tokenID string) *models.DBError {
	stmt := `UPDATE tokens SET used = TRUE WHERE id = $1`
	_, err := ds.db.Exec(ctx.Context, stmt, tokenID)

	return models.HandleDBError(ctx, err, "users.store.MarkEmailAsConfirmed", nil)
}

func (ds *DBStore) TokensGet(ctx *models.Context, tokenID string) (*pb.Token, *models.DBError) {
	stmt := `SELECT id, user_id, token, type, used, created_at, expires_at FROM tokens WHERE id = $1`

	var t pb.Token
	err := ds.db.QueryRow(ctx.Context, stmt, tokenID).Scan(
		&t.Id,
		&t.UserId,
		&t.Token,
		&t.Type,
		&t.Used,
		&t.CreatedAt,
		&t.ExpiresAt,
	)
	if err != nil {
		return nil, models.HandleDBError(ctx, err, "users.store.TokensGet", nil)
	}

	return &t, nil
}

func (ds *DBStore) TokensGetAllByUserID(ctx *models.Context, userID string) ([]*pb.Token, *models.DBError) {
	stmt := `SELECT id, user_id, token, type, used, created_at, expires_at FROM tokens WHERE user_id = $1`

	rows, err := ds.db.Query(ctx.Context, stmt, userID)
	if err != nil {
		return nil, models.HandleDBError(ctx, err, "users.store.TokensGetAllByUserID", nil)
	}
	defer rows.Close()

	result := []*pb.Token{}
	for rows.Next() {
		t := &pb.Token{}
		if err := rows.Scan(t.Id, t.UserId, t.Token, t.Type, t.Used, t.CreatedAt, t.ExpiresAt); err != nil {
			return nil, models.HandleDBError(ctx, err, "users.store.TokensGetAllByUserID", nil)
		}
		result = append(result, t)
	}

	return result, nil
}

func (ds *DBStore) TokensAdd(ctx *models.Context, userID string, token *utils.Token, tokenType intModels.TokenType, path string) *models.DBError {
	stmt := `
	  INSERT INTO tokens(id, user_id, token, type, created_at, expires_at) VALUES($1, $2, $3, $4, $5, $6)
	`

	args := []any{
		token.ID,
		userID,
		string(token.Hash),
		string(tokenType),
		utils.TimeGetMillis(),
		utils.TimeGetMillisFromTime(token.Expiry),
	}

	_, err := ds.db.Exec(ctx.Context, stmt, args...)
	if err != nil {
		return models.HandleDBError(ctx, err, path, nil)
	}
	return nil
}

// TokensDeleteAllPasswordResetByUserID returns the number of deleted rows(or 0), error
func (ds *DBStore) TokensDeleteAllPasswordResetByUserID(ctx *models.Context, userID string) (int64, *models.DBError) {
	stmt := `DELETE FROM tokens WHERE user_id = $1 AND type = $2`
	res, err := ds.db.Exec(ctx.Context, stmt, userID, string(intModels.TokenTypePasswordReset))
	if err != nil {
		return 0, models.HandleDBError(ctx, err, "users.store.TokensDeleteAllPasswordResetByUserID", nil)
	}

	return res.RowsAffected(), nil
}
