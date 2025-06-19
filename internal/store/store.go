package store

import (
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/user/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

type UsersStore interface {
	SignupSupplier(ctx *models.Context, s *pb.User) *DBError
}
