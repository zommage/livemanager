package base

import (
	"github.com/gin-gonic/gin"
)

/* 发往前台的公共接口
*  statusCode : http的状态码
*  errCode: 错误码
*  Msg: 信息
 */
func WebResp(c *gin.Context, statusCode, errCode int, data interface{}, Msg string) {
	respMap := map[string]interface{}{"code": errCode, "msg": Msg, "data": data}
	c.JSON(statusCode, respMap)
	return
}
