package models

import (
	"FileTransfer/pkg/logging"
	"github.com/jinzhu/gorm"
	"time"
)

type Files struct {
	Model
	Uid       int    `json:"uid" gorm:"index"`
	FileName  string `json:"file_name"`
	FileId    string `json:"file_id"`
	UploadId  string `json:"upload_id"`
	Password  string `json:"password"`
	Views     int    `json:"views"`
	Downloads int    `json:"downloads"`
	Expire    int    `json:"expire"`
	FileSize  string `json:"file_size"`
	FilePath  string `json:"file_path"`
}

func (files *Files) BeforeCreate(scope *gorm.Scope) error {
	nowTime := time.Now().Unix()
	scope.SetColumn("Ctime", nowTime)
	//设置文件一天的过期时间86400
	scope.SetColumn("Expire", nowTime+86400)
	return nil
}

func GetFilesTotal(maps interface{}) (count int) {
	db.Model(&Files{}).Where(maps).Count(&count)
	return
}

func GetFiles(pageNum int, pageSize int, maps interface{}) (file []Files) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&file)
	return
}

func GetFile(password string) (file Files) {
	db.Where("password = ?", password).First(&file)
	return
}

func AddFiles(data map[string]interface{}) (int, int) {
	files := &Files{
		Uid:       data["uid"].(int),
		FileName:  data["file_name"].(string),
		FileId:    data["file_id"].(string),
		UploadId:  data["upload_id"].(string),
		Password:  data["password"].(string),
		Views:     data["views"].(int),
		Downloads: data["downloads"].(int),
		FileSize:  data["file_size"].(string),
		FilePath:  data["file_path"].(string),
	}
	db.Create(files)
	return files.ID, files.Expire
}

func DeleteFiles(id int) bool {
	logging.Logger.Infof("删除文件记录：%d", id)
	db.Where("id = ?", id).Delete(Files{})
	return true
}
func DeleteFilesByFileId(fileId string) bool {
	db.Where("file_id = ?", fileId).Delete(Files{})
	return true
}

func ExistFileByPass(password string) bool {
	var files Files
	db.Select("id").Where("password = ?", password).First(&files)
	if files.ID > 0 {
		return true
	}
	return false
}

func UpdateViews(id int) {
	////// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = '2';
	db.Model(&Files{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + ?", 1))
}
func UpdateDownload(id int) {
	db.Model(&Files{}).Where("id = ?", id).UpdateColumn("downloads", gorm.Expr("downloads + ?", 1))
}
