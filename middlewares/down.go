package middlewares

import (
	"FileTransfer/models"
	"FileTransfer/pkg/ec"
	"FileTransfer/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

func DownCheck(c *gin.Context) {
	fileCode := com.StrTo(c.Param("fileCode")).String()
	s := models.ExistFileByPass(fileCode)
	if !s {
		utils.RespFunc(c, http.StatusForbidden, ec.ERROR_AUTH, "Wrong password")
		c.Abort()
		return
	}
	c.Next()
}
