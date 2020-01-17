package grmath

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

func Md5(str string) string {

	h := md5.New()
	io.WriteString(h, str)

	return hex.EncodeToString(h.Sum(nil))
}
func Md5ForAes(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

func HmacSha1(key, data string) string {

	fmt.Println(key, data)
	hmac := hmac.New(sha1.New, []byte(key))
	hmac.Write([]byte(data))

	//return hex.EncodeToString(hmac.Sum(nil))

	raw := hmac.Sum(nil)

	return base64.StdEncoding.EncodeToString(raw)
}

func Sha1(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}
