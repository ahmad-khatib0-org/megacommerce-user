package main

import (
	"fmt"
	"os"

	"github.com/ahmad-khatib0-org/megacommerce-user/internal/server"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/logger"
)

func main() {
	env := os.Getenv("ENV")
	if env != "dev" && env != "test" && env != "prod" {
		env = "dev"
	}
	config, err := server.LoadServiceConfig(fmt.Sprintf("config.%s.yaml", env))
	if err != nil {
		panic(err)
	}

	logger, err := logger.InitLogger(config.Service.Env)
	if err != nil {
		panic(err)
	}

	srv := &server.ServerArgs{Log: logger, Cfg: config}
	server.RunServer(srv)
}
