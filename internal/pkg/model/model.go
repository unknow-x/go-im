// Package model
/**
  @author:kk
  @data:2021/6/18
  @note
**/
package model

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"im_app/config"
	"runtime"
)

var DB *gorm.DB

type BaseModel struct {
	ID int64
}

// ConnectDB 初始化 grom
func ConnectDB() *gorm.DB {
	fmt.Println(runtime.Caller(2))
	var (
		host     = config.Conf.Database.Host
		port     = config.Conf.Database.Port
		database = config.Conf.Database.Database
		username = config.Conf.Database.Username
		password = config.Conf.Database.Password
		charset  = config.Conf.Database.Charset
		//loc  = config.Conf.Database.Loc
		err error
	)
	//&parseTime=True
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		username, password, host, port, database, charset)

	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
	//	username, password, host, port, database, charset, true, url.QueryEscape(loc))
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		zap.L().Fatal("db init error", zap.Error(err))
	}

	return DB
}
