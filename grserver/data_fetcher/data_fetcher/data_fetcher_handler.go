package main

import (
	//"database/sql"
	"fmt"
	"math"
	"reflect"

	"github.com/gdgrc/grutils/grapps/config/log"
	"github.com/gdgrc/grutils/grdatabase"
	"github.com/gdgrc/grutils/grframework"
	"github.com/gdgrc/grutils/grserver/data_fetcher/data_fetcherconf"
	model "github.com/gdgrc/grutils/grserver/data_fetcher/model"
	"github.com/gdgrc/grutils/grserver/data_fetcher/service"

	//"strings"
	"encoding/json"
	"time"
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

	err = SendDatabaseQueryRequest(req, &dataConf, rsp)
	if err != nil {
		log.Error("SendDatabaseQueryRequest fail, msg: ", err)

		return
	}

	log.Debug(fmt.Sprintf("FetchData req: %+v, rsp: %+v", req, rsp))

	return
}

func SendDatabaseQueryRequest(req *model.FetchDataReq, dataConf *data_fetcherconf.Query, rsp *model.FetchDataRsp) (err error) {

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

	log.Debug("conf: %+v,mainStatement.GetRecordPreparedStatement(): %s, mainStatement.GetParams(): %+v", dataConf, mainStatement.GetRecordPreparedStatement(req.PageNo, req.PageSize), mainStatement.GetParams())

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
	dataRows, err := databaseConn.Query(mainStatement.GetRecordPreparedStatement(req.PageNo, req.PageSize), mainStatement.GetParams()...)
	if err != nil {
		return
	}

	defer dataRows.Close()

	columns, err := dataRows.Columns()
	if err != nil {

		return
	}

	columnTypes, err := dataRows.ColumnTypes()
	if err != nil {
		return
	}

	rl := make([]map[string]interface{}, 0)
	values := make([]interface{}, len(columns), len(columns))
	scanArgs := make([]interface{}, len(columns), len(columns))

	for i := 0; i < len(scanArgs); i++ {
		//var valueInterface interface{}
		/*
			dtn := strings.ToUpper(columnTypes[i].DatabaseTypeName())
			if dtn == "JSON" {

				//valueInterface = map[string]interface{}{}

				scanArgs[i] = map[string]interface{}{}
			} else if strings.Index(dtn, "INT") > -1 {
				var t int64
				valueInterface = t
			} else if strings.Index(dtn, "CHAR") > -1 || strings.Index(dtn, "TIME") > -1 {
				var t string
				valueInterface = t
			} else {
				err = fmt.Errorf("UnRecognize DatabaseType: %s", dtn)
				return
			}*/

		scanArgs[i] = &values[i]
	}

	for dataRows.Next() {

		//scanArgs := make([]interface{}, len(columns))

		err = dataRows.Scan(scanArgs...)
		if err != nil {

			return
		}

		data := make(map[string]interface{})
		for i, v := range values {
			/*
				if v.Valid {
					data[columns[i]] = v.String()
				} else {
					data[columns[i]] = ""
				}*/
			vType := reflect.TypeOf(v).String()
			switch vType {

			case "[]uint8":

				if columnTypes[i].DatabaseTypeName() == "JSON" {
					tmpValue := make(map[string]interface{})
					err = json.Unmarshal([]byte(v.([]uint8)), &tmpValue)
					if err != nil {
						return
					}
					data[columns[i]] = tmpValue
				} else {
					//tmpValue := string(v.([]uint8))
					data[columns[i]] = string(v.([]uint8))
				}
			case "int64":
				data[columns[i]] = v
			case "time.Time":
				data[columns[i]] = (v.(time.Time)).Format("2006-01-02 15:04:05")
			default:
				err = fmt.Errorf("unrecognize type: %s", vType)
				return
			}
			//data[columns[i]] = v

			//log.Debug("column: %s type: %s, data: %s,columnScanType: %s, columnTypeName: %s,DatabaseTypeName: %s", columns[i], reflect.TypeOf(v).String(), data[columns[i]], columnTypes[i].ScanType(), columnTypes[i].Name(), columnTypes[i].DatabaseTypeName())

		}

		rl = append(rl, data)

	}

	rsp.RecordNum = len(rl)
	rsp.PageNo = req.PageNo
	rsp.RecordList = rl

	return

}
