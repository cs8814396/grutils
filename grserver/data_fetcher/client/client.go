package client

import (
	//"encoding/json"
	"github.com/gdgrc/grutils/grframework"
	dfModel "github.com/gdgrc/grutils/grserver/data_fetcher/model"
)

func CallDataFetcher(url string, req interface{}, rsp interface{}) (err error) {

	err = grframework.InternalCall(url, req, rsp)

	return

}
func Insert(dataName string, rl [][]interface{}) (err error) {

	req := dfModel.InsertDataReq{
		DataName:   dataName,
		RecordList: rl,
	}

	rsp := dfModel.InsertDataRsp{}
	url := "http://127.0.0.1:9096/insert_data"
	err = CallDataFetcher(url, req, rsp)
	return
}

func Query(req *dfModel.FetchDataReq, rsp *dfModel.FetchDataClientRsp) (err error) {

	url := "http://127.0.0.1:9096/fetch_data"
	err = CallDataFetcher(url, req, rsp)
	return
}
