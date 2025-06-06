package server

import (
	"sync"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/common"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

type App struct {
	commonClient *common.CommonClient
	configMux    sync.RWMutex
	config       *pb.Config
	done         chan *models.InternalError
}

func RunServer(c *models.Config) error {
	com, err := common.NewCommonClient(c)
	app := &App{
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
