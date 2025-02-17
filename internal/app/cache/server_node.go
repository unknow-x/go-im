// Package cache
/**
  @author:kk
  @data:2021/11/13
  @note
**/
package cache

import (
	"im_app/config"
	"im_app/internal/pkg/redis"
	"strconv"
	"sync"
)

// 必须有一个结构体去实现该接口的方法

type ServiceNodeHandler interface {
	SetUserServiceNode(ID int64)        // 设置用户节点
	GetUserServiceNode(ID int64) string // 获取用户节点
	DelUserServiceNode(ID int64)        // 删除用户节点
}

type ServiceNode struct {
	mu sync.Mutex
}

func getUserIdStr(ID int64) string {
	// 注意 ：写法 在可视化工具里面可以更好的看到缓存的结构体
	return "core:node:user:" + strconv.Itoa(int(ID))
}

func (node *ServiceNode) GetUserServiceNode(ID int64) string {
	var key = getUserIdStr(ID)
	StringCmd := redis.DB.Get(key)
	return StringCmd.Val()
}

func (node *ServiceNode) SetUserServiceNode(ID int64) {
	var key = getUserIdStr(ID)
	cacheNode := config.Conf.Node + ":" + strconv.Itoa(config.Conf.GrpcPort)
	var value = cacheNode
	node.mu.Lock()
	redis.DB.Set(key, value, 0)
	node.mu.Unlock()
}

func (node *ServiceNode) DelUserServiceNode(ID int64) {
	var key = getUserIdStr(ID)
	redis.DB.Del(key)
}
