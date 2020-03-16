package grdatabase

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gdgrc/grutils/grcommon"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"sync"
)

type DatabaseConn struct {
	*sql.DB
}

func (this *DatabaseConn) NewTableConn(tableName string) *TableConn {
	var tc TableConn
	tc.DB = this.DB
	tc.TableName = tableName

	return &tc
}

type TableConn struct {
	DatabaseConn
	TableName string
	Fields    []string
}

func (this *TableConn) QueryToMap(querySql string, params []interface{}, isPure bool) (retSlice []map[string]string, err error) {
	tablRows, err := this.Query(querySql, params...)
	if err != nil {
		return
	}
	defer tablRows.Close()

	columns, err := tablRows.Columns()
	if err != nil {

		return
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	retSlice = make([]map[string]string, 0)

	pureMd5Map := make(map[string]int)
	for tablRows.Next() {
		err = tablRows.Scan(scanArgs...)
		if err != nil {

			return
		}

		tmpMap := make(map[string]string) // origin_data

		for columnIndex, columnName := range columns {

			tmpMap[columnName] = string(values[columnIndex])
		}

		if isPure {
			var mapMd5 string
			mapMd5, err = grcommon.GetStringMapMd5(tmpMap)
			if err != nil {
				return
			}
			if _, ok := pureMd5Map[mapMd5]; ok {
				continue
			}

		}
		retSlice = append(retSlice, tmpMap)
	}

	return
}

func (this *TableConn) GetSelectSql() (sql string, err error) {
	sql = "SELECT "
	fields, err := this.GetFields()
	if err != nil {
		return
	}

	for _, field := range fields {

		sql += fmt.Sprintf("`%s`,", field)
	}

	lengthSql := len(sql)
	sql = sql[:lengthSql-1]

	sql += fmt.Sprintf(" FROM `%s` ", this.TableName)

	return

}

func (this *TableConn) GetFields() (fields []string, err error) {
	if this.Fields == nil {

		rawFuncDescSql := fmt.Sprintf("DESC `%s`", this.TableName)
		rows, err0 := this.Query(rawFuncDescSql)
		if err0 != nil {
			err = err0
			return

		}
		fields = make([]string, 0)
		for rows.Next() {
			var field, ttype string
			var t1, t2, t3, t4 interface{}
			err = rows.Scan(&field, &ttype, &t1, &t2, &t3, &t4)
			if err != nil {

				return
			}
			fields = append(fields, field)

		}
		this.Fields = fields
	}
	return this.Fields, nil
}

type XmlMysqlDsn struct {
	AdminDsn string `xml:"admin_dsn"`

	MaxIdleConn int `xml:"maxidleconn"`
}

type MysqlPoolMap struct {
	//mutex   *sync.RWMutex
	pool        sync.Map //map[string]*sql.DB
	dnsPool     sync.Map
	maxidlePool sync.Map
}

const (
	MINIMAL_DB_CONN = 50
)

var DefaultMysqlPool MysqlPoolMap

func OpenAndTestDB(name, dsn string, maxidle int) (db *sql.DB, err error) {

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.New("DB " + name + " init err=" + err.Error())
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.New("DB " + name + " ping err=" + err.Error())
	}

	db.SetMaxIdleConns(maxidle)
	return
}

func (mp *MysqlPoolMap) DBGetConn(dbname string, dsn string, maxidleconn int) (db *sql.DB, err error) {

	err = nil

	if idb, ok := mp.pool.Load(dbname); ok {
		db = idb.(*sql.DB)

		return
	}

	if dsn == "" {
		err = errors.New("DBGetConn not exist but no dsn: dbname: " + dbname)
		return
	}

	if maxidleconn < MINIMAL_DB_CONN {

		maxidleconn = MINIMAL_DB_CONN
	}

	db0, err0 := OpenAndTestDB(dbname, dsn, maxidleconn)
	if err0 != nil {

		err = errors.New("OpenAndTestDB dbname=" + dbname + ", dsn=" + dsn + " idle=" + strconv.Itoa(maxidleconn) + " err=" + err0.Error())
		return
	}

	db = db0

	mp.pool.Store(dbname, db)
	mp.dnsPool.Store(dbname, dsn)
	mp.maxidlePool.Store(dbname, maxidleconn)

	return
}
func MysqlClearTransaction(tx *sql.Tx, rollFlag *bool) {
	// fmt.Println("rollback rollFlag: ", *rollFlag)
	if *rollFlag {

		err := tx.Rollback()
		if err != sql.ErrTxDone && err != nil {

			panic("rollback err!")
		}
	} else {
		err := tx.Commit()
		if err != nil {

			panic("commit err!")
		}
	}
}
