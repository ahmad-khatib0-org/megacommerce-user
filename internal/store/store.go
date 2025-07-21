package store

import (
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
)

type UsersStore interface {
	SignupSupplier(ctx *models.Context, s *pb.User, token *utils.Token) *DBError
}
