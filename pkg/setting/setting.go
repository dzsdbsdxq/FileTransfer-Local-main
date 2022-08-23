package setting

import (
	"github.com/go-ini/ini"
	"log"
)

var (
	Cfg                 *ini.File
	UploadDir           string
	RunModel            string
	AddrIP              string
	HttpPort            int
	LimitCountPerSecond int
	JwtSecret           string
	WebName             string
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
	JwtSecret = app.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	UploadDir = app.Key("UPLOAD_DIR").MustString("./upload/")
	WebName = app.Key("WEB_NAME").MustString("点点笔记")
	DownDomain = app.Key("DOWN_DOMAIN").MustString("/d/")
}
