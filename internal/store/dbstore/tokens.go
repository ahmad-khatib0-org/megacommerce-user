package dbstore

import (
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

func (ds *DBStore) MarkEmailAsConfirmed(ctx *models.Context, tokenId string) *store.DBError {
	stmt := `UPDATE tokens SET used = TRUE WHERE id = $1`
	_, err := ds.db.Exec(ctx.Context, stmt, tokenId)

	return store.HandleDBError(ctx, err, "users.store.MarkEmailAsConfirmed", nil)
}

func (ds *DBStore) TokensGet(ctx *models.Context, tokenId string) (*pb.Token, *store.DBError) {
	stmt := `SELECT id, token, type, used, created_at, expires_at FROM tokens WHERE id = $1`

	var t pb.Token
	err := ds.db.QueryRow(ctx.Context, stmt, tokenId).Scan(
		&t.Id,
		&t.Token,
		&t.Type,
		&t.Used,
		&t.CreatedAt,
		&t.ExpiresAt,
	)
	if err != nil {
		return nil, store.HandleDBError(ctx, err, "users.store.TokensGet", nil)
	}

	return &t, nil
}
