package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

// ----------------------------------------------------------

type Mac struct {
	AccessKey string
	SecretKey []byte
}

func NewMac(accessKey, secretKey string) (mac *Mac) {
	return &Mac{accessKey, []byte(secretKey)}
}

func (mac *Mac) Sign(data []byte) (token string) {

	h := hmac.New(sha1.New, mac.SecretKey)
	h.Write(data)

	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return mac.AccessKey + ":" + sign[:27]
}
