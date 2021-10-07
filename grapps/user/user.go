package user

import (
	drsql "database/sql"
	"encoding/json"
	"fmt"

	"github.com/gdgrc/grutils/grapps/config"
	"github.com/gdgrc/grutils/grcache"
	"github.com/gdgrc/grutils/grdatabase"
)

const (
	REDIS_PREFIX_USER_FOR_PK = "RedisUser?appid=%s&uid=%d"
)

type User struct {
	Appid      string  `json:"appid"`
	Uid        int64   `json:"uid"`
	FrontUid   int64   `json:"front_uid"`
	BaseMoney  float64 `json:"base_money"`
	GiftMoney  float64 `json:"gift_money"`
	CreateTime string  `json:"create_time"`
	UpdateTime string  `json:"update_time"`

	ExtraInfo       UserExtraInfo `json:"extra_info" orm:"-"`
	ExtraInfoString string        `json:"-" orm:"type(json);column(extra_info)"`
	Status          int           `json:"status"`

	//UidStr      string `json:"uid_str" orm:"-"`
	//FrontUidStr string `json:"front_uid_str" orm:"-"`
	RequestId string `json:"-"`
}
type UserExtraInfo struct {
	Wechat         *UserExtraInfoWechat   `json:"wechat,omitempty"`
	WechatMp       *UserExtraInfoWechatMp `json:"wechat_mp,omitempty"`
	WechatOp       *UserExtraInfoWechatOp `json:"wechat_op,omitempty"`
	Cookie         *UserExtraInfoCookie   `json:"cookie,omitempty"`
	Achievement    map[string]int         `json:"achievement,omitempty"`
	LotteryChances map[string]int         `json:"lottery_chances,omitempty"`
	LotteryTimes   int                    `json:"lottery_times,omitempty"`
	ChnId          int64                  `json:"chn_id,omitempty"`
	FrontChnId     int64                  `json:"front_chn_id,omitempty"`
	Nickname       string                 `json:"nickname"`
	HeadImgUrl     string                 `json:"headimgurl"`
	//Ip string `json:"ip"`

	//Nickname string `json:"nickname"`
	//Sex      string `json:"sex"`
	//Province string `json:"province"`
	//City     string `json:"city"`
	//Country  string `json:"country"`
}
type UserExtraInfoWechat struct {
	WxAppid  string `json:"wx_appid"`
	WxOpenid string `json:"wx_openid"`
}
type UserExtraInfoWechatMp struct {
	WxAppid  string `json:"wx_appid"`
	WxOpenid string `json:"wx_openid"`
}
type UserExtraInfoWechatOp struct {
	WxOpAccount string `json:"wx_op_account"`
	WxUnionId   string `json:"wx_unionid"`
}

type UserExtraInfoCookie struct {
	Cookie string `json:"cookie"`
}

func (u *User) InsertDB() (err error) {

	newItemSlce := make([]*User, 0, 1)
	newItemSlce = append(newItemSlce, u)

	config.DefaultLogger.Debug("newItemSlce: %+v, User: %+v", newItemSlce, *u)

	_, _, err = UserInsert(newItemSlce)
	return

}
func (u *User) TxRewardMoney(tx *drsql.Tx, rewardBaseMoney float64, rewardGiftMoney float64) (err error) {

	if u.BaseMoney >= 0.0 && rewardBaseMoney >= 0.0 && u.GiftMoney >= 0.0 && rewardGiftMoney >= 0.0 {
		u.BaseMoney = u.BaseMoney + rewardBaseMoney
		u.GiftMoney = u.GiftMoney + rewardGiftMoney
	} else {
		err = fmt.Errorf("reward money, money format wrong, user: %+v, rewardBaseMoney: %f , rewardGiftMoney: %f", u, rewardBaseMoney, rewardGiftMoney)
		return
	}

	sql := "UPDATE `user` SET `base_money`= ?, `gift_money` = ?  WHERE `appid` = ? AND `uid` = ? LIMIT 1"
	//sql += "`extra_info`= JSON_SET(`extra_info`,'$.update_time', ? , $.head_url', ? ,'$.sex', ? )"

	paramSlice := make([]interface{}, 0)
	paramSlice = append(paramSlice, u.BaseMoney, u.GiftMoney, u.Appid, u.Uid)

	_, err = tx.Exec(sql, paramSlice...)
	config.DefaultLogger.Debug(" UserRewardMoney t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
	if err != nil {
		return
	}
	return

}
func (u *User) TxAddLotteryTimes(tx *drsql.Tx) (err error) {

	sql := "UPDATE `user` SET `extra_info`=JSON_SET(`extra_info`,'$.lottery_times', ?) WHERE `appid` = ? AND `uid` = ? LIMIT 1"
	//sql += "`extra_info`= JSON_SET(`extra_info`,'$.update_time', ? , $.head_url', ? ,'$.sex', ? )"

	paramSlice := make([]interface{}, 0)
	u.ExtraInfo.LotteryTimes = u.ExtraInfo.LotteryTimes + 1
	paramSlice = append(paramSlice, u.ExtraInfo.LotteryTimes, u.Appid, u.Uid)

	_, err = tx.Exec(sql, paramSlice...)
	config.DefaultLogger.Debug(" UserRewardMoney t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
	if err != nil {
		return
	}
	return

}
func UserAddLotteryChances(appid string, uid int64, payMoney float64) (err error) {
	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}
	tx, err := adminConn.Begin()
	if err != nil {
		return
	}

	rollFlag := true

	defer grdatabase.MysqlClearTransaction(tx, &rollFlag)

	u, err := GetUserForUpdate(tx, appid, uid)
	if err != nil {
		return
	}

	err = u.TxAddLotteryChances(tx, payMoney)
	if err != nil {
		return
	}

	rollFlag = false
	return
}
func (u *User) TxAddAchievement(tx *drsql.Tx, achievementKey string) (err error) {

	sql := fmt.Sprintf("UPDATE `user` SET `extra_info`=JSON_SET(`extra_info`,'$.achievement.%s', ?) WHERE `appid` = ? AND `uid` = ? LIMIT 1", achievementKey)
	//sql += "`extra_info`= JSON_SET(`extra_info`,'$.update_time', ? , $.head_url', ? ,'$.sex', ? )"

	paramSlice := make([]interface{}, 0)
	if u.ExtraInfo.Achievement == nil {

		u.ExtraInfo.Achievement = make(map[string]int)

		_, err = tx.Exec("UPDATE `user` SET `extra_info`=JSON_SET(`extra_info`,'$.achievement', JSON_MERGE('{}','{}')) WHERE `appid` = ? AND `uid` = ? LIMIT 1", u.Appid, u.Uid)
		config.DefaultLogger.Debug("TxAddAchievement t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
		if err != nil {
			return
		}

	}
	if _, ok := u.ExtraInfo.Achievement[achievementKey]; !ok {
		u.ExtraInfo.Achievement[achievementKey] = 0

	}

	u.ExtraInfo.Achievement[achievementKey] = u.ExtraInfo.Achievement[achievementKey] + 1
	paramSlice = append(paramSlice, u.ExtraInfo.Achievement[achievementKey], u.Appid, u.Uid)

	_, err = tx.Exec(sql, paramSlice...)
	config.DefaultLogger.Debug("TxAddAchievement t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
	if err != nil {
		return
	}
	return

}

func (u *User) TxAddLotteryChances(tx *drsql.Tx, payMoney float64) (err error) {
	lotteryChancesKey := fmt.Sprintf("m%d", int(payMoney))
	sql := fmt.Sprintf("UPDATE `user` SET `extra_info`=JSON_SET(`extra_info`,'$.lottery_chances.%s', ?) WHERE `appid` = ? AND `uid` = ? LIMIT 1", lotteryChancesKey)
	//sql += "`extra_info`= JSON_SET(`extra_info`,'$.update_time', ? , $.head_url', ? ,'$.sex', ? )"

	paramSlice := make([]interface{}, 0)
	if u.ExtraInfo.LotteryChances == nil {

		u.ExtraInfo.LotteryChances = make(map[string]int)

		_, err = tx.Exec("UPDATE `user` SET `extra_info`=JSON_SET(`extra_info`,'$.lottery_chances', JSON_MERGE('{}','{}')) WHERE `appid` = ? AND `uid` = ? LIMIT 1", u.Appid, u.Uid)
		config.DefaultLogger.Debug(" UserRewardMoney t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
		if err != nil {
			return
		}

	}
	if _, ok := u.ExtraInfo.LotteryChances[lotteryChancesKey]; !ok {
		u.ExtraInfo.LotteryChances[lotteryChancesKey] = 0

	}

	u.ExtraInfo.LotteryChances[lotteryChancesKey] = u.ExtraInfo.LotteryChances[lotteryChancesKey] + 1
	paramSlice = append(paramSlice, u.ExtraInfo.LotteryChances[lotteryChancesKey], u.Appid, u.Uid)

	_, err = tx.Exec(sql, paramSlice...)
	config.DefaultLogger.Debug(" UserRewardMoney t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
	if err != nil {
		return
	}
	return

}

func (u *User) CheckLotteryChances(payMoney float64) (success bool) {
	lotteryChancesKey := fmt.Sprintf("m%d", int(payMoney))

	if u.ExtraInfo.LotteryChances != nil {

		if _, ok := u.ExtraInfo.LotteryChances[lotteryChancesKey]; ok {

			if u.ExtraInfo.LotteryChances[lotteryChancesKey] >= 1 {
				return true
			}

		}
	}
	return false

}
func (u *User) TxUseLotteryChances(tx *drsql.Tx, payMoney float64) (err error) {
	lotteryChancesKey := fmt.Sprintf("m%d", int(payMoney))
	sql := fmt.Sprintf("UPDATE `user` SET `extra_info`=JSON_SET(`extra_info`,'$.lottery_chances.%s', ?) WHERE `appid` = ? AND `uid` = ? LIMIT 1", lotteryChancesKey)
	//sql += "`extra_info`= JSON_SET(`extra_info`,'$.update_time', ? , $.head_url', ? ,'$.sex', ? )"

	paramSlice := make([]interface{}, 0)
	success := u.CheckLotteryChances(payMoney)
	if success {
		u.ExtraInfo.LotteryChances[lotteryChancesKey] = u.ExtraInfo.LotteryChances[lotteryChancesKey] - 1

		paramSlice = append(paramSlice, u.ExtraInfo.LotteryChances[lotteryChancesKey], u.Appid, u.Uid)

		_, err = tx.Exec(sql, paramSlice...)
		config.DefaultLogger.Debug(" UserRewardMoney t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
		if err != nil {
			return
		}
		return

	} else {
		fmt.Errorf("no chance? user: %+v ,payMoney: %0.2f", u, payMoney)
	}
	return

}

func GetUserForUpdate(tx *drsql.Tx, appid string, uid int64) (u User, err error) {
	userSql := "SELECT `appid`,`uid`,`front_uid`,`base_money`,`gift_money`,date_format(`create_time`, '%Y-%m-%d %H:%i:%s'),date_format(`update_time`, '%Y-%m-%d %H:%i:%s'),`status`,`extra_info`  FROM `user` where `appid`= ? AND `uid`= ? limit 1 FOR UPDATE"

	err = tx.QueryRow(userSql, appid, uid).Scan(&u.Appid, &u.Uid, &u.FrontUid, &u.BaseMoney, &u.GiftMoney, &u.CreateTime, &u.UpdateTime, &u.Status, &u.ExtraInfoString)
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(u.ExtraInfoString), &u.ExtraInfo); err != nil {
		return
	}
	return
}

func UserChargeMoney(tx *drsql.Tx, appid string, uid int64, payBaseMoney float64, payGiftMoney float64) (u User, err error) {
	//var u User
	u, err = GetUserForUpdate(tx, appid, uid)
	if err != nil {
		return
	}

	if u.BaseMoney >= 0.0 && u.BaseMoney >= payBaseMoney && u.GiftMoney >= 0.0 && u.GiftMoney >= payGiftMoney {
		u.BaseMoney = u.BaseMoney - payBaseMoney
		u.GiftMoney = u.GiftMoney - payGiftMoney
	} else {
		err = fmt.Errorf("charge money, money format wrong, user: %+v, payBaseMoney: %f , payGiftMoney: %f", u, payBaseMoney, payGiftMoney)
		return
	}

	sql := "UPDATE `user` SET `base_money`= ?, `gift_money` = ?  WHERE `appid` = ? AND `uid` = ? LIMIT 1"
	//sql += "`extra_info`= JSON_SET(`extra_info`,'$.update_time', ? , $.head_url', ? ,'$.sex', ? )"

	paramSlice := make([]interface{}, 0)
	paramSlice = append(paramSlice, u.BaseMoney, u.GiftMoney, u.Appid, u.Uid)

	_, err = tx.Exec(sql, paramSlice...)
	config.DefaultLogger.Debug("UserChargeMoney t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
	if err != nil {
		return
	}

	return
}

func (u *User) UpdateBaseInfo() (err error) {
	//INSERT INTO `client` (`appid`,`client_id`,`pt_id`,`customer_id`,`create_time`, `status`,`extra_info`) VALUES('test1','test2','test3','test4','2018-05-05 00:00:00',1,'{}') ON DUPLICATE KEY UPDATE `extra_info` =JSON_SET(`extra_info`,'$.head_url','good2','$.sex', 2);
	sql := "UPDATE `user` SET `update_time`= ?, `status` = ? ,`extra_info` = JSON_SET(`extra_info`, '$.nickname', ? , '$.headimgurl', ? ) WHERE `appid` = ? AND `uid` = ? limit 1"
	//sql += "`extra_info`= JSON_SET(`extra_info`,'$.update_time', ? , $.head_url', ? ,'$.sex', ? )"

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}
	/*
		var extra_info_bytes []byte
		extra_info_bytes, err = json.Marshal(u.ExtraInfo)
		if err != nil {

			return
		}
		u.ExtraInfoString = string(extra_info_bytes)*/

	paramSlice := make([]interface{}, 0)
	paramSlice = append(paramSlice, u.UpdateTime, u.Status, u.ExtraInfo.Nickname, u.ExtraInfo.HeadImgUrl, u.Appid, u.Uid)

	_, err = adminConn.Exec(sql, paramSlice...)
	config.DefaultLogger.Debug("t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
	if err != nil {
		return
	}

	return

}
func (u *User) UpdateWechatOpInfo() (err error) {
	//INSERT INTO `client` (`appid`,`client_id`,`pt_id`,`customer_id`,`create_time`, `status`,`extra_info`) VALUES('test1','test2','test3','test4','2018-05-05 00:00:00',1,'{}') ON DUPLICATE KEY UPDATE `extra_info` =JSON_SET(`extra_info`,'$.head_url','good2','$.sex', 2);
	sql := "UPDATE `user` SET `extra_info`=JSON_SET(`extra_info`, '$.wechat_op', JSON_MERGE('{}',?)) WHERE `appid` = ? AND `uid` = ? LIMIT 1"
	//sql += "`extra_info`= JSON_SET(`extra_info`,'$.update_time', ? , $.head_url', ? ,'$.sex', ? )"

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}

	var wechatOpBytes []byte
	wechatOpBytes, err = json.Marshal(u.ExtraInfo.WechatOp)
	if err != nil {

		return
	}

	paramSlice := make([]interface{}, 0)
	paramSlice = append(paramSlice, string(wechatOpBytes), u.Appid, u.Uid)

	_, err = adminConn.Exec(sql, paramSlice...)
	config.DefaultLogger.Debug("t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
	if err != nil {
		return
	}

	return

}

func GetUserFromDBorCache(appid string, uid int64, slave bool) (newItem User, exist bool, err error) {
	err, conn := grcache.RedisConnOPGet(config.GRAPPS_REDIS, &config.GlobalConf.RedisPool)

	if err != nil {

		return

	}
	defer conn.Close()

	redisKey := fmt.Sprintf(REDIS_PREFIX_USER_FOR_PK, appid, uid)
	if slave {

		exist, err = conn.Get(redisKey, &newItem)
		if err != nil {
			return
		}
		if exist {
			return
		}
	}

	sql := "SELECT `appid`,`uid`,`front_uid`,`base_money`,`gift_money`,date_format(`create_time`, '%Y-%m-%d %H:%i:%s'),date_format(`update_time`, '%Y-%m-%d %H:%i:%s'),`status`,`extra_info`  FROM `user` where `appid`= ? AND `uid`= ?  limit 1"

	paramSlice := make([]interface{}, 0)
	paramSlice = append(paramSlice, appid)
	paramSlice = append(paramSlice, uid)

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
		err = rows.Scan(&newItem.Appid, &newItem.Uid, &newItem.FrontUid, &newItem.BaseMoney, &newItem.GiftMoney, &newItem.CreateTime, &newItem.UpdateTime, &newItem.Status, &newItem.ExtraInfoString)
		if err != nil {
			return
		}
		if err = json.Unmarshal([]byte(newItem.ExtraInfoString), &newItem.ExtraInfo); err != nil {
			return
		}
		exist = true

	}

	return

}
func (u *User) SaveDB() (err error) {

	newItemSlce := make([]*User, 0, 1)
	newItemSlce = append(newItemSlce, u)

	config.DefaultLogger.Debug("newItemSlce: %+v, User: %+v", newItemSlce, *u)

	_, _, err = UserInsertOnDuplicateKeyUpdate(newItemSlce)
	return

}
func UserInsert(newItemSlice []*User) (successSlice []string, successNum int, err error) {
	//INSERT INTO `client` (`appid`,`client_id`,`pt_id`,`customer_id`,`create_time`, `status`,`extra_info`) VALUES('test1','test2','test3','test4','2018-05-05 00:00:00',1,'{}') ON DUPLICATE KEY UPDATE `extra_info` =JSON_SET(`extra_info`,'$.head_url','good2','$.sex', 2);
	sql := "INSERT INTO `user` (`appid`,`uid`,`front_uid`,`base_money`,`gift_money`,`create_time`,`update_time`,`status`,`extra_info`) "
	sql += "VALUES(? ,? ,? ,? ,? ,? ,? ,? ,?)"
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
		paramSlice = append(paramSlice, newItem.Appid, newItem.Uid, newItem.FrontUid, newItem.BaseMoney, newItem.GiftMoney, newItem.CreateTime, newItem.UpdateTime, newItem.Status, newItem.ExtraInfoString)

		// newItem.Appid, newItem.ClientId, newItem.PtId, newItem.CustomerId, newItem.CreateTime, newItem.Status, extra_info_string,
		//newItem.CustomerId,
		//newItem.ExtraInfo.UpdateTime, newItem.ExtraInfo.HeadUrl, newItem.ExtraInfo.Sex
		var ret drsql.Result

		ret, err = adminConn.Exec(sql, paramSlice...)
		config.DefaultLogger.Debug("t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
		if err != nil {
			return
		}
		var rowsAffected int64
		rowsAffected, err = ret.RowsAffected()
		if err != nil {
			return
		}
		config.DefaultLogger.Debug("rowsAffected: %d", rowsAffected)

		successSlice = append(successSlice, newItem.RequestId)
		successNum += 1
	}

	return

}

func UserInsertOnDuplicateKeyUpdate(newItemSlice []*User) (successSlice []string, successNum int, err error) {
	//INSERT INTO `client` (`appid`,`client_id`,`pt_id`,`customer_id`,`create_time`, `status`,`extra_info`) VALUES('test1','test2','test3','test4','2018-05-05 00:00:00',1,'{}') ON DUPLICATE KEY UPDATE `extra_info` =JSON_SET(`extra_info`,'$.head_url','good2','$.sex', 2);
	sql := "UPDATE INTO `user` (`appid`,`uid`,`front_uid`,`base_money`,`gift_money`,`create_time`,`update_time`,`status`,`extra_info`) "
	sql += "VALUES(?,? ,? ,? ,? ,? ,? ,? ,?) ON DUPLICATE KEY UPDATE `update_time`= VALUES(`update_time`),`extra_info` = VALUES(`extra_info`)"
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
		paramSlice = append(paramSlice, newItem.Appid, newItem.Uid, newItem.FrontUid, newItem.BaseMoney, newItem.GiftMoney, newItem.CreateTime, newItem.UpdateTime, newItem.Status, newItem.ExtraInfoString)

		// newItem.Appid, newItem.ClientId, newItem.PtId, newItem.CustomerId, newItem.CreateTime, newItem.Status, extra_info_string,
		//newItem.CustomerId,
		//newItem.ExtraInfo.UpdateTime, newItem.ExtraInfo.HeadUrl, newItem.ExtraInfo.Sex
		var ret drsql.Result

		ret, err = adminConn.Exec(sql, paramSlice...)
		config.DefaultLogger.Debug("t_sql: %s,paramSlice: %+v, len: %d,err: %s", sql, paramSlice, len(paramSlice), err)
		if err != nil {
			return
		}
		var rowsAffected int64
		rowsAffected, err = ret.RowsAffected()
		if err != nil {
			return
		}
		config.DefaultLogger.Debug("rowsAffected: %d", rowsAffected)

		successSlice = append(successSlice, newItem.RequestId)
		successNum += 1
	}

	return

}

type CookieUser struct {
	Appid           string              `json:"-"`
	Uid             int64               `json:"uid"`
	Cookie          string              `json:"cookie"`
	CreateTime      string              `json:"create_time"`
	ExtraInfo       CookieUserExtraInfo `json:"extra_info"`
	ExtraInfoString string              `json:"-"`
	RequestId       string              `json:"-"`
}
type CookieUserExtraInfo struct {
	UpdateTime *string `json:"update_time"`
}

func (uo *CookieUser) InsertDB() (err error) {

	newItemSlce := make([]*CookieUser, 0, 1)
	newItemSlce = append(newItemSlce, uo)

	config.DefaultLogger.Debug("newItemSlce: %+v, CookieUser: %+v", newItemSlce, *uo)

	_, _, err = CookieUserInsert(newItemSlce)
	return

}

func GetCookieUserFromDBorCache(appid string, cookie string, slave bool) (newItem CookieUser, exist bool, err error) {
	if slave {

	}

	sql := "SELECT `appid`,`uid`,`cookie`,date_format(`create_time`, '%Y-%m-%d %H:%i:%s'),`extra_info` FROM `cookie_user` where `appid`= ? AND `cookie`= ? limit 1"

	paramSlice := make([]interface{}, 0)
	paramSlice = append(paramSlice, appid)
	paramSlice = append(paramSlice, cookie)

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
		err = rows.Scan(&newItem.Appid, &newItem.Uid, &newItem.Cookie, &newItem.CreateTime, &newItem.ExtraInfoString)
		if err != nil {
			return
		}
		if err = json.Unmarshal([]byte(newItem.ExtraInfoString), &newItem.ExtraInfo); err != nil {
			return
		}
		exist = true

	}

	return

}
func CookieUserInsert(newItemSlice []*CookieUser) (successSlice []string, successNum int, err error) {
	//INSERT INTO `client` (`appid`,`client_id`,`pt_id`,`customer_id`,`create_time`, `status`,`extra_info`) VALUES('test1','test2','test3','test4','2018-05-05 00:00:00',1,'{}') ON DUPLICATE KEY UPDATE `extra_info` =JSON_SET(`extra_info`,'$.head_url','good2','$.sex', 2);
	sql := "INSERT  INTO `cookie_user` (`appid`,`uid`, `cookie`,`create_time`,`extra_info`) "
	sql += "VALUES(?, ?, ?, ? ,? )"
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
		paramSlice = append(paramSlice, newItem.Appid, newItem.Uid, newItem.Cookie, newItem.CreateTime, newItem.ExtraInfoString)

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
