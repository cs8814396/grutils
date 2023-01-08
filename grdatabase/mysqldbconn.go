package grdatabase

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/gdgrc/grutils/grcommon"

	_ "github.com/go-sql-driver/mysql"
)

func NewDatabaseConn(s *sql.DB) *DatabaseConn {
	databaseConn := &DatabaseConn{s}
	return databaseConn
}

type DatabaseConn struct {
	*sql.DB
}

func (d *DatabaseConn) GetTableList() (tableList []string, err error) {

	rawFuncDescSql := "SHOW TABLES"
	rows, err := d.Query(rawFuncDescSql)
	if err != nil {
		return
	}
	tableList = make([]string, 0)
	for rows.Next() {
		var name string

		err = rows.Scan(&name)
		if err != nil {

			return
		}
		tableList = append(tableList, name)

	}
	return

}

func (d *DatabaseConn) NewTableConn(tableName string) *TableConn {
	var tc TableConn
	tc.DB = d.DB
	tc.TableName = tableName

	return &tc
}

type TableConn struct {
	DatabaseConn
	TableName         string
	Fields            []*Field
	IsTurnColumnUpper bool

	LastQueryReq *QueryReq
	// writer
	writerLock                sync.Mutex
	cacheWriteList            [][]interface{}
	writeSql                  string
	ignore                    bool
	WriteOnDuplicateFieldList []string
	TotalWriteNum             int64
	TotalAffectedNum          int64
	WriteNumPerTime           int
}

type QueryReq struct {
	ReadLineNum      int
	BeginIndex       int
	BeginIdIndex     interface{}
	BeginIdIndexName string
	ExtraSql         string

	SelectFields []string
}

// 可能会重复读，但是不会少读。
func (t *TableConn) ReadNextDataToMap(q *QueryReq) (rowsMap []map[string]string, err error) {
	err = t.DB.Ping()
	if err != nil {
		return
	}
	if q != nil {
		t.LastQueryReq = q
	}
	if t.LastQueryReq == nil {
		t.LastQueryReq = &QueryReq{ReadLineNum: 100000}
	}
	sql, args, err := t.GetSelectSql(t.LastQueryReq)
	if err != nil {
		return
	}

	rowsMap, err = t.QueryToMap(sql, args, false)
	if err != nil {
		return
	}

	//log.Printf("ReadNext.Sql: %s,args: %+v length: %d\n", sql, args, len(rowsMap))
	if len(rowsMap) == 0 {
		// read ot
		log.Printf("Table: %s Read Finish. sql: %s args: %+v \n", t.TableName, sql, args)
		return
	}
	if t.LastQueryReq.BeginIdIndexName != "" {
		newBeginIdIndex := rowsMap[len(rowsMap)-1][t.LastQueryReq.BeginIdIndexName]
		if t.LastQueryReq.BeginIdIndex == nil {
			// 没有开始id, 理论上同最后一个情况。至少跳过一行
			t.LastQueryReq.BeginIndex = 1
		} else if newBeginIdIndex == t.LastQueryReq.BeginIdIndex {
			// 假如最后一行的id值 和 查询条件开始的id值相同，证明这次的查询结果所有的 id值都是相同的，所以需要跳过这个查询部分，否则就会死循环
			t.LastQueryReq.BeginIndex = t.LastQueryReq.BeginIndex + t.LastQueryReq.ReadLineNum
		} else {
			// 除去上面两种情况，读到了不同于查询起始条件的值，又因为读到过，所以下次开始的时候至少跳过一行。 (的确可能重复读)
			// newBeginIdIndex!= t.LastQueryReq.BeginIdIndex
			t.LastQueryReq.BeginIndex = 1
		}
		t.LastQueryReq.BeginIdIndex = newBeginIdIndex

	} else {
		t.LastQueryReq.BeginIndex = t.LastQueryReq.BeginIndex + t.LastQueryReq.ReadLineNum
	}

	return
}

func (t *TableConn) QueryToMap(querySql string, params []interface{}, isPure bool) (retSlice []map[string]string, err error) {
	tablRows, err := t.Query(querySql, params...)
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

			if t.IsTurnColumnUpper {
				columnName = strings.ToUpper(columnName)
			}
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

func (t *TableConn) GetSelectSql(q *QueryReq) (sql string, args []interface{}, err error) {
	sql = "SELECT "

	if q != nil && len(q.SelectFields) > 0 {

		for _, field := range q.SelectFields {
			sql += fmt.Sprintf("`%s`,", field)
		}
		lengthSql := len(sql)
		sql = sql[:lengthSql-1]
	} else {
		var fields Fields
		fields, err = t.GetFields()
		if err != nil {
			return
		}

		for _, field := range fields {
			sql += fmt.Sprintf("`%s`,", field.Name)
		}
		lengthSql := len(sql)
		sql = sql[:lengthSql-1]
	}
	order := ""

	sql += fmt.Sprintf(" FROM `%s` ", t.TableName)

	args = make([]interface{}, 0, 0)

	if q != nil {

		if q.BeginIdIndexName != "" {

			if q.BeginIdIndex != nil {
				// 有可能读重复
				sql += fmt.Sprintf(" WHERE `%s` >= ? ", q.BeginIdIndexName)
				args = append(args, q.BeginIdIndex)

			}
			order = " order by " + q.BeginIdIndexName
		}

		if q.ExtraSql != "" {
			if strings.Contains(sql, "WHERE") {
				sql += (" AND " + q.ExtraSql)
			} else {
				sql += ("WHERE " + q.ExtraSql)
			}
		}
		sql += order
		// limit begin
		if q.ReadLineNum > 0 {
			sql += " LIMIT "
			if q.BeginIndex != 0 {
				//begin = begin_index if begin_index > -1 else 0
				sql += fmt.Sprintf(" %d, ", (q.BeginIndex))
			}

			sql += fmt.Sprintf(" %d", q.ReadLineNum)
		}
	}

	return

}

/*

try:




	# print(sql, args_list)
	rows = None
	rows_length = None
	if dict_query:
		rows, rows_length = self.db_conn_object.dict_query(sql, args_list)
	else:
		rows, rows_length = self.db_conn_object.query(sql, args_list)

	"""
	# i do not think we should travel
	for row in rows:
		tmp_list = []
		for data in row:
			data_tmp = ''
			if data is not None:
				data_tmp = str(data)

			tmp_list.append(data)
		data_list.append(tmp_list)
	"""

	# msg = "Finish reading table %s,read_index: %d,total_len: %d,one for fields' name" % (tbname, demoindex, len(data_list))
	# echomsg(msg, False)

	return rows, rows_length
except Exception as e:
	raise Exception("sql execute error: " + sql)
*/
type Fields []*Field

func (fs Fields) Fmt() string {
	s := ""
	for i, f := range fs {
		s += fmt.Sprintf("i: %d f: %+v\n", i, f)
	}
	return s
}
func (fs Fields) FindColumnName(name string) (ok bool) {
	ok = false
	for _, v := range fs {
		if v.Name == name {
			ok = true

		}
	}
	return
}

type Field struct {
	Name   string
	Type   string
	IsNull bool
}

func (t *TableConn) GetFields() (fields Fields, err error) {
	if t.Fields == nil {

		rawFuncDescSql := fmt.Sprintf("DESC `%s`", t.TableName)
		rows, err0 := t.Query(rawFuncDescSql)
		if err0 != nil {
			err = err0
			return

		}
		fields = make([]*Field, 0)
		for rows.Next() {
			var IsNull string
			var t2, t3, t4 interface{}

			var field = Field{}
			err = rows.Scan(&field.Name, &field.Type, &IsNull, &t2, &t3, &t4)
			if err != nil {

				return
			}
			if IsNull == "YES" {
				field.IsNull = true
			}
			if t.IsTurnColumnUpper {
				field.Name = strings.ToUpper(field.Name)
			}
			field.Type = strings.ToUpper(field.Type)
			fields = append(fields, &field)

		}
		t.Fields = fields
	}
	return t.Fields, nil
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
		db.Close()
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
	var db0 *sql.DB
	retryTimes := 2
	for i := 0; i <= retryTimes; i++ {
		i++
		db0, err = OpenAndTestDB(dbname, dsn, maxidleconn)
		if err == nil {
			break
		}

	}
	if err != nil {
		err = errors.New("OpenAndTestDB dbname=" + dbname + ", dsn=" + dsn + " idle=" + strconv.Itoa(maxidleconn) + " err=" + err.Error() + "retryTimes=" + strconv.Itoa(retryTimes))
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
