package common

import (
	"errors"
)

type LoginReq struct {
	User string `json:"user"`
	Pwd  string `json:"pwd"`
}

type LoginToken struct {
	Token string `json:"token"`
}

const (
	Success     = "success"
	ExprireTime = 60
)

var (
	TokenExprire = errors.New("token exprite")

	// 鉴权开关默认打开
	AuthSwitch = false
)
