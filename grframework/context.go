package grframework

import (
	"errors"
	"fmt"

	"github.com/valyala/fasthttp"
)

/*
type RequestContext struct {
	FasthttpCtx *fasthttp.RequestCtx
}

type ResponseContext struct {
	Headers map[string]string
	Errcode  int
	Msg     string
}*/

type Context struct {
	FasthttpCtx *fasthttp.RequestCtx
	Headers     map[string]string
}

/*
func MakeErrcodeWithMsg(result int, msg string) (rc *ResponseContext) {
	rc = &ResponseContext{Errcode: result, Msg: msg}
	return
}
*/
// Error .
type Error struct {
	ErrCode int    `json:"errcode"` // 错误码  五位数字
	ErrMsg  string `json:"errmsg"`  // 错误信息

	ServiceID string `json:"serviceid,omitempty"` // 服务ID
	TracerID  string `json:"tracerid,omitempty"`  // tracerID

}

func (e Error) Error() string {

	return fmt.Sprintf("%s(%d)", e.ErrMsg, e.ErrCode) //+ " (cause: " + e.Cause + ")"
}

// New new error
func NewError(result int, msg string) *Error {

	return &Error{
		ErrCode: result,
		ErrMsg:  msg,
	}
}

// AssertError .
func MakeError(e error) (err *Error) {
	if e == nil {
		return
	}

	var tmpErr *Error
	if errors.As(e, &tmpErr) {
		err = tmpErr
		return
	}
	err = NewError(ErrorSystemFail, e.Error())
	return
}

var (
	ErrorSystemFail = 500
)
