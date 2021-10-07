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
	REDIS_PREFIX_WECHAT_USER_FOR_PK = "RedisWechatUser?appid=%s&wxAppid=%s&wxOpenid=%s"
)

type WechatUser struct {
	Appid           string              `json:"-"`
	Uid             int64               `json:"uid"`
	WxAppid         string              `json:"-"`
	WxOpenid        string              `json:"-"`
	CreateTime      string              `json:"create_time"`
	ExtraInfo       WechatUserExtraInfo `json:"extra_info"`
	ExtraInfoString string              `json:"-"`
	RequestId       string              `json:"-"`
}

type WechatUserExtraInfo struct {
	UpdateTime *string   `json:"update_time,omitempty"`
	Ip         *string   `json:"ip,omitempty"`
	HeadImgUrl *string   `json:"headimgurl,omitempty"`
	Privilege  *[]string `json:"privilege,omitempty"`
	Nickname   *string   `json:"nickname,omitempty"`
	Sex        *int      `json:"sex,omitempty"`
	Province   *string   `json:"province,omitempty"`
	City       *string   `json:"city,omitempty"`
	Country    *string   `json:"country,omitempty"`
	//Openid     *string   `json:"openid,omitempty"`
}

func (uo *WechatUser) InsertDB() (err error) {

	newItemSlce := make([]*WechatUser, 0, 1)
	newItemSlce = append(newItemSlce, uo)

	config.DefaultLogger.Debug("newItemSlce: %+v, WechatUser: %+v", newItemSlce, *uo)

	_, _, err = WechatUserInsert(newItemSlce)
	return

}
func (uo *WechatUser) SaveDB() (err error) {

	newItemSlce := make([]*WechatUser, 0, 1)
	newItemSlce = append(newItemSlce, uo)

	config.DefaultLogger.Debug("newItemSlce: %+v, WechatUser: %+v", newItemSlce, *uo)

	_, _, err = WechatUserInsertOnDuplicateUpdate(newItemSlce)
	return

}
func GetWechatUserFromDBorCache(appid string, wxAppid string, wxOpenid string, slave bool) (newItem WechatUser, exist bool, err error) {
	err, conn := grcache.RedisConnOPGet(config.GRAPPS_REDIS, &config.GlobalConf.RedisPool)

	if err != nil {

		return

	}
	defer conn.Close()

	redisKey := fmt.Sprintf(REDIS_PREFIX_WECHAT_USER_FOR_PK, appid, wxAppid, wxOpenid)
	if slave {

		exist, err = conn.Get(redisKey, &newItem)
		if err != nil {
			return
		}
		if exist {
			return
		}
	}

	sql := "SELECT `appid`,`uid`,`wx_appid`,`wx_openid`,date_format(`create_time`, '%Y-%m-%d %H:%i:%s'),`extra_info` FROM `wechat_user` where `appid`= ? AND `wx_appid`= ? AND `wx_openid`= ? limit 1"

	paramSlice := make([]interface{}, 0)
	paramSlice = append(paramSlice, appid)
	paramSlice = append(paramSlice, wxAppid)
	paramSlice = append(paramSlice, wxOpenid)

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
		err = rows.Scan(&newItem.Appid, &newItem.Uid, &newItem.WxAppid, &newItem.WxOpenid, &newItem.CreateTime, &newItem.ExtraInfoString)
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
		err = adminConn.QueryRow(sql, paramSlice...).Scan(&newItem.Uid, &newItem.WxAppid, &newItem.WxOpenid, &newItem.ExtraInfoString)
		if err != nil {
			err = fmt.Errorf("client query no data err:%s,sql: %s, uid: %d, slave: %d", err, sql, uid, slave)
			return
		}
		if err = json.Unmarshal([]byte(newItem.ExtraInfoString), &newItem.ExtraInfo); err != nil {
			return
		}*/

	return

}
func WechatUserInsert(newItemSlice []*WechatUser) (successSlice []string, successNum int, err error) {
	//INSERT INTO `client` (`appid`,`client_id`,`pt_id`,`customer_id`,`create_time`, `status`,`extra_info`) VALUES('test1','test2','test3','test4','2018-05-05 00:00:00',1,'{}') ON DUPLICATE KEY UPDATE `extra_info` =JSON_SET(`extra_info`,'$.head_url','good2','$.sex', 2);
	sql := "INSERT  INTO `wechat_user` (`appid`,`uid`, `wx_appid`,`wx_openid`,`create_time`,`extra_info`) "
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
		paramSlice = append(paramSlice, newItem.Appid, newItem.Uid, newItem.WxAppid, newItem.WxOpenid, newItem.CreateTime, newItem.ExtraInfoString)

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

func WechatUserInsertOnDuplicateUpdate(newItemSlice []*WechatUser) (successSlice []string, successNum int, err error) {
	//INSERT INTO `client` (`appid`,`client_id`,`pt_id`,`customer_id`,`create_time`, `status`,`extra_info`) VALUES('test1','test2','test3','test4','2018-05-05 00:00:00',1,'{}') ON DUPLICATE KEY UPDATE `extra_info` =JSON_SET(`extra_info`,'$.head_url','good2','$.sex', 2);
	sql := "INSERT  INTO `wechat_user` (`appid`,`uid`,`wx_appid`,`wx_openid`,`create_time`,`extra_info`) "
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
		paramSlice = append(paramSlice, newItem.Appid, newItem.Uid, newItem.WxAppid, newItem.WxOpenid, newItem.CreateTime, newItem.ExtraInfoString)

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
