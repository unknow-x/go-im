/**
  @author:kk
  @data:2021/9/16
  @note
**/
package cache

import (
	"encoding/json"
	"im_app/internal/pkg/redis"
	"im_app/pkg/zaplog"
	"strconv"
)

func getGroupIdsStr(groupId int) string {
	return "core:group:" + strconv.Itoa(groupId)
}

// todo
func getGroup(gId int) map[int]int {
	groupId := make(map[int]int)
	str := getGroupIdsStr(gId)
	data := redis.DB.Get(str)
	if len(data.Val()) > 0 {
		byData, err := data.Bytes()
		if err != nil {
			zaplog.Error("----获取群组用户id失败", err)
		}
		json.Unmarshal(byData, &groupId)
	} else {

	}

	return groupId
}
