package grdatabase

import (
	sqllib "database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"log"
)

func InterfaceToSqlParam(dataStruct interface{}) (valueList []interface{}, fieldNameList []string, err error) {
	sv := reflect.ValueOf(dataStruct)
	st := reflect.TypeOf(dataStruct)
	if sv.Type().Kind() == reflect.Ptr {
		sv = sv.Elem()
		st = st.Elem()
	}
	if sv.Type().Kind() != reflect.Struct {
		err = fmt.Errorf("dataStruct must be a struct")
		return
	}
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
		case reflect.Struct, reflect.Slice:
			//newL = SetLoggerWithStructFields(ctx, l, reflectValue.Interface())
			panic("not support")
			break

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

	//sql = fmt.Sprintf("INSERT INTO `%s` ",tableName)

	return
	//logger.Debugf("SetLoggerWithSt

}

const (
	WRITE_NUM = 2000
)

func (t *TableConn) WriteRow(dataStruct interface{}, force bool) (rowsMap []map[string]string, err error) {
	valueList, fieldNameList, err := InterfaceToSqlParam(dataStruct)
	if err != nil {
		return
	}
	fields, err := t.GetFields()
	if err != nil {
		return
	}
	for i, f := range fields {
		if f.Name != fieldNameList[i] {
			err = fmt.Errorf("fields are not equal. fields: %+v fieldsnamelist: %v i: %d", fields, fieldNameList, i)
			return
		}
	}
	t.writerLock.Lock()
	t.cacheWriteList = append(t.cacheWriteList, valueList)
	t.writerLock.Unlock()

	if force || len(t.cacheWriteList) > WRITE_NUM {
		err = t.Flush()
		if err != nil {
			return
		}
	}
	return
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

			rowsAffected, err = func() (rowsAffected int64, innerErr error) {
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

				rowsAffected = 0
				for _, valueList := range t.cacheWriteList {
					var insertResult sqllib.Result
					insertResult, innerErr = stmt.Exec(valueList...)
					if innerErr != nil {
						return
					}
					insertResult.RowsAffected()
					rowsAffected = rowsAffected + 1
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

			log.Printf("insert error: %s. tryTimes: %d", err.Error(), tryTimes)
			lowerErrMsg := strings.ToLower(err.Error())
			if strings.Contains(lowerErrMsg, "try restarting transaction") ||
				strings.Contains(lowerErrMsg, "lost connection") ||
				strings.Contains(lowerErrMsg, "has gone away") {
				log.Printf("writeRows retry Exception: %s,tryTimes: %d", lowerErrMsg, tryTimes)
				time.Sleep(time.Second * 2)
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

		log.Printf("Finishing inserting Table: %s data num: %d. timecost: %s", t.TableName, cacheDataLength, cost.String())
		t.cacheWriteList = t.cacheWriteList[0:0]
	}
	return
}
