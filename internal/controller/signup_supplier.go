package controller

import (
	"context"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/user/v1"
)

func (c *Controller) CreateSupplier(context context.Context, s *pb.SupplierCreateRequest) (*pb.SupplierCreateResponse, error) {
	ctx, err := getContext(context)
	if err != nil {
		return nil, err
	}

	c.log.InfoStruct("incoming context is ", ctx)

	return &pb.SupplierCreateResponse{}, nil
}
