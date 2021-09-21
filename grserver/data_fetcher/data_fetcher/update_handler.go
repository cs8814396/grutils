package main

import (
	"data_fetcher/data_fetcherconf"
	model "data_fetcher/model"
	"data_fetcher/service"

	"github.com/gdgrc/grutils/grdatabase"

	//"data_fetcher/service"
	//"database/sql"
	"fmt"

	"github.com/gdgrc/grutils/grapps/config/log"
	"github.com/gdgrc/grutils/grframework"
	//	"math"
	//"data_fetcher/pb/data_fetcher"
)

func UpdateData(c *grframework.Context, req *model.UpdateDataReq, rsp *model.UpdateDataRsp) (err error) {
	//	rsp = &model.InsertDataRsp{}

	// page size should not be empty

	dataName := req.DataName

	dataConf, ok := data_fetcherconf.GlobalDataFetcherConf.Updates[dataName]
	if !ok {
		log.Warn("Can not find dataname: ", dataName)
		err = fmt.Errorf("find data but match conf failed")
		return
	}
	log.Debug(fmt.Sprintf("dataConf: %+v", dataConf))
	err = SendDatabaseUpdateRequest(req, &dataConf, rsp)
	if err != nil {
		log.Error("ConstructAndSendDatabaseRequest fail, msg: ", err)
		//rc = grframework.MakeResultWithMsg(-1, "data execute fail "+err.Error())
		return
	}

	log.Debug(fmt.Sprintf("%+v, %+v", req, rsp))

	return
}
func SendDatabaseUpdateRequest(req *model.UpdateDataReq, dataConf *data_fetcherconf.Update, rsp *model.UpdateDataRsp) (err error) {

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

	mainStatement, err := service.ConstructMainStatment(req.Condition, dataConf.Statement, dataConf.Conditions)
	if err != nil {
		log.Error("ConstructMainStatment return err: ", err)
		return
	}

	log.Debug("update conf: %+v, statement: %s, mainStatement.GetParams(): %+v", dataConf, mainStatement.PreparedStatement, mainStatement.GetParams())

	result, err := databaseConn.Exec(mainStatement.PreparedStatement, mainStatement.GetParams()...)
	if err != nil {
		log.Error("update Exec return err: ", err)
		return
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("update RowsAffected err: ", err)
		return
	}
	rsp.RowAffected = rowAffected
	lastId, err := result.LastInsertId()
	if err != nil {
		log.Error("update LastInsertId err: ", err)
		return
	}

	rsp.IdList = append(rsp.IdList, lastId)

	return

}
