package server

import (
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/common"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"google.golang.org/grpc"
)

type App struct {
	conn   *grpc.ClientConn
	client pb.CommonServiceClient
	done   chan *utils.AppError
}

func RunServer(c *models.Config) *utils.AppError {
	com, err := common.NewCommonClient(c)

	app := &App{
		conn: com.Conn(),
		done: make(chan *utils.AppError, 1),
	}

	if err != nil {
		app.done <- err
	}

	_, err = com.ConfigGet()
	if err != nil {
		app.done <- err
	}

	err = <-app.done
	if err != nil {
		// TODO: cleanup things
	}

	return err
}
