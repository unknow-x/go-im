// Package im
/**
  @author:kk
  @data:2021/8/12
  @note
**/
package im

import (
	"github.com/gin-gonic/gin"
	"im_app/config"
	"im_app/internal/app/utils"
	"im_app/pkg/response"
)

type UploadController struct{}

// @BasePath /api

// UploadVoiceFile
// @Summary 音频文件上传接口
// @Description 音频文件上传接口
// @Tags 音频文件上传接口
// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name Authorization
// @Param Authorization	header string true "Bearer 31a165baebe6dec616b1f8f3207b4273"
// @Param voice formData file true "图片上传"
// @Produce json
// @Success 200
// @Router /UploadVoiceFile [post]
func (*UploadController) UploadVoiceFile(cxt *gin.Context) {
	voice, _ := cxt.FormFile("voice")
	dir := utils.GetCurrentDirectory()
	// 上传文件至指定目录 没找到第三方免费的第三方存储 先用自己的吧
	path := dir + "/voice/" + voice.Filename
	cxt.SaveUploadedFile(voice, path)
	response.SuccessResponse(map[string]interface{}{
		"url": config.Conf.Ym + "voice/" + voice.Filename,
	}).ToJson(cxt)
	return
}
