package common

import (
	com "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/logger"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"google.golang.org/grpc"
)

type CommonClient struct {
	com.UnimplementedCommonServiceServer
	cfg    *models.Config
	conn   *grpc.ClientConn
	client com.CommonServiceClient
	log    *logger.Logger
}

type CommonArgs struct {
	Config *models.Config
	Log    *logger.Logger
}

// NewCommonClient runs the CommonService client
func NewCommonClient(ca *CommonArgs) (*CommonClient, *models.InternalError) {
	c := &CommonClient{cfg: ca.Config, log: ca.Log}
	if err := c.initCommonClient(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CommonClient) Close() error {
	return c.conn.Close()
}

func (c *CommonClient) Conn() *grpc.ClientConn {
	return c.conn
}
