package common

import (
	"fmt"
	"net"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

func (cc *CommonClient) initCommonClient() *utils.AppError {
	target := fmt.Sprintf("%s:%d", cc.cfg.Service.CommonServiceGrpcHost, cc.cfg.Service.CommonServiceGrpcPort)

	if _, err := net.ResolveTCPAddr("tcp", target); err != nil {
		return utils.NewAppError("user.common.initCommonClient", "invalid_grpc_url", nil, fmt.Sprintf("invalid grpc url %s, : %v ", target, err), int(codes.Internal))
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return utils.NewAppError("user.common.initCommonClient", "common_service_connect_error", nil, fmt.Sprintf("invalid grpc url %s, : %v ", target, err), int(codes.Internal))
	}

	fmt.Println("user service is listening on common grpc service on: " + target)
	cc.client = pb.NewCommonServiceClient(conn)
	cc.conn = conn

	return nil
}
