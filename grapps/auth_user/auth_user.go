package auth_user

import (
	"encoding/json"

	"github.com/gdgrc/grutils/grapps/application"
	"github.com/gdgrc/grutils/grapps/config"
)

//update `auth_user` set `extra_info`='{"applications":["xxx"]}' where id= xxx limit 1;

type AuthUser struct {
	Id              int               `orm:"pk"`
	Username        string            `orm:"column(username)"`
	Password        string            `orm:"column(password)"`
	UpdateTime      string            `json:"update_time"`
	CreateTime      string            `json:"create_time"`
	Status          int               `orm:"column(status)"`
	ExtraInfo       AuthUserExtraInfo `json:"extra_info" orm:"-"`
	ExtraInfoString string            `json:"-" orm:"type(json);column(extra_info)"`
}

type AuthUserExtraInfo struct {
	IsSuperUser  int      `json:"is_superuser"`
	IsStaff      int      `json:"is_staff"`
	IsActive     int      `json:"is_active"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Email        string   `json:"email"`
	Applications []string `json:"applications,omitempty"`

	AppChn map[string][]int `json:"app_chn,omitempty"`
	//Channels     []int64  `json:"channels,omitempty"`
}

func (this *AuthUser) GetApplications() (newItemSlice []application.Application, err error) {

	//appSlice := make([]string)

	return application.GetApplications(this.ExtraInfo.Applications, true)

}

// insert into `auth_user` (`username`,`password`,`status`,`create_time`,`update_time`,`extra_info`) VALUES('aaa','123',1,'2018-01-01 00:00:00','2018-01-01 00:00:00','{}');
func GetAuthUserFromDBorCacheById(id int, slave bool) (newItem AuthUser, err error) {

	if slave {
		//get from cache and return if exists
	}

	sql := "SELECT `id`,`username`,`password`,date_format(`update_time`, '%Y-%m-%d %H:%i:%s'), date_format(`create_time`, '%Y-%m-%d %H:%i:%s'),`status`,`extra_info` FROM `auth_user` where `id`=? "

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}

	var extraInfoStr string
	err = adminConn.QueryRow(sql, id).Scan(&newItem.Id, &newItem.Username, &newItem.Password, &newItem.UpdateTime, &newItem.CreateTime, &newItem.Status, &extraInfoStr)
	if err != nil {
		//config.DefaultLogger.Error("admin query  fail, err: %s", err)
		return
	}
	if err = json.Unmarshal([]byte(extraInfoStr), &newItem.ExtraInfo); err != nil {
		return
	}

	return

}
func GetAuthUserFromDBorCacheByUserPwd(username string, password string, slave bool) (newItem AuthUser, err error) {

	if slave {
		//get from cache and return if exists
	}

	sql := "SELECT `id`,`username`,`password`,date_format(`update_time`, '%Y-%m-%d %H:%i:%s'), date_format(`create_time`, '%Y-%m-%d %H:%i:%s'),`status`,`extra_info` FROM `auth_user` where `username`=? AND `password` = ?"

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}

	var extraInfoStr string
	err = adminConn.QueryRow(sql, username, password).Scan(&newItem.Id, &newItem.Username, &newItem.Password, &newItem.UpdateTime, &newItem.CreateTime, &newItem.Status, &extraInfoStr)
	if err != nil {
		//config.DefaultLogger.Error("admin query  fail, err: %s", err)
		return
	}
	if err = json.Unmarshal([]byte(extraInfoStr), &newItem.ExtraInfo); err != nil {
		return
	}

	return

}
