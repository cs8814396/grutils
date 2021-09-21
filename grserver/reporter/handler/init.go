package handler

import (
	"device_filter/reporter/conf"
	"device_filter/reporter/model"
	//"encoding/json"

	//"encoding/json"
	"fmt"
	"github.com/gdgrc/grutils/grapps/config/log"
	//"github.com/gdgrc/grutils/grnetwork"
	dfClient "github.com/gdgrc/grutils/grserver/data_fetcher/client"
	"os"
	"path"
	"sync"
	"time"
	//	"github.com/gogo/protobuf/test/data"
)

var channelMap sync.Map
var channelMutex sync.Mutex

func CloseFd(f **os.File) {

	if f != nil {
		(*f).Close()
		(*f) = nil
	}
}
func CreateFd(dirPath string, name string) (f *os.File, err error) {

	if err = os.MkdirAll(dirPath, 0755); err != nil {
		err = fmt.Errorf("mkdir: %s err: %s", dirPath, err)
		return
	}
	filePath := path.Join(dirPath, name)
	f, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		err = fmt.Errorf("open %v fail: %v", dirPath, err.Error())

		return
	}
	return
}

func getChannelByEvent(event string) (channel chan model.ReportData, err error) {
	ch, ok := channelMap.Load(event)

	if !ok {

		err = fmt.Errorf("Can not found this channel event: %s", event)
		return
	}

	channel = ch.(chan model.ReportData)

	return
}

func Init() bool {

	for _, eventName := range conf.GlobalReporterConf.Events {

		ch := make(chan model.ReportData, 65535)
		channelMap.Store(eventName, ch)
		log.Info("Init event channel: %s ", eventName)
		go WriteLog(eventName, ch)

	}

	return true

}
func FileWriter(dataText []byte, event string) {

}

func WriteLog(event string, channel chan model.ReportData) {
	openTime := int64(0)
	var logFile *os.File

	LogDir := conf.GlobalReporterConf.LogDir
	defer CloseFd(&logFile)

	LogPath := ""
	fileName := ""

	timeoutTicker := time.NewTicker(time.Second * 1) // 每隔1s进行一次打印
	for {

		now := time.Now() //TODO: reduce the system call time of now

		dataMapSlice := make([]model.ReportData, 0, 0)
		count := 0
	Loop:
		for {
			select {
			case <-timeoutTicker.C:
				break Loop
			case dataMap := <-channel:

				dataMap.SetDataTime(now.Format("2006-01-02 15:04:05"))

				dataMapSlice = append(dataMapSlice, dataMap)
				count = count + 1
				//log.Debug("Get Data: %+v ", dataMap, " count: ", count)
				if count >= 10 {

					break Loop
				}
			}
		}
		//log.Debug("begin to write: %+v ", dataMapSlice)

		if len(dataMapSlice) <= 0 {
			continue
		}

		dataName:=dataMapSlice[0].GetDataName()
		//过滤掉可能导致sql注入的字符

		//发送日志json数组到适配器

		//continue
		var dataSlice [][]interface{}
		var dataText []byte

		isContainedValidData := false
		for _, dataMap := range dataMapSlice {
			err := dataMap.Prepare()
			if err != nil {
				log.Error("Data Depreated! data prepare fail: %s, dataMap: %+v", err.Error(), dataMap)
				continue
			}
			text, e := dataMap.DumpBytes()
			if e != nil {
				log.Error("Data Depreated! Unexpected behavior while DumpBytes!!! ", e, " dataMap: ", dataMap)
				continue
			}
			dataText = append(dataText, text...)
			dataText = append(dataText, byte('\n'))
			//----
			rl, err := dataMap.DumpOrderedList()
			if err!=nil{
				log.Error("Data Depreated! Unexpected behavior while DumpOrderedList!!! ", e, " dataMap: ", dataMap)
				continue

			}
			if len(rl)>0{
				dataSlice = append(dataSlice, rl)
			}
			
			isContainedValidData = true
		}



		if !isContainedValidData {
			log.Error("Contain invalid Data and be depreacated ")
			continue
		}

		if now.Unix()/3600 != openTime/3600 {
			CloseFd(&logFile)
		}
		if logFile == nil {
			var err error
			LogPath = path.Join(LogDir, event)
			fileName = event + "_" + now.Format("2006010215.log")
			logFile, err = CreateFd(LogPath, fileName)
			if err != nil {
				log.Error("CreateFd fail: ", err)
			}

			openTime = now.Unix()
		}

		logFile.Write(dataText)

		if len(dataSlice)>0{
			if err := dfClient.Insert(dataName, dataSlice); err != nil {
				log.Error("Unexpected behavior while insert!!! ", err, dataSlice)
				continue
			}
		}

		

	}
}

