package users

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zommage/livemanager/controllers/base"
)

// 用户登录
func Login(c *gin.Context) {
	fmt.Println("hello login......................")

	base.WebResp(c, 200, 200, nil, "success")
	return
}
