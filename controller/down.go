package controller

import (
	"FileTransfer/models"
	"FileTransfer/pkg/ec"
	"FileTransfer/pkg/setting"
	"FileTransfer/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

func Down(c *gin.Context) {
	fileCode := com.StrTo(c.Param("fileCode")).String()

	if models.ExistFileByPass(fileCode) {
		fileDetail := models.GetFile(fileCode)
		models.UpdateDownload(fileDetail.ID)

		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+fileDetail.FileName) // 用来指定下载下来的文件名
		c.Header("Content-Transfer-Encoding", "binary")
		c.File(setting.UploadDir + fileDetail.FilePath)

		return
	} else {
		utils.RespFunc(c, http.StatusBadRequest, ec.ERROR_NOT_EXIST_ARTICLE, "")
		return
	}

}
