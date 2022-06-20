// Package im
/**
  @author:kk
  @data:2021/8/10
  @note
**/
package im

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"im_app/config"
	"im_app/internal/app/utils"
	"im_app/internal/pkg/redis"
	"im_app/pkg/response"
	log2 "im_app/pkg/zaplog"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type (
	SmApiController struct{}
	ResponseData    struct {
		Success   bool   `json:"success"`
		Code      string `json:"code"`
		Message   string `json:"message"`
		Data      Data   `json:"data"`
		RequestId string `json:"RequestId"`
	}
	Data struct {
		Token string `json:"token"`
	}
	ResponseUploadData struct {
		Success   bool        `json:"'success'"`
		Code      string      `json:"code"`
		Message   string      `json:"message"`
		Data      DataSuccess `json:"data"`
		RequestId string      `json:"RequestId"`
	}

	DataSuccess struct {
		FileId    int    `json:"file_id"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		Filename  string `json:"filename"`
		Storename string `json:"storename"`
		Size      int    `json:"size"`
		Path      string `json:"path"`
		Hash      string `json:"hash"`
		Url       string `json:"url"`
		Delete    string `json:"delete"`
		Page      string `json:"page"`
	}
)

func (*SmApiController) GetApiToken(cxt *gin.Context) {
	stringCmd := redis.DB.Get("sm_token")
	if len(stringCmd.Val()) != 0 {

		resp := new(ResponseData)
		resp.Code = "success"
		resp.Data.Token = stringCmd.Val()
		resp.Success = true

		response.SuccessResponse(resp).ToJson(cxt)
		return
	}
	data := url.Values{"username": {config.Conf.SmName}, "password": {config.Conf.SmPassword}}
	j, err := http.PostForm("https://sm.ms/api/v2/token", data)
	log2.Warning(err.Error())
	defer j.Body.Close()
	bodyC, _ := ioutil.ReadAll(j.Body)
	resp := new(ResponseData)
	json.Unmarshal(bodyC, resp)
	if resp.Success {
		response.FailResponse(http.StatusInternalServerError, resp.Message)
		return
	}
	redis.DB.Set("sm_token", resp.Data.Token, time.Hour*1)

	response.SuccessResponse(resp).ToJson(cxt)
	return
}

// @BasePath /api

// UploadImg
// @Summary 图片上传接口
// @Description 图片上传接口
// @Tags 图片上传接口
// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Param Smfile formData file true "图片上传"
// @Produce json
// @Success 200
// @Router /UploadImg [post]
func (*SmApiController) UploadImg(cxt *gin.Context) {

	file, _ := cxt.FormFile("Smfile")
	dir := utils.GetCurrentDirectory()
	path := dir + "/docs/" + file.Filename
	err := cxt.SaveUploadedFile(file, path)
	log2.LogError(err)
	header := new(utils.Header)
	header.Authorization = "Authorization"
	header.Token = config.Conf.SmToken
	resp, err := utils.PostFile(path, "https://sm.ms/api/v2/upload", header)
	log2.LogError(err)
	bodyC, _ := ioutil.ReadAll(resp.Body)
	data := new(ResponseUploadData)
	json.Unmarshal(bodyC, data)

	response.SuccessResponse(data).ToJson(cxt)
	return
}
