package fasthttp

import (
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/gdgrc/grutils/grapps/config"
	"github.com/gdgrc/grutils/grframework"
	"github.com/valyala/fasthttp"

	"reflect"
	"sync"
)

//Register3 Register3
var fhrInit sync.Once
var fhr *fasthttprouter.Router

func ResponseMap(c *fasthttp.RequestCtx, result *interface{}, isBeauty bool) {
	var data []byte
	var err error
	if isBeauty {
		data, err = json.MarshalIndent(*result, "", "      ")
		if err != nil {
			msg := fmt.Sprintf(`{"result":0, "msg":"last decode err"}`)
			//c.String(http.StatusOK, msg)
			c.Write([]byte(msg))
			config.DefaultLogger.Error(msg + err.Error())
			return
		}

	} else {
		data, err = json.Marshal(*result)
		if err != nil {
			msg := fmt.Sprintf(`{"result":0, "msg":"last decode err"}`)
			//c.String(http.StatusOK, msg)
			c.Write([]byte(msg))
			config.DefaultLogger.Error(msg + err.Error())
			return
		}

	}
	rspBody := []byte(data)

	config.DefaultLogger.Debug("req: %s, rspBody: %s", string(c.Request.Body()), rspBody)

	c.Write(rspBody)

}
func Register(funcPath string, h interface{}) {
	fhrInit.Do(func() {
		fhr = fasthttprouter.New()
	})
	v := reflect.ValueOf(h)
	t := v.Type()
	out1ElemName := t.Out(1).Elem().String()
	if out1ElemName != "grframework.ResponseContext" {
		msg := fmt.Sprintf("Register func wrong format,out1 element should be ResponseContext, rather than %s", out1ElemName)
		panic(msg)
	}

	handler := func(c *fasthttp.RequestCtx) {
		reqT := t.In(1).Elem()
		reqV := reflect.New(reqT)

		var result interface{}

		defaultResult := map[string]interface{}{}

		result = defaultResult
		defaultResult[grframework.RESULT] = 0
		defaultResult[grframework.MSG] = ""

		defer ResponseMap(c, &result, false)

		ctx := &grframework.RequestContext{FasthttpCtx: c}

		if c.IsGet() {

			defaultResult[grframework.RESULT] = grframework.RESULT_FRAMEWORK_EXAMINE_FAIL
			defaultResult[grframework.MSG] = "should not be get method"
			return

		}

		err := json.Unmarshal(c.PostBody(), reqV.Interface()) //c.PostBody()
		if err != nil {

			defaultResult[grframework.RESULT] = grframework.RESULT_FRAMEWORK_EXAMINE_FAIL
			defaultResult[grframework.MSG] = fmt.Sprintf("body umarshal fail. body: %s", c.PostBody())
			return
		}

		ret := v.Call([]reflect.Value{reflect.ValueOf(ctx), reqV})
		e := ret[1].Interface().(*grframework.ResponseContext)
		if e != nil {
			defaultResult[grframework.RESULT] = e.Result
			defaultResult[grframework.MSG] = e.Msg
			return
		}
		result = ret[0].Interface()

		c.Response.Header.Add("Content-Type", "application/json")
		//ResponseMap(c, result, false)

	}
	fhr.POST(funcPath, handler)

}
