package models

type Config struct {
	Service Service `mapstructure:"service"`
}

type Service struct {
	Env                   string `mapstructure:"env"`
	GrpcPort              int    `mapstructure:"grpc_port"`
	GrpcHost              string `mapstructure:"grpc_host"`
	CommonServiceGrpcHost string `mapstructure:"common_service_grpc_host"`
	CommonServiceGrpcPort int    `mapstructure:"common_service_grpc_port"`
}
