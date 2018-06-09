package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/zommage/livemanager/conf"
)

var (
	rsaSertKey []byte // rsa 私钥文件
)

func InitRsaKey() error {
	var err error

	rsaSertKey, err = ioutil.ReadFile(conf.Conf.BaseConf.RsaSertKey)
	if err != nil {
		return err
	}

	return nil
}

// 私钥解密的时候，有可能是 s1 和 s8 两种对齐的, 需要轮流试
func RsaS8Decrypt(baseEncryMsg string) ([]byte, error) {
	block, _ := pem.Decode(rsaSertKey)
	if block == nil {
		return nil, fmt.Errorf("private key error!")
	}

	/// s8 对齐方式的解密
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		tmpStr := fmt.Sprintf("private s8 key err: %v", err)
		return nil, fmt.Errorf(tmpStr)
	}

	// s8 对齐的
	privKey := priv.(*rsa.PrivateKey)

	//base64 解码
	decryptMsg, err := base64.StdEncoding.DecodeString(baseEncryMsg)
	if err != nil {
		tmpStr := fmt.Sprintf("base64 decode err: %v", err)
		return nil, fmt.Errorf(tmpStr)
	}

	return rsa.DecryptPKCS1v15(rand.Reader, privKey, decryptMsg)
}

// 私钥解密 s1对齐方式解密
func RsaS1Decrypt(baseEncryMsg string) ([]byte, error) {
	block, _ := pem.Decode(rsaSertKey)
	if block == nil {
		return nil, fmt.Errorf("private key error!")
	}

	// s1对齐的方式解密
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		tmpStr := fmt.Sprintf("s1 private key err: %v", err)
		return nil, fmt.Errorf(tmpStr)
	}

	// base64 解码
	decryptMsg, err := base64.StdEncoding.DecodeString(baseEncryMsg)
	if err != nil {
		tmpStr := fmt.Sprintf("base64 decode err: %v", err)
		fmt.Println(tmpStr)
		return nil, fmt.Errorf(tmpStr)
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, decryptMsg)
}
