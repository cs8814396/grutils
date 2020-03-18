package grframework

import (
	"github.com/valyala/fasthttp"
)

type RequestContext struct {
	FasthttpCtx *fasthttp.RequestCtx
}

type ResponseContext struct {
	Headers map[string]string
	Result  int
	Msg     string
}

func MakeResultWithMsg(result int, msg string) (rc *ResponseContext) {
	rc = &ResponseContext{Result: result, Msg: msg}
	return
}
