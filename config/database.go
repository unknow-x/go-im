// Package config
/**
  @author:kk
  @data:2021/6/18
  @note
**/
package config

type DatabaseST struct {
	Mysql
}

type Mysql struct {
	Host               string
	Port               int
	Database           string
	Username           string
	Password           string
	Charset            string
	Loc                string
	MaxIdleConnections int `mapstructure:"max_idle_connections"`
	MaxOpenConnections int `mapstructure:"max_open_connections"`
	MaxLifeSeconds     int `mapstructure:"max_life_seconds"`
}
