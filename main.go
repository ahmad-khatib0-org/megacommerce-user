package main

import (
	"fmt"
	"os"

	"github.com/ahmad-khatib0-org/megacommerce-user/internal/server"
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

	err = server.RunServer(config)
	if err != nil {
		panic(err)
	}
}
