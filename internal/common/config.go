package common

import (
	"context"
	"fmt"
	"time"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"google.golang.org/grpc/status"
)

func (cc *CommonClient) ConfigGet() (*pb.Config, *utils.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := cc.client.ConfigGet(ctx, &pb.ConfigGetRequest{})
	if err != nil {
		sc, ok := status.FromError(err)
		if ok {
			return nil, utils.NewAppError("user.common.ConfigGet", "failed_get_common_config", nil, fmt.Sprintf("failed to get configurations %v", err), int(sc.Code()))
		}
	}

	switch res := res.Response.(type) {
	case *pb.ConfigGetResponse_Data:
		return res.Data, nil
	case *pb.ConfigGetResponse_Error:
		err := utils.AppErrorFromProto(res.Error)
		return nil, err
	}

	return nil, nil
}
