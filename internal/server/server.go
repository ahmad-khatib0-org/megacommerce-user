package server

import (
	"sync"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/common"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/logger"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

type Server struct {
	commonClient *common.CommonClient
	configMux    sync.RWMutex
	config       *pb.Config
	done         chan *models.InternalError
}

type ServerArgs struct {
	Log *logger.Logger
	Cfg *models.Config
}

func RunServer(s *ServerArgs) error {
	com, err := common.NewCommonClient(s.Cfg)
	app := &Server{
		commonClient: com,
		done:         make(chan *models.InternalError, 1),
	}

	if err != nil {
		app.done <- err
	}

	app.initSharedConfig()
	app.initTrans()

	err = <-app.done
	if err != nil {
		// TODO: cleanup things
	}

	return err
}
