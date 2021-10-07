package trade

import (
	"database/sql"
	"encoding/json"

	//"grutils/grmath"
	"fmt"
	"time"

	"github.com/gdgrc/grutils/grapps/config"
)

type TradeLog struct {
	Appid string `json:"appid"`
	Tid   string `json:"tid"`
	Uid   int64  `json:"uid"`

	Type int `json:"type"`

	Way        int    `json:"way"`
	CreateTime string `json:"create_time"`

	ExtraInfoString string            `json:"-"`
	ExtraInfo       TradeLogExtraInfo `json:"extra_info"`
}
type TradeLogExtraInfo struct {
	BaseMoney   *float64 `json:"base_money,omitempty"`
	GiftMoney   *float64 `json:"gift_money,omitempty"`
	RelUid      *int64   `json:"rel_uid,omitempty"`
	RelNickname *string  `json:"rel_nickname,omitempty"`
	OrderCode   *string  `json:"order_code,omitempty"`
	Desc        *string  `json:"desc,omitempty"`
	IsHidden    *int     `json:"is_hidden,omitempty"`
}

const (
	TRADE_LOG_TYPE_RECHARGE       = 1
	TRADE_LOG_TYPE_REWARD_PRODUCT = 2

	TRADE_LOG_TYPE_CHAIN_REWARD = 3

	TRADE_LOG_TYPE_ORDER_REFUND = 4
)
const (
	TRADE_LOG_WAY_IN  = 1
	TRADE_LOG_WAY_OUT = 2
)

func (tl *TradeLog) TradeLogTxInsert(tx *sql.Tx) (err error) {
	sql := "INSERT INTO `trade_log` (`appid`,`tid`,`uid`,`type` , `way` ,`create_time`, `extra_info`) "
	sql += " VALUES (? ,? ,? ,? ,? ,? ,? )"
	//sql += "`extra_info`= JSON_SET(`extra_info`,'$.update_time', ? , $.head_url', ? ,'$.sex', ? )"
	var extra_info_bytes []byte
	extra_info_bytes, err = json.Marshal(tl.ExtraInfo)
	if err != nil {

		return
	}
	//tl.CreateTime = grmath.GetCurTime()
	tl.ExtraInfoString = string(extra_info_bytes)

	paramSlice := make([]interface{}, 0)
	paramSlice = append(paramSlice, tl.Appid, tl.Tid, tl.Uid, tl.Type, tl.Way, tl.CreateTime, tl.ExtraInfoString)

	_, err = tx.Exec(sql, paramSlice...)
	config.DefaultLogger.Debug("UserChargeMoney t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
	if err != nil {
		return
	}
	return
}

type TradeLogStat struct {
	SumBaseMoney float64 `json:"sum_base_money"`
	SumGiftMoney float64 `json:"sum_gift_money"`
	Count        int     `json:"count"`
}

func GetTradeLogStatFromDBorCache(appid string, uid int64, iType int, beginTime string, endTime string, slave bool) (tls TradeLogStat, err error) {
	_, err = time.Parse("2006-01-02 15:04:05", beginTime)
	if err != nil {
		return
	}

	if slave {
		//get from cache and return if exists
	}

	paramSlice := make([]interface{}, 0)

	sql := "SELECT count(*),FORMAT(IFNULL(sum(`extra_info`->>'$.base_money'),0.0),2) as sum_base_money,FORMAT(IFNULL(sum(`extra_info`->>'$.gift_money'),0.0),2) as sum_gift_money "
	sql += " FROM `trade_log` where `appid`= ? AND `uid`= ? AND (`extra_info`->>'$.is_hidden' is NULL or `extra_info`->>'$.is_hidden' =0 )"
	paramSlice = append(paramSlice, appid)
	paramSlice = append(paramSlice, uid)

	sql += " AND `create_time` >= ? "
	paramSlice = append(paramSlice, beginTime)

	if endTime != "" {
		_, err = time.Parse("2006-01-02 15:04:05", endTime)
		if err != nil {
			return
		}
		sql += " AND `create_time`<= ? "
		paramSlice = append(paramSlice, endTime)
	}

	if iType > 0 {
		sql += " AND `type` = ? "

		paramSlice = append(paramSlice, iType)
	}

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}

	config.DefaultLogger.Debug("sql: %s,param: %+v", sql, paramSlice)

	rows, err := adminConn.Query(sql, paramSlice...)
	if err != nil {
		config.DefaultLogger.Error("admin query  fail, err: %s", err)
		return
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&tls.Count, &tls.SumBaseMoney, &tls.SumGiftMoney)
		if err != nil {
			return
		}

		break
		//psSlice = append(psSlice, ps)

	}

	config.DefaultLogger.Debug("TradeLogStat: %+v", tls)

	return

}

func GetTradeLogSliceFromDBorCache(appid string, uid int64, iType int, beginTime string, endTime string, beginIndex int, count int, sortType int, slave bool) (dataSlice *[]TradeLog, err error) {
	_, err = time.Parse("2006-01-02 15:04:05", beginTime)
	if err != nil {
		return
	}

	if slave {
		//get from cache and return if exists
	}

	paramSlice := make([]interface{}, 0)

	sql := "SELECT `appid`,`tid`,`uid`,`type` , `way` ,date_format(`create_time`, '%Y-%m-%d %H:%i:%s'), `extra_info`"
	sql += " FROM `trade_log` where `appid`= ? AND `uid`= ? AND (`extra_info`->>'$.is_hidden' is NULL or `extra_info`->>'$.is_hidden' =0 )"
	paramSlice = append(paramSlice, appid)
	paramSlice = append(paramSlice, uid)

	sql += " AND `create_time` >= ? "
	paramSlice = append(paramSlice, beginTime)

	if endTime != "" {
		_, err = time.Parse("2006-01-02 15:04:05", endTime)
		if err != nil {
			return
		}
		sql += " AND `create_time`<= ? "
		paramSlice = append(paramSlice, endTime)
	}

	if iType > 0 {
		sql += " AND `type` = ? "

		paramSlice = append(paramSlice, iType)
	}
	if sortType == 0 {
		sql += " ORDER BY `create_time`"
	}

	sql += fmt.Sprintf(" limit %d, %d ", beginIndex, count)

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}

	config.DefaultLogger.Debug("sql: %s,param: %+v", sql, paramSlice)

	rows, err := adminConn.Query(sql, paramSlice...)
	if err != nil {
		config.DefaultLogger.Error("admin query  fail, err: %s", err)
		return
	}

	data0Slice := make([]TradeLog, 0)

	defer rows.Close()

	for rows.Next() {

		var tl TradeLog

		err = rows.Scan(&tl.Appid, &tl.Tid, &tl.Uid, &tl.Type, &tl.Way, &tl.CreateTime, &tl.ExtraInfoString)
		if err != nil {
			return
		}

		if err = json.Unmarshal([]byte(tl.ExtraInfoString), &tl.ExtraInfo); err != nil {
			return
		}

		data0Slice = append(data0Slice, tl)

		//psSlice = append(psSlice, ps)

	}

	config.DefaultLogger.Debug("data0Slice: %+v", data0Slice)

	dataSlice = &data0Slice

	return

}
