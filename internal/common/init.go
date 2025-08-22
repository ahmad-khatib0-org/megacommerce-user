package common

import (
	"context"
	"fmt"
	"net"
	"time"

	com "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (cc *CommonClient) initCommonClient() *models.InternalError {
	target := cc.cfg.Service.CommonServiceGrpcURL

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

	fmt.Println("user service connected to common service at: " + target)
	cc.client = com.NewCommonServiceClient(conn)
	cc.conn = conn

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = cc.client.Ping(ctx, &com.PingRequest{})
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
