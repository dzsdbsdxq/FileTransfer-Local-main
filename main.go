package main

import (
	_ "FileTransfer/middlewares"
	"FileTransfer/pkg/logging"
	"FileTransfer/pkg/setting"
	"FileTransfer/routers"
	"fmt"
)

func init() {
	//err := middlewares.SetupIPRateLimiter()
	//if err != nil {
	//	logging.Logger.Fatalf("init.SetupIPRateLimiter err: %v", err)
	//}
}
func main() {
	r := routers.InitRouter()
	base := fmt.Sprintf("%s:%d", setting.AddrIP, setting.HttpPort)
	logging.Logger.Infof("start server @ %s", base)
	err := r.Run(base)
	if err != nil {
		logging.Logger.Errorf("failed to start: %s", err.Error())
	}
}
