// Package config
/**
  @author:kk
  @data:2021/6/18
  @note
**/
package config

type Cache struct {
	Redis
}

type Redis struct {
	Addr     string
	Port     int
	Password string
	DB       int `mapstructure:"db"`
}
