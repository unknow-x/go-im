// Package config
/**
  @author:kk
  @data:2021/6/18
  @note
**/
package config

import (
	"im_app/pkg/config"
)

type ConfST struct {
	Core
	Cache
	Database DatabaseST
	Mail
	Oauth
	Rabbitmp
}

var Conf = ConfST{}

func Initialize() {
	config.Initialize(&Conf)
}
