package store

import (
	"context"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/user/v1"
)

type UsersStore interface {
	SignupSupplier(ctx context.Context, s *pb.SupplierCreateRequest) error
}
