package version

import (
	//drsql "database/sql"
	"encoding/json"

	"github.com/gdgrc/grutils/grapps/config"
	//"grutils/grmath"
	//"math/rand"
	//"time"
)

type Version struct {
	Appid     string `json:"appid" orm:"pk"`
	VersionId int    `json:"version_id"`

	CreateTime      string           `json:"create_time"`
	Status          int              `json:"status"`
	ExtraInfo       VersionExtraInfo `json:"extra_info" orm:"-"`
	ExtraInfoString string           `json:"-" orm:"column(extra_info);type(json)"`
	RequestId       string           `json:"-" orm:"-"`
}
type VersionExtraInfo struct {
	Desc   *string `json:"desc,omitempty"`
	Danger *int    `json:"danger,omitempty"`
}

func GetVersionFromDBorCache(appid string, versionId int, slave bool) (newItem Version, exist bool, err error) {

	if slave {
		//get from cache and return if exists
	}

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}
	//date_format(`create_time`, '%Y-%m-%d %H:%i:%s')
	sql := "SELECT `appid`,`version_id`, date_format(`create_time`, '%Y-%m-%d %H:%i:%s'),`status`,`extra_info` FROM `version` where `appid` = ? AND `version_id` = ?"

	rows, err := adminConn.Query(sql, appid, versionId)
	if err != nil {

		return
	}

	defer rows.Close()

	for rows.Next() {

		//fmt.Println("goood")

		err = rows.Scan(&newItem.Appid, &newItem.VersionId, &newItem.CreateTime, &newItem.Status, &newItem.ExtraInfoString)
		if err != nil {
			//config.DefaultLogger.Error("admin query  fail, err: %s", err)
			return
		}

		if err = json.Unmarshal([]byte(newItem.ExtraInfoString), &newItem.ExtraInfo); err != nil {
			return
		}
		exist = true

	}

	return

}
