// Package app
/**
  @author:kk
  @data:2021/9/8
  @note
**/
package app

import (
	"github.com/gin-gonic/gin"
	conf "im_app/config"
	"im_app/internal/app/router"
	"im_app/internal/app/ws"
	"im_app/pkg/zaplog"
)

func StartHttp() {

	app := gin.Default()

	SetupPool()

	go ws.ImManager.Start()

	if conf.Conf.AppClusterModel == true {
		go ws.StartRpc()
	}

	router.RegisterApiRoutes(app)
	router.RegisterIMRouters(app)

	app.Use(zaplog.Recover)

	_ = app.Run(":" + string(conf.Conf.Core.Port))

}
