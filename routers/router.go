package routers

import (
	"FileTransfer/controller"
	"FileTransfer/middlewares"
	"FileTransfer/pkg/setting"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	//r.Use(bootstrap.LimitMiddleware())

	gin.SetMode(setting.RunModel)

	Cors(r)
	r.StaticFS("/web", http.Dir("./web"))
	r.LoadHTMLGlob("./web/*.html")
	r.GET("/", controller.IndexPage)
	r.POST("/api/check", controller.CheckAuth)
	r.GET("/api/d/:fileCode", middlewares.DownCheck, controller.Down)
	r.POST("/api/detail", controller.ShowDetailFiles)
	//上传
	r.POST("/api/upload", controller.UploadFiles)
	r.POST("/api/copyright", controller.CopyRight)

	return r
}

func Cors(r *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization", "range")
	r.Use(cors.New(config))
}
