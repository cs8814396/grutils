package grframework

import (
	"encoding/json"
	"fmt"

	//"github.com/gdgrc/grutils/grfile"
	"github.com/gdgrc/grutils/grnetwork"
)

const (
	HEADER_CALL_WAY             = "call_way"
	HEADER_VALUES_INTERNAL_CALL = "internal_call"
)

type FrameworkResponse struct {
	ErrCode int         `json:"errcode"`
	ErrMsg  string      `json:"errmsg"`
	Data    interface{} `json:"data"`
}

func InternalCall(url string, req interface{}, rsp interface{}) (err error) {

	b, err := json.Marshal(req)

	header := make(map[string]string)
	header[HEADER_CALL_WAY] = HEADER_VALUES_INTERNAL_CALL

	rspBytes, _, err := grnetwork.HttpPostRetry(url, b, header, 50000, 1)
	if err != nil {
		return
	}

	defaultResult := &FrameworkResponse{Data: rsp}
	err = json.Unmarshal(rspBytes, defaultResult)
	if err != nil {
		err = fmt.Errorf("InternalCall Rsp unmarshal fail: %s, RspData: %s", err.Error(), string(rspBytes))
		return
	}
	if defaultResult.ErrCode != 0 {
		err = Error{ErrCode: defaultResult.ErrCode, ErrMsg: defaultResult.ErrMsg}
		return
	}

	return

}
