package router

import (
	"fmt"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CommonClient struct {
	conn   *grpc.ClientConn
	client pb.CommonServiceClient
}

func InitGrpcServer(c *models.Config) (*CommonClient, error) {
	target := fmt.Sprintf("%s:%d", c.Service.GrpcHost, c.Service.GrpcPort)
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	fmt.Println("common grpc service is listening on: " + target)

	client := pb.NewCommonServiceClient(conn)
	cc := &CommonClient{conn: conn, client: client}
	cc.initRouter()
	return cc, nil
}

func (c *CommonClient) initRouter() {
}

func (c *CommonClient) Close() {
	c.conn.Close()
}

func (c *CommonClient) Conn() *grpc.ClientConn {
	return c.conn
}
