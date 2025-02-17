/**
  @author:kk
  @data:2021/7/13
  @note
**/
package im

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"im_app/internal/app/http/models/group"
	"im_app/internal/app/http/models/group_user"
	userModel "im_app/internal/app/http/models/user"
	"im_app/internal/app/http/validates"
	"im_app/internal/pkg/model"
	"im_app/pkg/helpler"
	"im_app/pkg/response"
	"im_app/pkg/zaplog"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type (
	GroupController struct{}
	Groups          struct {
		GroupId string `json:"group_id"`
	}
)

// @BasePath /api

// List
// @Summary 获取群聊列表
// @Description 获取群聊列表
// @Tags 获取群聊列表
// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Produce json
// @Success 200
// @Router /GetGroupList [get]
func (*GroupController) List(c *gin.Context) {

	user := userModel.AuthUser
	var groupId []Groups
	err := model.DB.Table("im_group_users").
		Select("group_id").
		Where("user_id=?", user.ID).
		Group("group_id").
		Find(&groupId).Error
	if err != nil {
		fmt.Println(err)
	}
	v := reflect.ValueOf(groupId)
	groupSlice := make([]string, v.Len())
	for key, value := range groupId {
		groupSlice[key] = value.GroupId
	}
	fmt.Println(groupSlice)
	list, err := group.GetGroupUserList(groupSlice)

	if err != nil {
		zaplog.Error("----获取群聊列表异常", err)
		response.FailResponse(http.StatusInternalServerError, "服务器错误").ToJson(c)
		return
	}
	response.SuccessResponse(list).ToJson(c)
	return
}

// @BasePath /api

// Show
// @Summary 获取群聊详情
// @Description 获取群聊详情
// @Tags 获取群聊详情
// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Param group_id query string true "群聊id"
// @Produce json
// @Success 200
// @Router /GetGroupDetails [get]
func (*GroupController) Show(c *gin.Context) {
	var groups group.ImGroups
	group_id := c.Query("group_id")
	err := model.DB.Preload("Users").Where("id=?", group_id).First(&groups).Error
	if err != nil {
		zaplog.Error("----获取群聊详情异常", err)
		response.FailResponse(http.StatusInternalServerError, "服务器错误").ToJson(c)
		return
	}
	response.SuccessResponse(groups).ToJson(c)
	return
}

// @BasePath /api

// Create
// @Summary 创建群聊
// @Description 创建群聊
// @Tags 创建群聊
// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Param group_name formData string true "群聊名称"
// @Param user_id formData array true "群聊用户"
// @Produce json
// @Success 200
// @Router /CreateGroup [post]
func (*GroupController) Create(c *gin.Context) {
	user := userModel.AuthUser

	_groups := validates.CreateGroupParams{
		GroupName: c.PostForm("group_name"),
		UserId:    c.PostFormMap("user_id"),
	}
	fmt.Println(_groups)
	rules := govalidator.MapData{
		"group_name": []string{"required", "between:2,20"},
		// "user_id": []string{"required"},
	}
	opts := govalidator.Options{
		Data:          &_groups,
		Rules:         rules,
		TagIdentifier: "valid",
	}
	errs := govalidator.New(opts).ValidateStruct()

	if len(errs) > 0 {

		data, _ := json.MarshalIndent(errs, "", "  ")
		var result = helpler.JsonToMap(data)
		response.ErrorResponse(http.StatusInternalServerError, "参数不合格", result).ToJson(c)
		return
	}
	if len(_groups.UserId) > 50 {
		response.ErrorResponse(http.StatusInternalServerError, "默认只能邀请50人入群").ToJson(c)
	}

	id, err := group.Created(user.ID, _groups.GroupName)
	if err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "创建异常").ToJson(c)
		return
	}
	err = group_user.CreatedAll(_groups.UserId, id, user.ID)
	if err != nil {
		zaplog.Error("----创建群聊异常", err)
		response.ErrorResponse(http.StatusInternalServerError, "创建异常").ToJson(c)
		return
	}
	response.SuccessResponse(map[string]interface{}{
		"group_id": id,
	}).ToJson(c)
	return
}

// @BasePath /api

// RemoveGroup
// @Summary 删除群聊
// @Description 删除群聊
// @Tags 删除群聊
// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Param group_id formData string true "群聊id"
// @Produce json
// @Success 200
// @Router /RemoveGroup [post]
func (*GroupController) RemoveGroup(cxt *gin.Context) {
	group_id := cxt.PostForm("group_id")
	if len(group_id) == 0 {
		response.ErrorResponse(http.StatusInternalServerError, "参数不合格").ToJson(cxt)
		return
	}
	model.DB.Where("id=?", group_id).Delete(&group.ImGroups{})
	model.DB.Where("group_id=?", group_id).Delete(&group_user.ImGroupUsers{})

	response.SuccessResponse().ToJson(cxt)
	return
}

// @BasePath /api

// RemovedUserFromGroup
// @Summary 移除群聊用户
// @Description 移除群聊用户
// @Tags 移除群聊用户
// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Param group_id formData string true "群聊id"
// @Param user_id formData string true "用户id"
// @Produce json
// @Success 200
// @Router /RemovedUserFromGroup [post]
func (*GroupController) RemovedUserFromGroup(c *gin.Context) {

	_group := validates.GroupFrom{
		GroupId: c.PostForm("group_id"),
		UserId:  c.PostForm("user_id"),
	}
	errs := validates.ValidateGroupForm(_group)

	if len(errs) > 0 {
		response.FailResponse(http.StatusUnauthorized, "error", errs)
	}
	g_id, _ := group.GetGroupUserId(_group.GroupId)

	if userModel.AuthUser.ID != g_id {
		response.FailResponse(http.StatusUnauthorized, "没有权限删除群成员！").ToJson(c)
		return
	}

	model.DB.Table("im_group_users").
		Where("user_id and group_id=?", _group.UserId, _group.GroupId).
		Delete(&group_user.ImGroupUsers{})

	response.SuccessResponse().ToJson(c)
	return
}

// @BasePath /api

// JoinGroup
// @Summary 添加用户到指定群聊
// @Description 添加用户到指定群聊
// @Tags 添加用户到指定群聊
// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Param group_id formData string true "群聊id"
// @Param user_id formData string true "用户id"
// @Produce json
// @Success 200
// @Router /JoinGroup [post]
func (*GroupController) JoinGroup(c *gin.Context) {

	_group := validates.GroupFrom{
		GroupId: c.PostForm("group_id"),
		UserId:  c.PostForm("user_id"),
	}
	errs := validates.ValidateGroupForm(_group)

	if len(errs) > 0 {
		response.FailResponse(http.StatusUnauthorized, "error", errs).ToJson(c)
		return
	}

	u_id, _ := group.GetGroupUserId(_group.GroupId)

	if userModel.AuthUser.ID != u_id {
		response.FailResponse(http.StatusUnauthorized, "没有资格邀请群成员！").ToJson(c)
		return
	}

	isExist := group_user.GetGroupUser(_group.GroupId, _group.UserId)

	if isExist == true {
		response.FailResponse(http.StatusUnauthorized, "已经在群聊里面了！").ToJson(c)
		return
	}

	var user userModel.Users
	err := model.DB.First(&user, _group.UserId).Error
	if err != nil {
		fmt.Println(err)
	}
	userID, _ := strconv.Atoi(_group.UserId)
	groupID, _ := strconv.Atoi(_group.GroupId)

	model.DB.Table("im_group_users").Create(&group_user.ImGroupUsers{
		UserId:    int64(userID),
		GroupId:   int64(groupID),
		Remark:    c.PostForm("remark"),
		CreatedAt: time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"),
		Avatar:    user.Avatar,
		Name:      user.Name,
	})

	response.SuccessResponse().ToJson(c)
	return

}
