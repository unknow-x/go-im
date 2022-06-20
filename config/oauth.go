// Package config
/**
  @author:kk
  @data:2021/6/18
  @note
**/
package config

type Oauth struct {
	WbClientId     string `mapstructure:"wb_client_id"`
	WbClientSecret string `mapstructure:"wb_client_secret"`
	WbRedirectUri  string `mapstructure:"wb_redirect_uri"`
	GeClientId     string `mapstructure:"ge_client_id"`
	GeClientSecret string `mapstructure:"ge_client_secret"`
	GeRedirectUri  string `mapstructure:"ge_redirect_uri"`
}
