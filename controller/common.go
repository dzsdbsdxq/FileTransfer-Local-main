package controller

import (
	"FileTransfer/pkg/setting"
	"FileTransfer/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "",
	})
}

func CopyRight(c *gin.Context) {

	data := map[string]interface{}{
		"web_site":      setting.WebName,
		"web_url":       setting.WebUrl,
		"web_sub_title": setting.WebSubTitle,
		"web_notice":    setting.WebNotice,
		"upload_title":  setting.UploadTitle,
		"upload_desc":   setting.UploadDesc,
		"mit":           setting.WebMit,
	}
	utils.RespFunc(c, http.StatusOK, 200, data)

}
