package users

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zommage/livemanager/common"
	"github.com/zommage/livemanager/controllers/base"
	. "github.com/zommage/livemanager/logs"
	models "github.com/zommage/livemanager/models"
)

type LoginResp struct {
	User string `json:"user"`
	Role string `json:"role"`
}

// 用户登录
func Login(c *gin.Context) {
	fmt.Println("hello login......................")

	tokenBytes, err := c.GetRawData()
	if err != nil {
		tmpStr := fmt.Sprintf("get token err: %v", err)
		fmt.Println(tmpStr)
		Log.Info(tmpStr)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	tokenReq := &common.LoginToken{}
	err = json.Unmarshal(tokenBytes, &tokenReq)
	if err != nil {
		tmpStr := fmt.Sprintf("fmt json token err: %v", err)
		fmt.Println(tmpStr)
		Log.Info(tmpStr)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	fmt.Println("token========= ", tokenReq.Token)

	// 对 token 进行 ras s1 对齐方式 私钥解密
	loginBytes, err := common.RsaS1Decrypt(tokenReq.Token)
	if err != nil {
		tmpStr := fmt.Sprintf("descrypt s1 token err: %v", err)
		fmt.Println(tmpStr)
		Log.Info(tmpStr)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	loginReq := common.LoginReq{}
	err = json.Unmarshal(loginBytes, &loginReq)
	if err != nil {
		tmpStr := fmt.Sprintf("fmt login json err: %v", err)
		fmt.Println(tmpStr)
		Log.Info(tmpStr)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	fmt.Println("login data=========", string(loginBytes))

	err = common.NumLetterLine(loginReq.User)
	if err != nil {
		tmpStr := fmt.Sprintf("username invalid: %v", err)
		Log.Info(tmpStr)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	row, err := models.QueryUserByUsername(loginReq.User, 1)
	if err != nil {
		tmpStr := fmt.Sprintf("user not exist or user is unuse")
		Log.Infof("query by user name err: %v", err)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	// 对数据库中用户的密码进行解密
	pwdBytes, err := common.RsaS1Decrypt(row.Pwd)
	if err != nil {
		tmpStr := fmt.Sprintf("descrypt user pwd err: %v", err)
		Log.Info(tmpStr)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	if string(pwdBytes) != loginReq.Pwd {
		tmpStr := fmt.Sprintf("user pwd not match")
		Log.Infof(tmpStr)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	fmt.Println("pwd bytes: ", string(pwdBytes))

	// 更新用户信息
	nowTime := time.Now()
	row.UpdatedAt = nowTime
	row.Online = 1

	err = models.UpdateDbs(row)
	if err != nil {
		tmpStr := fmt.Sprintf("user update err: %v", err)
		Log.Infof(tmpStr)
		base.WebResp(c, 500, 500, nil, tmpStr)
		return
	}

	// 插入用户记录表
	historyRow := &models.LiveManagerUserToken{}
	historyRow.Username = loginReq.User
	historyRow.Role = row.Role
	historyRow.Token = tokenReq.Token
	historyRow.Expire = nowTime.Add(common.ExprireTime * time.Minute)
	historyRow.CreatedAt = nowTime
	err = models.InsertDbs(historyRow)
	if err != nil {
		tmpStr := fmt.Sprintf("insert user token err: %v", err)
		Log.Infof(tmpStr)
		base.WebResp(c, 500, 500, nil, tmpStr)
		return
	}

	resp := &LoginResp{
		User: row.Username,
		Role: row.Role,
	}

	base.WebResp(c, 200, 200, resp, common.Success)
	return
}
