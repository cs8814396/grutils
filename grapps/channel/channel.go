package channel

import (
	"database/sql"
	"encoding/json"
)

type Channel struct {
	Appid string `json:"appid" orm:"pk;column(appid)"`
	ChnId int64  `json:"chn_id" orm:"column(chn_id)"`

	CreateTime string `json:"create_time" orm:"column(create_time)"`

	ExtraInfo       ChannelExtraInfo `json:"extra_info" orm:"-"`
	ExtraInfoString string           `json:"-" orm:"type(json);column(extra_info)"`
	Status          int              `json:"status" orm:"column(status)"`
	RequestId       string           `json:"-" orm:"-"`
}
type ChannelExtraInfo struct {
	HiddenChance *int `json:"hidden_chance,omitempty"`
}

func GetChannelForUpdate(tx *sql.Tx, appid string, chnId int64) (u Channel, exist bool, err error) {
	selectSql := "SELECT `appid`,`chn_id`,date_format(`create_time`, '%Y-%m-%d %H:%i:%s'),`status`,`extra_info`  FROM `channel` where `appid`= ? AND `chn_id`= ? limit 1 FOR UPDATE"

	err = tx.QueryRow(selectSql, appid, chnId).Scan(&u.Appid, &u.ChnId, &u.CreateTime, &u.Status, &u.ExtraInfoString)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		return
	}
	if err = json.Unmarshal([]byte(u.ExtraInfoString), &u.ExtraInfo); err != nil {
		return
	}

	exist = true
	return
}
