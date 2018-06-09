package base

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zommage/livemanager/common"
	. "github.com/zommage/livemanager/logs"
	models "github.com/zommage/livemanager/models"
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

// 针对请求进行鉴权与签名校验
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if common.AuthSwitch != false {
			c.Next()
			return
		}

		var err error
		urlArr := strings.Split(c.Request.URL.Path, "?")
		preUrl := urlArr[0]
		fmt.Println("pre url=============", preUrl)

		// 如果为登录接口, 直接跳到 controllers
		if preUrl == "/login" {
			c.Next()
			return
		}

		// 鉴权
		err = AuthFunc(c)
		if err != nil {
			if err == common.TokenExprire {
				fmt.Println(common.TokenExprire)
				Log.Error(common.TokenExprire)
				WebResp(c, 403, 777, nil, fmt.Sprintf("%v", common.TokenExprire))
				c.Abort()
				return
			}
			tmpStr := fmt.Sprintf("auth fail: %v", err)
			fmt.Println(tmpStr)
			Log.Errorf(tmpStr)
			WebResp(c, 401, 401, nil, tmpStr)
			c.Abort()
			return
		}

		c.Next()
		return
	}
}

// 鉴权函数
func AuthFunc(c *gin.Context) error {
	token := c.GetHeader("token")

	if token == "" {
		return fmt.Errorf("token is nil")
	}

	// 对 token 进行解析
	err := ParseToken(token)
	if err != nil {
		return err
	}

	return nil
}

// 对 token 进行校验
func ParseToken(token string) error {
	// 对数据库中用户的密码进行解密
	tokenRow, err := models.QueryByToken(token)
	if err != nil {
		return fmt.Errorf("token is not exist")
	}

	// 判断 token 是否已经过期
	if tokenRow.Expire.Before(time.Now()) != false {
		return common.TokenExprire
	}

	return nil
}

// 获取签名的公共参数
func ComSigParam(c *gin.Context) (map[string]interface{}, string, error) {
	params := make(map[string]interface{})

	TimeStamp := c.Query("TimeStamp")
	if TimeStamp == "" {
		return nil, "", fmt.Errorf("signature timestamp is nil")
	}

	SignatureNonce := c.Query("SignatureNonce")
	if SignatureNonce == "" {
		return nil, "", fmt.Errorf("signature noce is nil")
	}

	// 签名
	signature := c.Query("Signature")
	if signature == "" {
		return nil, "", fmt.Errorf("signature is nil")
	}

	// 如果签名中带有 + 号, url会将 + 解析成空格
	signature = strings.Replace(signature, " ", "+", -1)

	// 增加公共参数, 时间格式为YYYY-MM-DDThh:mm:ssZ,例如，2014-11-11T12:00:00Z
	params["TimeStamp"] = TimeStamp

	// 随机字符串
	params["SignatureNonce"] = SignatureNonce

	return params, signature, nil
}
