// Package store define interfaces and helpers for database
package store

import (
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
)

type UsersStore interface {
	SignupSupplier(ctx *models.Context, s *pb.User, token *utils.Token) *DBError
	MarkEmailAsConfirmed(ctx *models.Context, tokenID string) *DBError
	UsersGetByEmail(ctx *models.Context, email string) (*pb.User, *DBError)
	TokensGet(ctx *models.Context, tokenID string) (*pb.Token, *DBError)
	TokensGetAllByUserID(ctx *models.Context, userID string) ([]*pb.Token, *DBError)
	TokensAdd(ctx *models.Context, userID string, token *utils.Token, tokenType models.TokenType, path string) *DBError
	// TokensDeleteAllPasswordResetByUserID returns the number of deleted rows(or 0), error
	TokensDeleteAllPasswordResetByUserID(ctx *models.Context, userID string) (int64, *DBError)
}
