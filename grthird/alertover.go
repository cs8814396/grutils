package grthird

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gdgrc/grutils/grnetwork"
	"time"
)

var monitTimeStamp int64

type AlertOverReq struct {
	Source   string `json:"source"`
	Receiver string `json:"receiver"`
	Content  string `json:"content"`
	Title    string `json:"title"`
}

type AlertOverRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

const (
	ALERTOVER_SOUND_SILENT     = "silent"
	ALERTOVER_URL              = "https://api.alertover.com/v1/alert"
	ALERTOVER_POST_TIMEOUT     = 3000
	ALERTOVER_DEFAULT_DURATION = 10
)

type XmlAlertOver struct {
	Source   string `xml:"source"`
	Receiver string `xml:"receiver"`
}

func AlertOverNotify(source string, receiver string, title string, content string) (err error) {

	now := time.Now().Unix()
	/*
		lgd.Trace("monit notify [%s], now[%d] monitLastTime[%d]", game, now, monitTimeStamp)
		if now-monitTimeStamp < config.GlobalConf.Monit.Duration {
			return true
		}
	*/
	if now-monitTimeStamp < ALERTOVER_DEFAULT_DURATION { // 10s
		return
	}

	monitTimeStamp = now
	payload := AlertOverReq{
		Source:   source,
		Receiver: receiver,
		Content:  content,
		Title:    title,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	heads := make(map[string]string)
	heads["content-type"] = "application/json"

	rspStr, _, err := grnetwork.HttpPostExTimeout(ALERTOVER_URL, string(data), heads, ALERTOVER_POST_TIMEOUT)
	if err != nil {
		return
	}

	var rsp AlertOverRsp
	err = json.Unmarshal([]byte(rspStr), &rsp)
	if err != nil || rsp.Code != 0 {
		msg := fmt.Sprintf("monit %+v json unmarshal fail %s rsp %s", payload, err, rspStr)
		err = errors.New(msg)
		return
	}
	return
}

func init() {
	monitTimeStamp = time.Now().Unix() - ALERTOVER_DEFAULT_DURATION
}
