package controller

import (
	"FileTransfer/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckAuth(c *gin.Context) {
	ch := utils.PostDataCheckAuth(c, []string{"fileCode"})
	fmt.Println(ch)
	utils.RespFunc(c, http.StatusOK, 200, ch)
}
