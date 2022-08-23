package operate

import (
	"FileTransfer/pkg/drivers"
	"FileTransfer/pkg/setting"
	"net/http"
)

func Upload(r *http.Request) (string, error) {
	u, err := drivers.New(setting.UploadDir, "2006/01/02/", 1024*1024*1024)
	if err != nil {
		return "", err
	}
	files, err := u.Do("files", r)
	if err != nil {
		return "", err
	}
	return files[0], err
}

func Delete(path string) bool {
	return drivers.DeleteFile(setting.UploadDir + path)
}

func GetFileMd5(path string) (string, error) {
	return drivers.GetFileMd5(setting.UploadDir + path)
}

func PersonalInfo() (map[string]interface{}, error) {
	return nil, nil
	//driver, ok := base.GetDriver(setting.Drivers)
	//if !ok {
	//	return base.Json{}, fmt.Errorf("no [%s] driver", setting.Drivers)
	//}
	//
	//req, err := driver.PersonalInfo()
	//if err != nil {
	//	logging.Logger.Errorf("upload error: %s", err.Error())
	//}
	//debug.FreeOSMemory()
	//return map[string]interface{}{
	//	"total_size": req.TotalSize,
	//	"used_size":  req.UsedSize,
	//}, err
}
