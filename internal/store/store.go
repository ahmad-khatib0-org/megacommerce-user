// Package store define interfaces and helpers for database
package store

import (
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/utils"
	intModels "github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

type UsersStore interface {
	SignupSupplier(ctx *models.Context, s *pb.User, token *utils.Token) *models.DBError
	SignupCustomer(ctx *models.Context, c *pb.User, token *utils.Token) *models.DBError
	MarkEmailAsConfirmed(ctx *models.Context, tokenID string) *models.DBError
	UsersGetByEmail(ctx *models.Context, email string) (*pb.User, *models.DBError)
	TokensGet(ctx *models.Context, tokenID string) (*pb.Token, *models.DBError)
	TokensGetAllByUserID(ctx *models.Context, userID string) ([]*pb.Token, *models.DBError)
	TokensAdd(ctx *models.Context, userID string, token *utils.Token, tokenType intModels.TokenType, path string) *models.DBError
	// TokensDeleteAllPasswordResetByUserID returns the number of deleted rows(or 0), error
	TokensDeleteAllPasswordResetByUserID(ctx *models.Context, userID string) (int64, *models.DBError)
}
