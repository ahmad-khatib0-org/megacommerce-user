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

func (cc *CommonClient) ConfigListener(clientID string) *utils.AppError {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := cc.client.ConfigListener(ctx, &pb.ConfigListenerRequest{ClientId: clientID})
	if err != nil {
		sc, _ := status.FromError(err)
		return utils.NewAppError("user.common.ConfigListener", "failed_listen_common_config", nil, fmt.Sprintf("failed to register a listener call: %v", err), int(sc.Code()))
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			fmt.Println("config listener stream closed")
			break
		}

		switch x := res.Response.(type) {
		case *pb.ConfigListenerResponse_Data:
			fmt.Println("Config changed: ", x.Data)
		case *pb.ConfigListenerResponse_Error:
			fmt.Println("Error received: ", x.Error.Message)
		}
	}

	return nil
}
