package application

import (
	"encoding/json"
	"fmt"

	"github.com/gdgrc/grutils/grapps/config"
)

type Application struct {
	Appid  string `json:"appid" orm:"pk"`
	Appkey string `json:"appkey"`

	CreateTime      string               `json:"create_time" orm:"column(create_time)"`
	Status          int                  `json:"status"`
	ExtraInfo       ApplicationExtraInfo `json:"extra_info" orm:"-"`
	ExtraInfoString string               `json:"-" orm:"type(json);column(extra_info)"`
}

type ApplicationExtraInfo struct {
	IsProxyParam  *int    `json:"is_proxy_param,omitempty"` // it shows params for wechat applications vary from proxy
	WxAppid       *string `json:"wx_appid,omitempty"`
	WxAppSecret   *string `json:"wx_appsecret,omitempty"`
	WxMpAppid     *string `json:"wxmp_appid,omitempty"`
	WxMpAppSecret *string `json:"wxmp_appsecret,omitempty"`
	WxOpAccount   *string `json:"wx_op_account,omitempty"`
	WxOpSecret    *string `json:"wx_op_secret,omitempty"`
}

func GetValidApplicationFromDBorCache(slave bool) (newItemSlice []Application, err error) {

	if slave {
		//get from cache and return if exists
	}

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}
	//date_format(`create_time`, '%Y-%m-%d %H:%i:%s')
	sql := "SELECT `appid`,`appkey`, date_format(`create_time`, '%Y-%m-%d %H:%i:%s'),`status`,`extra_info` FROM `application` where `status` = 1 "

	rows, err := adminConn.Query(sql)
	if err != nil {

		return
	}

	defer rows.Close()

	newItemSlice = make([]Application, 0)

	for rows.Next() {

		//fmt.Println("goood")

		var newItem Application
		var extraInfoStr string
		err = rows.Scan(&newItem.Appid, &newItem.Appkey, &newItem.CreateTime, &newItem.Status, &extraInfoStr)
		if err != nil {
			//config.DefaultLogger.Error("admin query  fail, err: %s", err)
			return
		}

		//fmt.Println("asasdasdasdasdasdasd %+v", newItem)
		if err = json.Unmarshal([]byte(extraInfoStr), &newItem.ExtraInfo); err != nil {
			return
		}
		newItemSlice = append(newItemSlice, newItem)

	}

	return

}
func GetApplications(appidSlice []string, slave bool) (newItemSlice []Application, err error) {

	if slave {
		//get from cache and return if exists
	}

	if appidSlice == nil {
		return
	}

	jsonBytes, err := json.Marshal(appidSlice)
	if err != nil {
		return
	}

	jsonString := string(jsonBytes)
	jsonStringLength := len(jsonString)

	if jsonStringLength < 2 {
		err = fmt.Errorf("jsonStringLength is less than 2? jsonString: %S", jsonString)
		return
	}

	//fmt.Println(string(a[:]), appidSlice == nil, a == nil, a, err)

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}
	//date_format(`create_time`, '%Y-%m-%d %H:%i:%s')
	sql := "SELECT `appid`,`appkey`, date_format(`create_time`, '%Y-%m-%d %H:%i:%s'),`status`,`extra_info` FROM `application` where `appid` in "
	sql += fmt.Sprintf("(%s)", jsonString[1:jsonStringLength-1])

	rows, err := adminConn.Query(sql)
	if err != nil {

		return
	}

	defer rows.Close()

	newItemSlice = make([]Application, 0)

	for rows.Next() {

		//fmt.Println("goood")

		var newItem Application
		var extraInfoStr string
		err = rows.Scan(&newItem.Appid, &newItem.Appkey, &newItem.CreateTime, &newItem.Status, &extraInfoStr)
		if err != nil {
			//config.DefaultLogger.Error("admin query  fail, err: %s", err)
			return
		}

		//fmt.Println("asasdasdasdasdasdasd %+v", newItem)
		if err = json.Unmarshal([]byte(extraInfoStr), &newItem.ExtraInfo); err != nil {
			return
		}
		newItemSlice = append(newItemSlice, newItem)

	}

	return

}

func GetApplicationFromDBorCache(appid string, slave bool) (newItem Application, err error) {

	if slave {
		//get from cache and return if exists
	}

	sql := "SELECT `appid`,`appkey`, date_format(`create_time`, '%Y-%m-%d %H:%i:%s'),`status`,`extra_info` FROM `application` where `appid`=? "

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}

	var extraInfoStr string
	err = adminConn.QueryRow(sql, appid).Scan(&newItem.Appid, &newItem.Appkey, &newItem.CreateTime, &newItem.Status, &extraInfoStr)
	if err != nil {
		//config.DefaultLogger.Error("admin query  fail, err: %s", err)
		return
	}
	if err = json.Unmarshal([]byte(extraInfoStr), &newItem.ExtraInfo); err != nil {
		return
	}

	return

}
