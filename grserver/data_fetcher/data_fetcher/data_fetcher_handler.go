package main

import (
	"database/sql"
	"fmt"
	"github.com/gdgrc/grutils/grapps/config/log"
	"github.com/gdgrc/grutils/grdatabase"
	"github.com/gdgrc/grutils/grframework"
	"github.com/gdgrc/grutils/grserver/data_fetcher/data_fetcherconf"
	model "github.com/gdgrc/grutils/grserver/data_fetcher/data_fetchermodel"
	"github.com/gdgrc/grutils/grserver/data_fetcher/service"
	"math"
	//"data_fetcher/pb/data_fetcher"
)

// FetchData FetchData
func FetchData(c *grframework.Context, req *model.FetchDataReq, rsp *model.FetchDataRsp) (err error) {
	//rsp = &model.FetchDataRsp{}

	// page size should not be empty
	if req.PageSize == 0 {
		err = fmt.Errorf("params error")
		return
	}

	if req.PageNo <= 1 {
		req.PageNo = 1
	}

	dataName := req.DataName

	dataConf, ok := data_fetcherconf.GlobalDataFetcherConf.Querys[dataName]
	if !ok {
		log.Warn("Can not find dataname: ", dataName)
		err = fmt.Errorf("find data but match conf failed")
		//rc = grframework.MakeResultWithMsg(-1, "find data but match conf failed")
		return
	}

	err = ConstructAndSendDatabaseRequest(req, &dataConf, rsp)
	if err != nil {
		log.Error("ConstructAndSendDatabaseRequest fail, msg: ", err)

		return
	}

	log.Debug(fmt.Sprintf("%+v, %+v", req, rsp))

	return
}

func ConstructAndSendDatabaseRequest(req *model.FetchDataReq, dataConf *data_fetcherconf.Query, rsp *model.FetchDataRsp) (err error) {

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
	mainStatement, err := service.ConstructMainStatment(req, dataConf)
	if err != nil {
		log.Error("ConstructMainStatment return err: ", err)
		return
	}

	countRows, err := databaseConn.Query(mainStatement.GetCountPreparedStatement(), mainStatement.GetParams()...)
	if err != nil {
		return
	}

	defer countRows.Close()

	for countRows.Next() {

		//fmt.Println("goood")

		err = countRows.Scan(&rsp.TotalRecordNum)
		if err != nil {

			return
		}
		break

	}

	rsp.PageSize = req.PageSize
	rsp.PageNo = req.PageNo
	rsp.TotalPageNum = int(math.Ceil(float64(rsp.TotalRecordNum) / float64(rsp.PageSize)))

	dataRows, err := databaseConn.Query(mainStatement.GetRecordPreparedStatement(), mainStatement.GetParams()...)
	if err != nil {
		return
	}

	defer dataRows.Close()

	columns, err := dataRows.Columns()
	if err != nil {

		return
	}

	rsp.RecordList = make([]map[string]string, 0)
	values := make([]sql.NullString, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for dataRows.Next() {

		//scanArgs := make([]interface{}, len(columns))

		err = dataRows.Scan(scanArgs...)
		if err != nil {

			return
		}

		data := make(map[string]string)
		for i, v := range values {
			if v.Valid {
				data[columns[i]] = v.String
			} else {
				data[columns[i]] = ""
			}
		}

		rsp.RecordList = append(rsp.RecordList, data)

	}
	rsp.RecordNum = len(rsp.RecordList)
	rsp.PageNo = req.PageNo

	return

}
