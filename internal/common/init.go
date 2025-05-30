package common

import (
	"fmt"
	"net"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (cc *CommonClient) initCommonClient() error {
	target := fmt.Sprintf("%s:%d", cc.cfg.Service.GrpcHost, cc.cfg.Service.GrpcPort)

	if _, err := net.ResolveTCPAddr("tcp", target); err != nil {
		return fmt.Errorf("invalid grpc url %s, : %v ", target, err)
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	fmt.Println("user service is listening on common grpc service on: " + target)

	cc.client = pb.NewCommonServiceClient(conn)
	cc.conn = conn

	return nil
}
