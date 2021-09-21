package main

import (
	"database/sql"

	"data_fetcher/data_fetcherconf"
	model "data_fetcher/model"

	//"data_fetcher/service"
	//"database/sql"
	"fmt"

	"github.com/gdgrc/grutils/grapps/config/log"
	"github.com/gdgrc/grutils/grdatabase"
	"github.com/gdgrc/grutils/grframework"
	//	"math"
	//"data_fetcher/pb/data_fetcher"
)

// FetchData FetchData
func InsertData(c *grframework.Context, req *model.InsertDataReq, rsp *model.InsertDataRsp) (err error) {
	//	rsp = &model.InsertDataRsp{}

	// page size should not be empty
	if len(req.RecordList) == 0 {
		err = fmt.Errorf("params error, record_list is 0 length")
		return
	}

	dataName := req.DataName

	dataConf, ok := data_fetcherconf.GlobalDataFetcherConf.Inserts[dataName]
	if !ok {
		log.Warn("Can not find dataname: ", dataName)
		err = fmt.Errorf("find data but match conf failed")
		return
	}
	log.Debug(fmt.Sprintf("dataConf: %+v", dataConf))
	err = SendDatabaseInsertRequest(req, &dataConf, rsp)
	if err != nil {
		log.Error("ConstructAndSendDatabaseRequest fail, msg: ", err)
		//rc = grframework.MakeResultWithMsg(-1, "data execute fail "+err.Error())
		return
	}

	log.Debug(fmt.Sprintf("%+v, %+v", req, rsp))

	return
}

func SendDatabaseInsertRequest(req *model.InsertDataReq, dataConf *data_fetcherconf.Insert, rsp *model.InsertDataRsp) (err error) {

	instanceName := dataConf.DatabaseInstance
	databaseName := dataConf.DatabaseName
	instance, ok := data_fetcherconf.GlobalDataFetcherConf.Instances[instanceName]
	if !ok {
		err = fmt.Errorf("can not find this instance: %s", instanceName)
		return
	}

	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&timeout=30s&loc=Local&autocommit=true&parseTime=true",
		instance.Username, instance.Password, instance.Ip, instance.Port, databaseName)

	maxIdleConn := 10
	databaseConn, err := grdatabase.DefaultMysqlPool.DBGetConn(instanceName, dsn, maxIdleConn)
	if err != nil {
		log.Error("DBGetConn return err: ", err)
		return
	}

	conn, err := databaseConn.Begin()
	if err != nil {
		log.Error("begin return err: ", err)
		return
	}
	defer conn.Rollback()
	prepareStatement, err := conn.Prepare(dataConf.Statement)
	if err != nil {
		return
	}

	defer prepareStatement.Close()

	for _, data := range req.RecordList {
		var result sql.Result
		result, err = prepareStatement.Exec(data...)
		if err != nil {
			log.Error("data exec fail: %+v", data)
			return

		}
		var lastId int64
		lastId, err = result.LastInsertId()
		if err != nil {
			log.Error("data exec LastInsertId: %s", err)
			return
		}
		rsp.IdList = append(rsp.IdList, lastId)

	}
	conn.Commit()

	return

}
