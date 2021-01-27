package fasthttp

import (
	"encoding/json"
	//	"errors"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/gdgrc/grutils/grapps/config"
	"github.com/gdgrc/grutils/grframework"
	"github.com/valyala/fasthttp"
	"reflect"
	"runtime/debug"

	"sync"
)

//Register3 Register3
var fhrInit sync.Once
var fhr *fasthttprouter.Router

func ResponseMap(c *grframework.Context, result *interface{}, isBeauty bool) {
	var data []byte
	var err error
	data = c.GetRawResponse()
	if data != nil {
		c.FasthttpCtx.Write(data)
		config.DefaultLogger.Debug("req: %s, ===================================== rawRsp: %s", string(c.FasthttpCtx.Request.Body()), string(data))
		return
	}

	if isBeauty {
		data, err = json.MarshalIndent(*result, "", "      ")
		if err != nil {
			msg := fmt.Sprintf(`{"result":0, "msg":"last decode err"}`)
			//c.String(http.StatusOK, msg)
			c.FasthttpCtx.Write([]byte(msg))
			config.DefaultLogger.Error(msg + err.Error())
			return
		}

	} else {
		data, err = json.Marshal(*result)
		if err != nil {
			msg := fmt.Sprintf(`{"result":0, "msg":"last decode err"}`)
			//c.String(http.StatusOK, msg)
			c.FasthttpCtx.Write([]byte(msg))
			config.DefaultLogger.Error(msg + err.Error())
			return
		}

	}
	rspBody := []byte(data)

	config.DefaultLogger.Debug("req: %s, ===================================== rspBody: %s", string(c.FasthttpCtx.Request.Body()), string(rspBody))

	c.FasthttpCtx.Write(rspBody)

}

func NewGroup() {

}

func Register(funcPath string, h interface{}, middleware ...func(*grframework.Context) *grframework.Error) { //TODO: h is can be umarshal
	fhrInit.Do(func() {
		fhr = fasthttprouter.New()
	})
	v := reflect.ValueOf(h)
	t := v.Type()

	//var tmpErr error
	//tmpErr = errors.New("")

	errorType := reflect.TypeOf((*error)(nil)).Elem()
	out1Type := t.Out(0)
	if !out1Type.Implements(errorType) {
		//log.Panic("XX(*engin.Context, proto.Message)(proto.Message, error): first out arg must be proto.Message")

		msg := fmt.Sprintf("Register func wrong format,out1 element should be %s rather than %s", errorType.String(), out1Type.String())
		panic(msg)
	}

	/*out1ElemName := t.Out(0).Elem().String()
	if out1ElemName != "error" {
		msg := fmt.Sprintf("Register func wrong format,out1 element should be error, rather than %s", out1ElemName)
		panic(msg)
	}*/

	handler := func(c *fasthttp.RequestCtx) {

		reqT := t.In(1).Elem()
		rspT := t.In(2).Elem()
		reqV := reflect.New(reqT)
		rspV := reflect.New(rspT)

		var result interface{}

		defaultResult := &grframework.FrameworkResponse{}

		result = defaultResult
		defaultResult.ErrCode = 0
		defaultResult.ErrMsg = ""

		ctx := &grframework.Context{FasthttpCtx: c, Headers: make(map[string]string)}

		//ctx.SetRawResponse(nil)

		defer func() {
			if err := recover(); err != nil {
				config.DefaultLogger.Error("ProcePanic：", err, "\nfull stack: ", string(debug.Stack()))
				//commRsp.ErrMsg = CreateEnginErrWithHint(KProcessPanic, genSessionID())
				//responeRetCodeByHeader(c, commRsp.ErrMsg.ErrCode)
				//ginrender.Response(c, &commRsp)
				//collector.ReportLocalCall(commRsp.ErrMsg.ErrCode, modelName, funcPath)
			}
		}()

		defer ResponseMap(ctx, &result, false)

		f := func(key []byte, value []byte) {
			ctx.Headers[string(key)] = string(value)
		}
		ctx.FasthttpCtx.Request.Header.VisitAll(f)

		for _, m := range middleware {
			grErr := m(ctx)
			if grErr != nil {
				defaultResult.ErrCode = grErr.ErrCode
				defaultResult.ErrMsg = grErr.ErrMsg
				return
			}
		}

		if c.IsGet() {

			defaultResult.ErrCode = grframework.RESULT_FRAMEWORK_EXAMINE_FAIL
			defaultResult.ErrMsg = "should not be get method"
			return

		}

		err := json.Unmarshal(c.PostBody(), reqV.Interface()) //c.PostBody()
		if err != nil {

			defaultResult.ErrCode = grframework.RESULT_FRAMEWORK_EXAMINE_FAIL
			defaultResult.ErrMsg = fmt.Sprintf("body umarshal fail. body: %s", c.PostBody())
			return
		}

		ret := v.Call([]reflect.Value{reflect.ValueOf(ctx), reqV, rspV})
		//e := ret[0].Interface().(error)
		e := ret[0]
		if !e.IsNil() {
			tmpErr := e.Interface().(error)
			grError := grframework.MakeError(tmpErr)
			defaultResult.ErrCode = grError.ErrCode
			defaultResult.ErrMsg = grError.ErrMsg
			return
		}

		defaultResult.Data = rspV.Interface()
		config.DefaultLogger.Debug("rspV.Interface(): ", rspV.Interface())
		//result = rspV.Interface() //ret[0].Interface()

		c.Response.Header.Add("Content-Type", "application/json")
		//ResponseMap(c, result, false)

	}
	fhr.POST(funcPath, handler)

}
