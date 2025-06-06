package common

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (cc *CommonClient) initCommonClient() *models.InternalError {
	target := fmt.Sprintf("%s:%d", cc.cfg.Service.CommonServiceGrpcHost, cc.cfg.Service.CommonServiceGrpcPort)

	if _, err := net.ResolveTCPAddr("tcp", target); err != nil {
		return &models.InternalError{
			Temp: false,
			Err:  err,
			Msg:  "failed to init common service client, invalid grpc url",
			Path: "user.config.initCommonClient",
		}
	}

	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return &models.InternalError{
			Temp: false,
			Err:  err,
			Msg:  "failed to connect to the shared common service",
			Path: "user.config.initCommonClient",
		}
	}

	fmt.Println("user service is listening on common grpc service on: " + target)
	cc.client = pb.NewCommonServiceClient(conn)
	cc.conn = conn

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = cc.client.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		return &models.InternalError{
			Temp: false,
			Err:  err,
			Msg:  "failed to ping the common service",
			Path: "user.config.initCommonClient",
		}
	}

	return nil
}
