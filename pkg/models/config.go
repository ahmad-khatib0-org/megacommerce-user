package models

type Config struct {
	Service Service `mapstructure:"service"`
}

type Service struct {
	Env      string `mapstructure:"env"`
	GrpcPort int    `mapstructure:"grpc_port"`
	GrpcHost string `mapstructure:"grpc_host"`
}
