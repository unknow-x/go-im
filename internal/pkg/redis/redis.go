// Package redis
/**
  @author:kk
  @data:2021/6/18
  @note
**/
package redis

import (
	"github.com/go-redis/redis"
	"im_app/config"
	"strconv"
	"time"
)

var DB *redis.Client

// InitClient redis 连接
func InitClient() (err error) {

	DB = redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         config.Conf.Cache.Redis.Addr + ":" + strconv.Itoa(config.Conf.Cache.Redis.Port),
		Password:     config.Conf.Cache.Redis.Password,
		DB:           config.Conf.Cache.Redis.DB,
		PoolSize:     15, //连接池 默认为4倍cpu数
		MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolTimeout:  5 * time.Second,
	})
	_, err = DB.Ping().Result()

	return err
}
