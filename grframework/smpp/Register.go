package smpp

import (
	"encoding/json"
	"fmt"
	"github.com/ajankovic/smpp"
	"github.com/ajankovic/smpp/pdu"
	"github.com/gdgrc/grutils/grapps/config"
	"github.com/gdgrc/grutils/grcommon"
	"github.com/gdgrc/grutils/grframework"
	"net/url"
	"reflect"
	"runtime/debug"
	"sync"
)

/*

fhrInit.Do(func() {
		fhr = fasthttprouter.New()
	})

*/

var smppInit sync.Once

var sessionConf *smpp.SessionConf

var server *smpp.Server

type SmppServer struct {
	*smpp.Server
}

func (s *SmppServer) ListenAndBlock(addr string) {
	server = smpp.NewServer(addr, *sessionConf)

	err := server.ListenAndServe()

	panic(err)

}

var SmppContextKey = "smppContext"
var SmppSubmitSmContextKey = "smppSubmitSmContext"

type SubmitSm pdu.SubmitSm

func GetSubmitSmContext(c *grframework.Context) *SubmitSm {
	return (c.GetData(SmppSubmitSmContextKey)).(*SubmitSm)

}

func (s *SmppServer) Register(funcPath string, h interface{}) {
	smppInit.Do(func() {
		sessionConf = &smpp.SessionConf{}
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

	sessionConf.Handler = smpp.HandlerFunc(func(ctx *smpp.Context) {

		defer func() {
			if err := recover(); err != nil {
				config.DefaultLogger.Error("ProcePanic：", err, "\nfull stack: ", string(debug.Stack()))
				//commRsp.ErrMsg = CreateEnginErrWithHint(KProcessPanic, genSessionID())
				//responeRetCodeByHeader(c, commRsp.ErrMsg.ErrCode)
				//ginrender.Response(c, &commRsp)
				//collector.ReportLocalCall(commRsp.ErrMsg.ErrCode, modelName, funcPath)
			}
		}()

		reqT := t.In(1).Elem()
		rspT := t.In(2).Elem()
		reqV := reflect.New(reqT)
		rspV := reflect.New(rspT)

		var result interface{}

		defaultResult := &grframework.FrameworkResponse{}

		result = defaultResult
		defaultResult.ErrCode = 0
		defaultResult.ErrMsg = ""

		fwCtx := &grframework.Context{RequestCtx: new(grframework.RequestContext)}

		fwCtx.SetData(SmppContextKey, ctx)

		switch ctx.CommandID() {
		case pdu.BindTransceiverID:
			btrx, err := ctx.BindTRx()
			if err != nil {
				config.DefaultLogger.Error("err: ", err)
			}
			config.DefaultLogger.Debug(fmt.Sprintf("BindTRx: btrx: %+v", btrx))
			resp := btrx.Response("good")
			if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
				config.DefaultLogger.Error("err: ", err)
			}
		case pdu.SubmitSmID:

			sm, err := ctx.SubmitSm()
			if err != nil {
				config.DefaultLogger.Error("err: ", err)
			}
			fwCtx.SetData(SmppSubmitSmContextKey, (*SubmitSm)(sm))

			//fwCtx.SetData(smppContextKey, sm)

			defer ResponseMap(fwCtx, &result, sm.ShortMessage, ctx.CommandID(), sm, false)

			gbk, _ := grcommon.GBKToUTF8([]byte(sm.ShortMessage))
			iso88598 := string(grcommon.ISO88598ToUtf8([]byte(sm.ShortMessage)))
			ucs2 := grcommon.UCS2([]byte(sm.ShortMessage))
			ucs2message := string(ucs2.Decode())

			config.DefaultLogger.Debug(fmt.Sprintf("SubmitSmID: sm: %+v, iso88598: %s gbk %s, ucs2message: %s", sm, string(iso88598), gbk, ucs2message))

			//config.DefaultLogger.Debug("get short message: ", sm.ShortMessage, " prototype: ", sm.ProtocolID, " DataCoding: ", sm.DataCoding)
			sm.ShortMessage = ucs2message
			switch sm.DataCoding {
			case 0x1001:
				var urlValues url.Values
				urlValues, err = url.ParseQuery(sm.ShortMessage)
				if err != nil {
					defaultResult.ErrCode = grframework.RESULT_FRAMEWORK_EXAMINE_FAIL
					defaultResult.ErrMsg = fmt.Sprintf("body ParseQuery fail. body: %s", sm.ShortMessage)
					return
				}

				tmpMap := make(map[string]string)

				for key, values := range urlValues {
					if len(values) > 0 {
						tmpMap[key] = values[0]
					}

				}
				var b []byte
				b, err = json.Marshal(tmpMap)
				if err != nil {
					defaultResult.ErrCode = grframework.RESULT_FRAMEWORK_EXAMINE_FAIL
					defaultResult.ErrMsg = fmt.Sprintf("body Marshal fail. body: %s", sm.ShortMessage)
					return
				}

				err = json.Unmarshal(b, reqV.Interface())
				if err != nil {
					defaultResult.ErrCode = grframework.RESULT_FRAMEWORK_EXAMINE_FAIL
					defaultResult.ErrMsg = fmt.Sprintf("body umarshal fail. body: %s", sm.ShortMessage)
					return
				}

			case 0x1002:
				err := json.Unmarshal([]byte(sm.ShortMessage), reqV.Interface()) //c.PostBody()
				if err != nil {
					defaultResult.ErrCode = grframework.RESULT_FRAMEWORK_EXAMINE_FAIL
					defaultResult.ErrMsg = fmt.Sprintf("body umarshal fail. body: %s", sm.ShortMessage)
					return
				}
			default:
				//raw text

			}

			ret := v.Call([]reflect.Value{reflect.ValueOf(fwCtx), reqV, rspV})
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

			//msgID++

		case pdu.UnbindID:
			unb, err := ctx.Unbind()
			if err != nil {
				config.DefaultLogger.Error("err: ", err)
			}
			config.DefaultLogger.Debug(fmt.Sprintf("UnbindID: unb: %+v", unb))
			resp := unb.Response()
			if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
				config.DefaultLogger.Error("err: ", err)
			}
			ctx.CloseSession()
		}

	})

}
func Response(c *grframework.Context, commandId pdu.CommandID, rawCtx interface{}, data string) {

	switch commandId {
	case pdu.SubmitSmID:
		sm := rawCtx.(*pdu.SubmitSm)

		resp := sm.Response(data)
		if err := ((c.GetData(SmppContextKey)).(*smpp.Context)).Respond(resp, pdu.StatusOK); err != nil {
			config.DefaultLogger.Error("err: ", err, " commandId: ", commandId.String(), " sm: ", sm)
			return
		}
	}

}
func ResponseMap(c *grframework.Context, result *interface{}, reqBody string, commandId pdu.CommandID, rawCtx interface{}, isBeauty bool) {

	var data []byte
	var err error
	data = c.GetRawResponse()
	if data != nil {
		Response(c, commandId, rawCtx, string(data))
		return
	}

	if isBeauty {
		data, err = json.MarshalIndent(*result, "", "      ")
		if err != nil {
			msg := fmt.Sprintf(`{"result":0, "msg":"last decode err"}`)
			//c.String(http.StatusOK, msg)
			Response(c, commandId, rawCtx, msg)
			config.DefaultLogger.Error(msg + err.Error())
			return
		}

	} else {
		data, err = json.Marshal(*result)
		if err != nil {
			msg := fmt.Sprintf(`{"result":0, "msg":"last decode err"}`)
			Response(c, commandId, rawCtx, msg)
			config.DefaultLogger.Error(msg + err.Error())
			return
		}

	}
	rspBody := string(data)

	config.DefaultLogger.Debug("Smpp Finish req: %s , rsp: %s ", reqBody, rspBody)

	Response(c, commandId, rawCtx, rspBody)

	return
}
