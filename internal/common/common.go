// Package common connect to the common service for trans, config ...
package common

import (
	com "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/logger"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	intModels "github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"google.golang.org/grpc"
)

type CommonClient struct {
	com.UnimplementedCommonServiceServer
	cfg    *intModels.Config
	conn   *grpc.ClientConn
	client com.CommonServiceClient
	log    *logger.Logger
}

type CommonArgs struct {
	Config *intModels.Config
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

func (cc *CommonClient) Close() error {
	return cc.conn.Close()
}

func (cc *CommonClient) Conn() *grpc.ClientConn {
	return cc.conn
}
