package server

import (
	com "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/spf13/viper"
)

func (a *Server) initSharedConfig() {
	config, err := a.commonClient.ConfigGet()
	if err != nil {
		a.errors <- err
	}

	a.configFn = func() *com.Config { return config }
	a.configMux.Lock()
	a.config = config
	a.configMux.Unlock()
}

func LoadServiceConfig(fileName string) (*models.Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigFile(fileName)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var c models.Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
