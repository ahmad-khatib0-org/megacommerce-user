package server

import (
	"sync"

	commonPb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/common"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
)

type App struct {
	commonClient *common.CommonClient
	configMux    sync.RWMutex
	config       *commonPb.Config
	done         chan *utils.AppError
	utils        *utils.Utils
}

func RunServer(c *models.Config) *utils.AppError {
	com, err := common.NewCommonClient(c)
	app := &App{
		commonClient: com,
		done:         make(chan *utils.AppError, 1),
	}

	if err != nil {
		app.done <- err
	}

	app.initSharedConfig()
	trans := app.initTrans()

	utils, err := utils.NewUtils(&utils.UtilsArgs{AllTrans: trans})
	if err != nil {
		app.done <- err
	}
	app.utils = utils

	err = <-app.done
	if err != nil {
		// TODO: cleanup things
	}

	return err
}
