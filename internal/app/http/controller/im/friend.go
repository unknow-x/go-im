/**
  @author:kk
  @data:2021/9/4
  @note
**/
package im

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"im_app/internal/app/enum"
	"im_app/internal/app/http/models/friend"
	"im_app/internal/app/http/models/friend_record"
	userModel "im_app/internal/app/http/models/user"
	"im_app/internal/app/http/validates"
	"im_app/internal/app/service"
	"im_app/internal/pkg/model"
	"im_app/pkg/response"
	"im_app/pkg/zaplog"
	"net/http"
	"strconv"
	"time"
)

type FriendController struct{}

// GetList
// @Summary 获取好友列表
// @Description 获取好友列表
// @Tags 获取好友列表
// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Produce json
// @Success 200
// @Router /FriendList [get]
func (*FriendController) GetList(c *gin.Context) {
	user := userModel.AuthUser
	list, err := userModel.GetFriendListV2(user.ID)
	if err != nil {
		zaplog.Error("----获取好友列表异常", err)
		response.FailResponse(http.StatusInternalServerError, "服务器错误")
		return
	}
	response.SuccessResponse(map[string]interface{}{
		"list": list,
	}, http.StatusOK).ToJson(c)
	return
}

// GetFriendForRecord
// @Summary 获取好友申请记录
// @Description 获取好友申请记录
// @Tags 获取好友申请记录
// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Produce json
// @Success 200
// @Router /GetFriendForRecord [get]
func (*FriendController) GetFriendForRecord(c *gin.Context) {

	list, err := friend_record.GetFriendRecordList(userModel.AuthUser.ID)
	if err != nil {
		response.FailResponse(http.StatusInternalServerError, "获取好友申请记录异常").ToJson(c)
		return
	}
	response.SuccessResponse(list).ToJson(c)
	return
}

// @BasePath /api

// SendFriendRequest
// @Summary 发送好友请求
// @Description 发送好友请求接口
// @Tags 发送好友请求接口
// @Accept multipart/form-data
// @Produce json
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Param information formData string true "请求描述"
// @Param f_id formData string true "用户id"
// @Param client_type formData string false "客户端类型 0.网页端登录 1.设备端登录 2.app端"
// @Success 200
// @Router /SendFriendRequest [post]
func (*FriendController) SendFriendRequest(c *gin.Context) {
	f_id := c.PostForm("f_id")
	FId, _ := strconv.Atoi(f_id)
	if int64(FId) == userModel.AuthUser.ID {
		response.FailResponse(http.StatusUnauthorized, "请勿添加自己为好友").ToJson(c)
		return
	}
	_send := validates.SendFriendRequestFrom{
		MId: strconv.Itoa(int(userModel.AuthUser.ID)),
		FId: f_id,
	}
	errs := validates.ValidateSendFriendRequestFrom(_send)

	if len(errs) > 0 {
		response.FailResponse(http.StatusInternalServerError, "error", errs).ToJson(c)
		return
	}

	information := c.PostForm("information")

	var friend friend.ImFriends

	model.DB.Table("im_friends").
		Where("status=1 and f_id=? and m_id=?", f_id, userModel.AuthUser.ID).
		Find(&friend)

	if friend.ID == 0 {
		err := friend_record.AddRecords(userModel.AuthUser.ID, f_id, information)
		if err != nil {
			response.FailResponse(http.StatusInternalServerError, "添加失败").ToJson(c)
			return
		}

		service.SendMessage(enum.ADD_FRIEND_REQUEST, int(userModel.AuthUser.ID), FId,
			fmt.Sprintf("系统:%s请求添加你为好友", userModel.AuthUser.Name))

		response.SuccessResponse().ToJson(c)
		return
	} else {
		response.FailResponse(http.StatusUnauthorized, "已经是好友关系了，请勿重复添加")
		return
	}

}

type ImFriendRecords struct {
	ID     int64 `json:"id"`
	UserId int64 `json:"user_id"`
	FId    int64 `json:"f_id"`
	Status int   `json:"status"`
}

// @BasePath /api

// @Summary 同意好友请求接口
// @Description 同意好友请求接口
// @Tags 同意好友请求接口
// @Accept multipart/form-data
// @Produce json
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Param information formData string true "请求描述"
// @Param id formData string true "请求记录id"
// @Param status formData bool true  "1.同意 0 拒绝"
// @Success 200
// @Router /ByFriendRequest [post]
func (*FriendController) ByFriendRequest(c *gin.Context) {

	id := c.PostForm("id")
	sta := c.PostForm("status")
	status, _ := strconv.Atoi(sta)

	var friends ImFriendRecords
	err := model.DB.Where("id=?", id).
		First(&friends).Error
	if err != nil {
		response.FailResponse(http.StatusInternalServerError, "添加失败").ToJson(c)
		return
	}

	if status == 0 {
		friends.Status = 2
		model.DB.Save(&friends)
		// 投递一条消息
		service.SendMessage(enum.ADD_FRIEDN_ERROR, int(userModel.AuthUser.ID), int(friends.UserId),
			fmt.Sprintf("%s已拒绝添加你为好友", userModel.AuthUser.Name))
		response.FailResponse(http.StatusInternalServerError, "已经拒绝了~").ToJson(c)
		return
	} else {
		friend.AddFriends(friends.UserId, friends.FId)
		friend.AddFriends(friends.FId, friends.UserId)
		friends.Status = 1
		model.DB.Save(&friends)
		service.SendMessage(enum.ADD_FRIEDN_SUCCESS, int(userModel.AuthUser.ID), int(friends.UserId),
			fmt.Sprintf("%s已通过添加你为好友,一起聊天吧～", userModel.AuthUser.Name))
		// 投递一条消息
		response.SuccessResponse().ToJson(c)
		return

	}

}

// @BasePath /api

// @Summary 删除好友接口
// @Description 删除好友接口
// @Tags 删除好友接口
// @Accept multipart/form-data
// @Produce json
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Param user_id formData string true "好友id"
// @Success 200
// @Router /RemoveFriend [post]
func (*FriendController) RemoveFriend(c *gin.Context) {
	user_id := c.PostForm("user_id")

	if len(user_id) == 0 {
		response.FailResponse(http.StatusInternalServerError, "user_id不能为空~").ToJson(c)
		return
	}
	user := userModel.AuthUser

	model.DB.Where("m_id=? and f_id=?", user.ID, user_id).Delete(&friend.ImFriends{})
	model.DB.Where("m_id=? and f_id=?", user_id, user.ID).Delete(&friend.ImFriends{})

	f_id, _ := strconv.Atoi(user_id)

	service.SendMessage(enum.DELETE_FRIEND, int(user.ID), f_id, fmt.Sprintf("系统:%s已将你删除", user.Name))

	response.FailResponse(http.StatusOK, "删除成功~").ToJson(c)
	return
}

// @BasePath /api

// @Summary 好友置顶功能接口
// @Description 好友置顶功能接口
// @Tags 好友置顶功能接口
// @Accept multipart/form-data
// @Produce json
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Param user_id formData string true "好友id"
// @Param status formData int true "状态 1.置顶 0.取消置顶"
// @Success 200
// @Router /FriendPlacedTop [post]
func (*FriendController) FriendPlacedTop(c *gin.Context) {

	user_id := c.PostForm("user_id")
	if len(user_id) < 1 {
		response.ErrorResponse(http.StatusInternalServerError, "用户id不能为空").ToJson(c)
		return
	}

	status := c.DefaultPostForm("status", "0")

	_status, _ := strconv.Atoi(status)
	newBool := _status != 0

	id := userModel.AuthUser.ID

	model.DB.Model(&friend.ImFriends{}).Where("m_id=? and f_id=? ", id, user_id).
		Updates(map[string]interface{}{
			"status": newBool, "TopTime": time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"),
		})

	response.SuccessResponse().ToJson(c)

	return
}

// @BasePath /api

// @Summary 更新好友备注接口
// @Description 更新好友备注接口
// @Tags 更新好友备注接口
// @Accept multipart/form-data
// @Produce json
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Param user_id formData string true "好友id"
// @Param note formData string true "备注"
// @Success 200
// @Router /UpdateFriendNote [post]
func (*FriendController) UpdateFriendNote(c *gin.Context) {
	user_id := c.PostForm("user_id")
	note := c.PostForm("note")
	if len(user_id) < 1 || len(note) > 20 || len(note) < 2 {
		response.ErrorResponse(http.StatusInternalServerError, "参数不合格").ToJson(c)
		return
	}
	id := userModel.AuthUser.ID

	model.DB.Model(&friend.ImFriends{}).Where("m_id=? and f_id=? ", id, user_id).
		Update("note", note)

	response.SuccessResponse().ToJson(c)

	return
}
