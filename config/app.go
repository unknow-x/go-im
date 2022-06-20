// Package config
/**
  @author:kk
  @data:2021/6/18
  @note
**/
package config

type Core struct {
	Name         string
	Ym           string
	GoCoroutines int `mapstructure:"go_coroutines"`
	Node         string
	Env          string
	FileDisk     string `mapstructure:"file_disk"`
	Port         string
	GrpcPort     int    `mapstructure:"grpc_port"`
	TcpPort      int    `mapstructure:"tcp_port"`
	SwaggerPort  int    `mapstructure:"swagger_port"`
	LogAddress   string `mapstructure:"log_address"`
	GiteeApiKey  string `mapstructure:"gitee_api_key"`
	GaodeKey     string `mapstructure:"gaode_key"`
	Jwt
	Base64          string
	SmName          string `mapstructure:"sm_name"`
	SmPassword      string `mapstructure:"sm_password"`
	SmToken         string `mapstructure:"sm_token"`
	AppYpId         string `mapstructure:"app_yp_id"`
	AppYpKey        string `mapstructure:"app_yp_key"`
	AppYpSecretKey  string `mapstructure:"app_yp_secret_key"`
	AppYpSignKey    string `mapstructure:"app_yp_sign_key"`
	AppClusterModel bool   `mapstructure:"app_cluster_mode"`
}

type Jwt struct {
	SignKey        string `mapstructure:"sign_key"`
	ExpirationTime int64  `mapstructure:"expiration_time"`
}
