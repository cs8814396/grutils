package grframework

import (
	"errors"
	"github.com/valyala/fasthttp"
)

/*
type RequestContext struct {
	FasthttpCtx *fasthttp.RequestCtx
}

type ResponseContext struct {
	Headers map[string]string
	Result  int
	Msg     string
}*/

type Context struct {
	FasthttpCtx *fasthttp.RequestCtx
	Headers     map[string]string
}

/*
func MakeResultWithMsg(result int, msg string) (rc *ResponseContext) {
	rc = &ResponseContext{Result: result, Msg: msg}
	return
}
*/
// Error .
type Error struct {
	Result int    `json:"result"` // 错误码  五位数字
	Msg    string `json:"msg"`    // 错误信息

	ServiceID string `json:"serviceid,omitempty"` // 服务ID
	TracerID  string `json:"tracerid,omitempty"`  // tracerID

}

func (e Error) Error() string {

	return e.Msg //+ " (cause: " + e.Cause + ")"
}

// New new error
func NewError(result int, msg string) *Error {

	return &Error{
		Result: result,
		Msg:    msg,
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
