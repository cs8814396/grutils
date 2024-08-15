package grdatabase

import (
	sqllib "database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"log"
)

func InterfaceToSqlParam(dataStruct interface{}, fields Fields) (valueList []interface{}, fieldNameList []string, err error) {
	sv := reflect.ValueOf(dataStruct)
	st := reflect.TypeOf(dataStruct)
	if st.Kind() == reflect.Ptr {
		sv = sv.Elem()
		st = st.Elem()
	}

	switch st.Kind() {
	case reflect.Map:
		/*
			CREATE TABLE `CHINA_KUAISHOU_3Y1` (
				`Id` int(11) NOT NULL AUTO_INCREMENT,
				`ANAME` varchar(128) DEFAULT NULL,
				`cell` varchar(128) DEFAULT NULL,
				`ADDRESS` varchar(255) DEFAULT NULL,
				`GOODS` varchar(255) DEFAULT NULL,
				`AMOUNT` varchar(255) DEFAULT NULL,
				`STIME` varchar(64) DEFAULT NULL,
				PRIMARY KEY (`Id`),
				KEY `IDX_GOODS` (`GOODS`(6))
			) ENGINE=InnoDB AUTO_INCREMENT=51777837 DEFAULT CHARSET=utf8mb4
			2024/03/05 15:48:05 value: 28099 f.name: Id type: string
			2024/03/05 15:48:05 value: <invalid reflect.Value> f.name: NAME type: invalid
		*/
		for _, f := range fields {
			fieldNameList = append(fieldNameList, f.Name)
			value := sv.MapIndex(reflect.ValueOf(f.Name))
			//log.Printf("value: %+v f.name: %+v type: %+v datasturct: %+v field: %s", value, f.Name, value.Kind(), dataStruct, fields.Fmt())
			targetValue := value.Interface()
			/*	if value.IsNil() {
				data := ""
				targetValue = reflect.ValueOf(data)
			}*/
			valueList = append(valueList, targetValue)
		}

	case reflect.Struct:
		numFields := st.NumField()
		for i := 0; i < numFields; i++ {
			field := st.Field(i)

			fieldName := field.Name
			tag, ok := field.Tag.Lookup("json")
			if ok {
				ss := strings.SplitN(tag, ",", 2)
				if len(ss) > 0 {
					fieldName = ss[0]
				}
				if fieldName == "-" {
					continue
				}
			}
			if !ok {
				err = fmt.Errorf("no tag. field: %s", fieldName)
				return
			}
			fieldNameList = append(fieldNameList, fieldName)
			reflectValue := sv.Field(i)

			kind := field.Type.Kind()
			var data interface{}
			switch kind {
			case reflect.Struct, reflect.Slice, reflect.Map:
				//newL = SetLoggerWithStructFields(ctx, l, reflectValue.Interface())
				var dataBytes []byte
				dataBytes, err = json.Marshal(reflectValue.Interface())
				if err != nil {
					err = fmt.Errorf("struct slice not support. err: %s", err)
					return
				}
				data = string(dataBytes)

			case reflect.String:
				data = reflectValue.String()
			case reflect.Float32, reflect.Float64:
				data = reflectValue.Float()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				data = reflectValue.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				data = reflectValue.Uint()
			default:
				data = reflectValue.String()
			}

			valueList = append(valueList, data)

		}

		return
	default:
		err = fmt.Errorf("dataStruct must be a struct or map")
		return
	}

	return

}

const (
	WRITE_NUM = 2000
)

func (t *TableConn) WriteRow(dataStruct interface{}, force bool) (result bool, err error) {
	fields, err := t.GetFields()
	if err != nil {
		return
	}
	valueList, fieldNameList, err := InterfaceToSqlParam(dataStruct, fields)
	if err != nil {
		return
	}

	for i, f := range fields {
		if f.Name != fieldNameList[i] {
			err = fmt.Errorf("fields are not equal. fields: %+v fieldsnamelist: %v i: %d f.Name: %s fieldNameList[i]: %s", fields, fieldNameList, i, f.Name, fieldNameList[i])
			return
		}
	}

	//beforeLength := len(t.cacheWriteList)
	t.writerLock.Lock()
	t.cacheWriteList = append(t.cacheWriteList, valueList)
	t.writerLock.Unlock()

	writeNum := 3000
	if t.WriteNumPerTime != 0 {
		writeNum = t.WriteNumPerTime
	}

	//log.Printf("writeNum: %d len(t.cacheWriteList): %d valueList: %+v beforeLength: %d", writeNum, len(t.cacheWriteList), valueList, beforeLength)

	if force || len(t.cacheWriteList) >= writeNum {
		result = true
		err = t.Flush()
		if err != nil {
			return
		}
	}
	return
}
func (t *TableConn) GetDataNum() (length int) {
	t.writerLock.Lock()
	defer t.writerLock.Unlock()
	return len(t.cacheWriteList)
}
func (t *TableConn) Flush() (err error) {
	t.writerLock.Lock()
	defer t.writerLock.Unlock()
	if len(t.cacheWriteList) > 0 {
		sql := t.writeSql
		start := time.Now()
		if sql == "" {
			sql = "INSERT "
			if t.ignore {
				sql += " ignore "
			}
			sql += fmt.Sprintf(" into %s (", t.TableName)
			placeHolder := "("
			duplicateString := ""
			for i, f := range t.Fields {
				if i != 0 {
					sql = sql + ","
					placeHolder = placeHolder + ","
					duplicateString = duplicateString + ","
				}
				sql += fmt.Sprintf("`%s`", f.Name)
				placeHolder += "?"
				duplicateString += fmt.Sprintf("`%s`=VALUES(`%s`)", f.Name, f.Name)
			}
			placeHolder = placeHolder + ")"
			sql = sql + ") values " + placeHolder
			sql = sql + " ON DUPLICATE KEY UPDATE "

			if len(t.WriteOnDuplicateFieldList) > 0 {
				duplicateString = ""
				for i, f := range t.WriteOnDuplicateFieldList {
					if i != 0 {
						duplicateString = duplicateString + ","
					}
					duplicateString += fmt.Sprintf("`%s`=VALUES(`%s`)", f, f)
				}
			}
			sql = sql + duplicateString

		}

		tryTimes := 1000
		var rowsAffected int64
		for {

			rowsAffected, err = func() (tmpRowsAffected int64, innerErr error) {
				tx, innerErr := t.DB.Begin()
				if innerErr != nil {
					return
				}
				defer func() {
					if innerErr != nil {
						rollBackErr := tx.Rollback()
						if rollBackErr != nil {
							innerErr = fmt.Errorf("insert error: %s and rollback err: %s", innerErr.Error(), rollBackErr.Error())
							return
						}
					}
				}()

				stmt, innerErr := tx.Prepare(sql)
				if innerErr != nil {
					return
				}

				tmpRowsAffected = 0
				for _, valueList := range t.cacheWriteList {
					var insertResult sqllib.Result
					insertResult, innerErr = stmt.Exec(valueList...)
					if innerErr != nil {
						return
					}
					var ef int64
					ef, innerErr = insertResult.RowsAffected()
					if err != nil {
						return
					}
					tmpRowsAffected = tmpRowsAffected + ef
				}
				innerErr = tx.Commit()
				if innerErr != nil {
					return
				}
				return
			}()
			if err == nil {
				break
			}

			log.Printf("insert tableName: %s error: %s. tryTimes: %d", t.TableName, err.Error(), tryTimes)
			lowerErrMsg := strings.ToLower(err.Error())
			if strings.Contains(lowerErrMsg, "try restarting transaction") ||
				strings.Contains(lowerErrMsg, "lost connection") ||
				strings.Contains(lowerErrMsg, "has gone away") ||
				strings.Contains(lowerErrMsg, "invalid connection") {
				log.Printf("writeRows retry Exception: %s,tryTimes: %d", lowerErrMsg, tryTimes)
				time.Sleep(time.Second * time.Duration(rand.Intn(10)+2))
				if tryTimes > 0 {
					tryTimes = tryTimes - 1
				} else {
					log.Printf("retry %d times,but did not pass,err: %s", tryTimes, lowerErrMsg)
					break
				}

			} else if strings.Contains(lowerErrMsg, "is full") {
				break
			} else {
				log.Printf("writeRows fatal Exception: %s,try_times: %d,ignore: %v sql: %s demodata: %s(%d)",
					lowerErrMsg, tryTimes, t.ignore, sql, t.cacheWriteList[0], len(t.cacheWriteList[0]))
				break
			}

		}
		cacheDataLength := int64(len(t.cacheWriteList))
		cost := time.Since(start)
		t.TotalWriteNum += cacheDataLength
		t.TotalAffectedNum += rowsAffected
		t.TotalCost = t.TotalCost + cost

		log.Printf("Finishing inserting Table: %s data num: %d. affected: %d timecost: %s. Total: WriteNum: %d Affected: %d Cost: %s",
			t.TableName, cacheDataLength, rowsAffected, cost.String(), t.TotalWriteNum, t.TotalAffectedNum, t.TotalCost.String())
		t.cacheWriteList = t.cacheWriteList[0:0]
	}
	return
}
