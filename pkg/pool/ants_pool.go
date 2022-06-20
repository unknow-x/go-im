// Package pool
/**
  @author:kk
  @data:2021/8/18
  @note
**/
package pool

import (
	"github.com/panjf2000/ants/v2"
	"im_app/config"
)

var AntsPool *ants.Pool

func ConnectPool() *ants.Pool {
	//设置数量
	AntsPool, _ = ants.NewPool(config.Conf.GoCoroutines)
	return AntsPool
}
