package jtexpress

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
)

// Signer 签名接口
type Signer interface {
	Sign(content string) (string, error)
	SignStruct(data interface{}) (string, error)
}

// MD5Signer MD5签名实现
type MD5Signer struct {
	key string
}

// NewMD5Signer 创建MD5签名器
func NewMD5Signer(key string) *MD5Signer {
	return &MD5Signer{key: key}
}

// Sign 签名字符串
func (s *MD5Signer) Sign(content string) (string, error) {
	str := content + s.key
	md5Sum := md5.Sum([]byte(str))
	return base64.StdEncoding.EncodeToString(md5Sum[:]), nil
}

// SignStruct 签名结构体
func (s *MD5Signer) SignStruct(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return s.Sign(string(jsonBytes))
}
