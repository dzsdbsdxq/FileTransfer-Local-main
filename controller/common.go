package controller

import (
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
		"web_site": "点点笔记",
		"web_url":  "https://mapi.net.cn/",
		"driver":   "阿里云盘",
		"miit":     "粤ICP备2020114467号",
	}
	utils.RespFunc(c, http.StatusOK, 200, data)

}
