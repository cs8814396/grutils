package handler

import (
	//"device_filter/reporter/conf"

	//"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"fmt"

	"github.com/gdgrc/grutils/grframework"
	//"github.com/mongodb/mongo-go-driver/core/result"

	//"math"
	"device_filter/reporter/model"

	//dfModel "github.com/gdgrc/grutils/grserver/data_fetcher/model"
	"encoding/json"
	//"os"
	//"path"
	"time"
	//"sync"
	"github.com/gdgrc/grutils/grapps/config/log"
	//dfClient "github.com/gdgrc/grutils/grserver/data_fetcher/client"
)

type CommReportReq struct {
	Datas []*model.CommReportData `json:"datas"`
	Event string                  `json:"event"`
}
type CommReportRsp struct {
	//	Ids []string `json:"ids"`
	//Result   int    `json:"result"`
	//ExtraMsg string `json:"extraMsg"`
}

func CommResponse(c *grframework.Context, inErr *error, data interface{}) {
	dataMap := make(map[string]interface{})
	dataMap["result"] = 1
	dataMap["extraMsg"] = ""
	if *inErr != nil {
		dataMap["result"] = -1
		dataMap["extraMsg"] = (*inErr).Error()
	} else {
		dataMap["results"] = data
	}

	b, err := json.Marshal(dataMap)
	if err != nil {
		log.Error("Response fail, err: ", err)
	}
	log.Info("SetRawResponse:  ", string(b))
	c.SetRawResponse(b)
}

// FetchData FetchData
func CommReport(c *grframework.Context, req *CommReportReq, rsp *ReportRsp) (err error) {

	// page size should not be empty
	defer CommResponse(c, &err, rsp)

	//req.Event = "device_info"
	channel, err := getChannelByEvent(req.Event)
	if err != nil {
		log.Error("getChannelByEvent fail, err: ", err)
		//rc = grframework.MakeResultWithMsg(-1, "event get fail")
		return
	}
	log.Debug("CommReport get data event: %s length: %d", req.Event, len(req.Datas))

	for _, d := range req.Datas {
		d.Event = req.Event

		timeOk := false
		timeInterface, ok := d.Data["time"]

		if ok {
			timeString := fmt.Sprintf("%v", timeInterface)
			_, err = time.ParseInLocation("2006-01-02 15:04:05", timeString, time.Local)
			if err == nil {
				timeOk = true
			}
			d.Time = timeString
		}
		if !timeOk {
			d.Time = time.Now().Format("2006-01-02 15:04:05")
		}

		channel <- d
	}

	return
}
