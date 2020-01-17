package grmath

import (
	"testing"
)

func Test_Sign(t *testing.T) {
	// go test -v -test.run Test_Sign
	m := make(map[string]interface{})
	appid := "xx"
	openid := "xxx"

	m["wx_appid"] = appid
	m["wx_openid"] = openid
	sign := SortOrderUpperMd5Sign(m, "xxx")

	t.Log(sign)

}
