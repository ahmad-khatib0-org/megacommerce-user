package server

import (
	"fmt"

	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
)

func (a *App) initSharedConfig() {
	config, err := a.commonClient.ConfigGet()
	if err != nil {
		a.done <- err
	}

	a.configMux.Lock()
	a.config = config
	a.configMux.Unlock()
}

func LoadServiceConfig(fileName string) (*models.Config, *utils.AppError) {
	viper.AddConfigPath(".")
	viper.SetConfigFile(fileName)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, utils.NewAppError("user.load.loadConfig", "failed_to_load_config", nil, fmt.Sprintf("failed to load configurations %v", err), int(codes.Internal))
	}

	var c models.Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, utils.NewAppError("user.load.loadConfig", "failed_to_parse_config", nil, fmt.Sprintf("failed to parse configurations %v", err), int(codes.Internal))
	}

	return &c, nil
}
