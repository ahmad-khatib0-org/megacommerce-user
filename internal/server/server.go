package server

import (
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"google.golang.org/grpc"
)

type App struct {
	conn   *grpc.ClientConn
	client pb.CommonServiceClient
	done   chan error
}

func RunServer(c *models.Config) error {
	// router, err := router.InitGrpcServer(c)
	// if err != nil {
	// 	return err
	// }
	//
	// app := &App{
	// 	conn: router.Conn(),
	// 	done: make(chan error, 1),
	// }
	//
	// err = <-app.done
	// return err

	return nil
}
