package setting

import (
	"github.com/go-ini/ini"
	"log"
)

var (
	Cfg                 *ini.File
	UploadTitle         string
	UploadDir           string
	UploadDesc          string
	RunModel            string
	AddrIP              string
	HttpPort            int
	LimitCountPerSecond int
	WebName             string
	WebUrl              string
	WebSubTitle         string
	WebMit              string
	WebNotice           string
	DownDomain          string
)

func init() {
	var err error
	Cfg, err = ini.Load("config/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini':%v", err)
	}
	loadConfig()
}

func loadConfig() {
	//加载项目运行模式
	RunModel = Cfg.Section("").Key("RUN_MODE").MustString("debug")
	//加载项目运行HTTP配置
	server, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server':%v", err)
	}
	AddrIP = server.Key("ADDR_IP").MustString("0.0.0.0")
	HttpPort = server.Key("HTTP_PORT").MustInt(8999)
	//加载项目APP配置
	app, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app':%v", err)
	}
	LimitCountPerSecond = app.Key("LIMIT_COUNT_PER_SECOND").MustInt(7)
	UploadDir = app.Key("UPLOAD_DIR").MustString("./runtime/upload/")
	UploadTitle = app.Key("UPLOAD_TITLE").MustString("")
	UploadDesc = app.Key("UPLOAD_DESC").MustString("")
	WebName = app.Key("WEB_NAME").MustString("点点笔记")
	WebUrl = app.Key("WEB_URL").MustString("https://mapi.net.cn")
	WebMit = app.Key("WEB_MIT").MustString("粤ICP备2020114467号")
	WebSubTitle = app.Key("WEB_SUB_TITLE").MustString("")
	WebNotice = app.Key("WEB_NOTICE").MustString("欢迎使用文件快传")
	DownDomain = app.Key("DOWN_DOMAIN").MustString("/api/d/")
}
