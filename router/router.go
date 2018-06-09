/** * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *
 *
 * 路由控制器
 *
 * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zommage/livemanager/controllers/users"
)

var (
	OnePoint = "/livemanager/v1"
)

func ApiRouter(router *gin.Engine) {
	// nested group
	version1 := router.Group(OnePoint)
	{
		//
		version1.POST("/login", users.Login) // 支付回调
	}
}
