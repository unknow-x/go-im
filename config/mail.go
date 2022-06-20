// Package config
/**
  @author:kk
  @data:2021/12/5
  @note
**/
package config

type Mail struct {
	Driver     string
	Host       string
	Port       int
	Name       string
	Password   string
	Encryption string
	Address    string
	FromName   string `mapstructure:"from_name"`
}
