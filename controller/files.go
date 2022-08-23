package controller

import (
	"FileTransfer/models"
	"FileTransfer/pkg/drivers/operate"
	"FileTransfer/pkg/ec"
	"FileTransfer/pkg/setting"
	"FileTransfer/pkg/utils"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

type RetData struct {
	Path string `json:"path"`
	St   string `json:"st"`
}

func UploadFiles(c *gin.Context) {

	form, err := c.MultipartForm()
	if err != nil {
		utils.RespFunc(c, http.StatusBadRequest, ec.ERROR, err)
		return
	}
	filePath, err := operate.Upload(c.Request)
	if err != nil {
		utils.RespFunc(c, http.StatusBadRequest, ec.ERROR, err)
		return
	}
	file := form.File["files"][0]
	password := utils.RandPass(6)
	fileMd5, _ := operate.GetFileMd5(filePath)

	data := make(map[string]interface{})
	data["file_name"] = file.Filename
	data["file_id"] = fileMd5
	data["upload_id"] = c.Request.Header.Get("X-Request-Id")
	data["password"] = password
	data["uid"] = 0
	data["views"] = 0
	data["downloads"] = 0
	data["file_size"] = com.HumaneFileSize(uint64(file.Size))
	data["file_path"] = filePath
	autoId, expire := models.AddFiles(data)

	if autoId > 0 {
		//组装返回数据Data
		var retData = make(map[string]interface{}, 0)
		retData["pwd"] = password
		retData["expire"] = expire
		retData["file_size"] = data["file_size"]
		retData["file_name"] = file.Filename
		retData["file_id"] = data["file_id"]
		retData["url"] = password
		retData["web_name"] = setting.WebName
		utils.RespFunc(c, http.StatusOK, ec.SUCCESS, retData)
		return
	}
	//如果入数据库失败，删除文件
	operate.Delete(filePath)
	utils.RespFunc(c, http.StatusBadRequest, ec.ERROR, "")
}

func ShowDetailFiles(c *gin.Context) {
	fileCode := com.StrTo(c.PostForm("fileCode")).String()
	valid := validation.Validation{}

	valid.Required(fileCode, "fileCode").Message("取件码不能为空")

	if valid.HasErrors() {
		utils.RespFunc(c, http.StatusBadRequest, 400, valid.Errors)
		return
	}
	isAuth := utils.PostDataCheckAuth(c, []string{"fileCode"})
	if !isAuth {
		utils.RespFunc(c, http.StatusBadRequest, 400, "参数校验失败")
		return
	}

	data := make(map[string]interface{}, 0)
	if models.ExistFileByPass(fileCode) {
		fileDetail := models.GetFile(fileCode)
		models.UpdateViews(fileDetail.ID)
		data["id"] = fileDetail.ID
		data["file_name"] = fileDetail.FileName
		data["password"] = fileDetail.Password
		data["ctime"] = fileDetail.Ctime
		data["views"] = fileDetail.Views
		data["downloads"] = fileDetail.Downloads
		data["expire"] = fileDetail.Expire
		data["file_id"] = fileDetail.FileId
		data["file_size"] = fileDetail.FileSize
		data["url"] = setting.DownDomain + fileDetail.Password

		utils.RespFunc(c, http.StatusOK, ec.SUCCESS, data)
		return
	} else {
		utils.RespFunc(c, http.StatusBadRequest, ec.ERROR_NOT_EXIST_ARTICLE, "")
		return
	}

}
