package server

import (
	"fmt"
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
	fmt.Println("after new common ", err)
	app := &App{
		commonClient: com,
		done:         make(chan *utils.AppError, 1),
	}

	if err != nil {
		app.done <- err
	}

	app.initConfig()
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

func (a *App) initConfig() {
	config, err := a.commonClient.ConfigGet()
	if err != nil {
		a.done <- err
	}

	a.configMux.Lock()
	a.config = config
	a.configMux.Unlock()
}

func (a *App) initTrans() map[string]*commonPb.TranslationElements {
	trans, err := a.commonClient.TranslationsGet()
	if err != nil {
		a.done <- err
	}

	return trans
}
