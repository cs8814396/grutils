package material

import (
	"database/sql"
	"encoding/json"

	//"grutils/grmath"

	"github.com/gdgrc/grutils/grapps/config"
)

type MaterialLog struct {
	Appid      string `json:"appid"`
	Mlid       int64  `json:"mlid"`
	Uid        int64  `json:"uid"`
	Type       int    `json:"type"`
	CreateTime string `json:"create_time"`

	ExtraInfoString string               `json:"-"`
	ExtraInfo       MaterialLogExtraInfo `json:"extra_info"`
}
type MaterialLogExtraInfo struct {
	AccessUrl *float64 `json:"access_url,omitempty"`
}

const (
	MATERIAL_LOG_TYPE_FISSION_VIDEO = 1
)

func (tl *MaterialLog) MaterialLogTxInsert(tx *sql.Tx) (err error) {
	sql := "INSERT INTO `material_log` (`appid`,`mlid`,`uid`,`type` ,`create_time`, `extra_info`) "
	sql += " VALUES (? ,? ,? ,? ,? ,?  )"
	//sql += "`extra_info`= JSON_SET(`extra_info`,'$.update_time', ? , $.head_url', ? ,'$.sex', ? )"
	var extra_info_bytes []byte
	extra_info_bytes, err = json.Marshal(tl.ExtraInfo)
	if err != nil {

		return
	}
	//tl.CreateTime = grmath.GetCurTime()
	tl.ExtraInfoString = string(extra_info_bytes)

	paramSlice := make([]interface{}, 0)
	paramSlice = append(paramSlice, tl.Appid, tl.Mlid, tl.Uid, tl.Type, tl.CreateTime, tl.ExtraInfoString)

	_, err = tx.Exec(sql, paramSlice...)
	config.DefaultLogger.Debug("MaterialLogTxInsert t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
	if err != nil {
		return
	}
	return
}
