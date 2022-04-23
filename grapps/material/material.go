package material

import (
	"encoding/json"
	//"grutils/grmath"
	"fmt"

	"github.com/gdgrc/grutils/grapps/config"
	"github.com/gdgrc/grutils/grcache"
)

const (
	REDIS_PREFIX_MATERIAL_FOR_PK = "RedisMaterial?appid=%s&type=%d&mid=%d"
)

type Material struct {
	Appid      string `json:"appid"`
	Type       int    `json:"type"`
	Mid        int64  `json:"mid"`
	Status     int    `json:"status"`
	CreateTime string `json:"create_time"`

	ExtraInfoString string            `json:"-"`
	ExtraInfo       MaterialExtraInfo `json:"extra_info"`
}
type MaterialExtraInfo struct {
	RemoteUrl *string `json:"remote_url,omitempty"`
}

const (
	MATERIAL_TYPE_FISSION_VIDEO = 1
)

func GetMaterialFromDBorCache(appid string, iType int, mid int64, slave bool) (newItem Material, exist bool, err error) {
	err, conn := grcache.RedisConnOPGet(config.GRAPPS_REDIS, &config.GlobalConf.RedisPool)

	if err != nil {

		return

	}
	defer conn.Close()

	redisKey := fmt.Sprintf(REDIS_PREFIX_MATERIAL_FOR_PK, appid, iType, mid)
	if slave {

		exist, err = conn.Get(redisKey, &newItem)
		if err != nil {
			return
		}
		if exist {
			return
		}
	}

	sql := "SELECT `appid`,`type`,`mid`,`status` ,`create_time`, `extra_info`"
	sql += " FROM `material` where `appid`= ? AND `type`= ? AND `mid` = ?"

	paramSlice := make([]interface{}, 0)
	paramSlice = append(paramSlice, appid)
	paramSlice = append(paramSlice, iType)

	paramSlice = append(paramSlice, mid)

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
		err = rows.Scan(&newItem.Appid, &newItem.Type, &newItem.Mid, &newItem.Status, &newItem.CreateTime, &newItem.ExtraInfoString)
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
