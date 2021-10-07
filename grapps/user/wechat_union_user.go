package user

import (
	//drsql "database/sql"
	"encoding/json"
	"fmt"

	"github.com/gdgrc/grutils/grapps/config"
	"github.com/gdgrc/grutils/grcache"
	//"grutils/grdatabase"
)

const (
	REDIS_PREFIX_WECHAT_UNION_USER_FOR_PK = "RedisWechatUnionUser?appid=%s&wxOpAccount=%s&wxUnionId=%s"
)

type WechatUnionUser struct {
	Appid           string                   `json:"-"`
	Uid             int64                    `json:"uid"`
	WxOpAccount     string                   `json:"wx_op_account"`
	WxUnionId       string                   `json:"wx_unionid"`
	CreateTime      string                   `json:"create_time"`
	ExtraInfo       WechatUnionUserExtraInfo `json:"extra_info"`
	ExtraInfoString string                   `json:"-"`
	RequestId       string                   `json:"-"`
}

type WechatUnionUserExtraInfo struct {
	UpdateTime *string `json:"update_time,omitempty"`
}

func (uo *WechatUnionUser) InsertDB() (err error) {

	newItemSlce := make([]*WechatUnionUser, 0, 1)
	newItemSlce = append(newItemSlce, uo)

	config.DefaultLogger.Debug("newItemSlce: %+v, WechatUnionUser: %+v", newItemSlce, *uo)

	_, _, err = WechatUnionUserInsert(newItemSlce)
	return

}
func (uo *WechatUnionUser) SaveDB() (err error) {

	newItemSlce := make([]*WechatUnionUser, 0, 1)
	newItemSlce = append(newItemSlce, uo)

	config.DefaultLogger.Debug("newItemSlce: %+v, WechatUnionUser: %+v", newItemSlce, *uo)

	_, _, err = WechatUnionUserInsertOnDuplicateUpdate(newItemSlce)
	return

}
func GetWechatUnionUserFromDBorCache(appid string, wxOpAccount string, wxUnionId string, slave bool) (newItem WechatUnionUser, exist bool, err error) {
	err, conn := grcache.RedisConnOPGet(config.GRAPPS_REDIS, &config.GlobalConf.RedisPool)

	if err != nil {

		return

	}
	defer conn.Close()

	redisKey := fmt.Sprintf(REDIS_PREFIX_WECHAT_USER_FOR_PK, appid, wxOpAccount, wxUnionId)
	if slave {

		exist, err = conn.Get(redisKey, &newItem)
		if err != nil {
			return
		}
		if exist {
			return
		}
	}

	sql := "SELECT `appid`,`uid`,`wx_op_account`,`wx_unionid`,date_format(`create_time`, '%Y-%m-%d %H:%i:%s'),`extra_info` FROM `wechat_union_user` where `appid`= ? AND `wx_op_account`= ? AND `wx_unionid`= ? limit 1"

	paramSlice := make([]interface{}, 0)
	paramSlice = append(paramSlice, appid)
	paramSlice = append(paramSlice, wxOpAccount)
	paramSlice = append(paramSlice, wxUnionId)

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}

	rows, err := adminConn.Query(sql, paramSlice...)
	if err != nil {

		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&newItem.Appid, &newItem.Uid, &newItem.WxOpAccount, &newItem.WxUnionId, &newItem.CreateTime, &newItem.ExtraInfoString)
		if err != nil {
			return
		}
		if err = json.Unmarshal([]byte(newItem.ExtraInfoString), &newItem.ExtraInfo); err != nil {
			return
		}
		exist = true

	}

	if exist {
		err = conn.SetEx(redisKey, newItem, 600)
		if err != nil {
			return
		}
	}

	/*
		err = adminConn.QueryRow(sql, paramSlice...).Scan(&newItem.Uid, &newItem.WxOpAccount, &newItem.WxUnionId, &newItem.ExtraInfoString)
		if err != nil {
			err = fmt.Errorf("client query no data err:%s,sql: %s, uid: %d, slave: %d", err, sql, uid, slave)
			return
		}
		if err = json.Unmarshal([]byte(newItem.ExtraInfoString), &newItem.ExtraInfo); err != nil {
			return
		}*/

	return

}
func WechatUnionUserInsert(newItemSlice []*WechatUnionUser) (successSlice []string, successNum int, err error) {
	//INSERT INTO `client` (`appid`,`client_id`,`pt_id`,`customer_id`,`create_time`, `status`,`extra_info`) VALUES('test1','test2','test3','test4','2018-05-05 00:00:00',1,'{}') ON DUPLICATE KEY UPDATE `extra_info` =JSON_SET(`extra_info`,'$.head_url','good2','$.sex', 2);
	sql := "INSERT  INTO `wechat_union_user` (`appid`,`uid`, `wx_op_account`,`wx_unionid`,`create_time`,`extra_info`) "
	sql += "VALUES(?, ?, ?, ? ,? ,?)"
	//sql += "`extra_info`= JSON_SET(`extra_info`,'$.update_time', ? , $.head_url', ? ,'$.sex', ? )"

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}

	successSlice = make([]string, 0)
	successNum = 0

	for _, newItem := range newItemSlice {
		//newItem := *pNewItem

		var extra_info_bytes []byte
		extra_info_bytes, err = json.Marshal(newItem.ExtraInfo)
		if err != nil {

			return
		}
		newItem.ExtraInfoString = string(extra_info_bytes)

		paramSlice := make([]interface{}, 0)
		paramSlice = append(paramSlice, newItem.Appid, newItem.Uid, newItem.WxOpAccount, newItem.WxUnionId, newItem.CreateTime, newItem.ExtraInfoString)

		// newItem.Appid, newItem.ClientId, newItem.PtId, newItem.CustomerId, newItem.CreateTime, newItem.Status, extra_info_string,
		//newItem.CustomerId,
		//newItem.ExtraInfo.UpdateTime, newItem.ExtraInfo.HeadUrl, newItem.ExtraInfo.Sex

		_, err = adminConn.Exec(sql, paramSlice...)
		config.DefaultLogger.Debug("t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
		if err != nil {
			return
		}
		successSlice = append(successSlice, newItem.RequestId)
		successNum += 1
	}

	return

}

func WechatUnionUserInsertOnDuplicateUpdate(newItemSlice []*WechatUnionUser) (successSlice []string, successNum int, err error) {
	//INSERT INTO `client` (`appid`,`client_id`,`pt_id`,`customer_id`,`create_time`, `status`,`extra_info`) VALUES('test1','test2','test3','test4','2018-05-05 00:00:00',1,'{}') ON DUPLICATE KEY UPDATE `extra_info` =JSON_SET(`extra_info`,'$.head_url','good2','$.sex', 2);
	sql := "INSERT  INTO `wechat_union_user` (`appid`,`uid`,`wx_op_account`,`wx_unionid`,`create_time`,`extra_info`) "
	sql += "VALUES(? ,? ,? , ?, ? ,?) ON DUPLICATE KEY UPDATE `extra_info`=VALUES(`extra_info`)"
	//sql += "`extra_info`= JSON_SET(`extra_info`,'$.update_time', ? , $.head_url', ? ,'$.sex', ? )"

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}

	successSlice = make([]string, 0)
	successNum = 0

	for _, newItem := range newItemSlice {
		//newItem := *pNewItem

		var extra_info_bytes []byte
		extra_info_bytes, err = json.Marshal(newItem.ExtraInfo)
		if err != nil {

			return
		}
		newItem.ExtraInfoString = string(extra_info_bytes)

		paramSlice := make([]interface{}, 0)
		paramSlice = append(paramSlice, newItem.Appid, newItem.Uid, newItem.WxOpAccount, newItem.WxUnionId, newItem.CreateTime, newItem.ExtraInfoString)

		// newItem.Appid, newItem.ClientId, newItem.PtId, newItem.CustomerId, newItem.CreateTime, newItem.Status, extra_info_string,
		//newItem.CustomerId,
		//newItem.ExtraInfo.UpdateTime, newItem.ExtraInfo.HeadUrl, newItem.ExtraInfo.Sex

		_, err = adminConn.Exec(sql, paramSlice...)
		config.DefaultLogger.Debug("t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
		if err != nil {
			return
		}
		successSlice = append(successSlice, newItem.RequestId)
		successNum += 1
	}

	return

}
