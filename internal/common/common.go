package common

import (
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"google.golang.org/grpc"
)

type CommonClient struct {
	pb.UnimplementedCommonServiceServer
	cfg    *models.Config
	conn   *grpc.ClientConn
	client pb.CommonServiceClient
}

// NewCommonClient runs the CommonService client
func NewCommonClient(config *models.Config) (*CommonClient, *models.InternalError) {
	c := &CommonClient{cfg: config}
	if err := c.initCommonClient(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CommonClient) Close() {
	c.conn.Close()
}

func (c *CommonClient) Conn() *grpc.ClientConn {
	return c.conn
}
