package controller

import (
	"context"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/user/v1"
)

func (c Controller) CreateSupplier(ctx context.Context, s *pb.SupplierCreateRequest) (*pb.SupplierCreateResponse, error) {
	return &pb.SupplierCreateResponse{}, nil
}
