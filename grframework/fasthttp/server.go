package fasthttp

import (
	//"github.com/gdgrc/grutils/grapps/config"
	"github.com/valyala/fasthttp"
	//	"runtime/debug"
	"sync"
)

var fhsInit sync.Once
var fhs *fasthttp.Server

func ListenAndBlock(addr string) error {

	fhsInit.Do(func() { fhs = &fasthttp.Server{} })

	fhs.Handler = fhr.Handler
	return fhs.ListenAndServe(addr)
}
