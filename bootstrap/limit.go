package bootstrap

import (
	"FileTransfer/middlewares"
	"FileTransfer/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//限流器
func LimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//得到ip地址
		ipAddr := utils.GetRealIp(c)

		//ipAddr:="127.0.0.1"
		limiter := middlewares.RateLimiter.GetLimiter(ipAddr)
		if !limiter.Allow() {
			utils.RespFunc(c, http.StatusBadRequest, 2004, "访问超出限制")
			c.Abort()
			return
		} else {
			fmt.Println("allow,next")
			c.Next()
		}
	}
}
