package config

// config
import (
	"fmt"
	//"grutils/grcache"
	"github.com/gdgrc/grutils/grdatabase"
	"github.com/gdgrc/grutils/grfile"
)

// config
var GlobalConf xmlConfig

//default log

var DefaultLogger *grfile.Logger

var DefaultMysqlPool grdatabase.MysqlPoolMap

func BaseInit() bool {
	DefaultLogger = new(grfile.Logger)
	err := grfile.MakeFilePathDirIfNotExist(GlobalConf.DefaultLog.LogFile)
	if err != nil {
		panic("grutils BaseInit err: " + err.Error())
	}

	DefaultLogger.CreateLogger(GlobalConf.DefaultLog.LogFile, "grdfl", GlobalConf.DefaultLog.Level)

	DefaultLogger.SetMonitorAlertOver(GlobalConf.AlertOver)

	//grthird.XXWeiXinUrl = GlobalConf.XxWx.Url

	// init db

	/*

		_, err := DataAdminGet(false)
		if err != nil {
			DefaultLogger.Error("DataAdminGet error err: %s,config: %+v", err, GlobalConf.DataAdmin)
			return false
		}*/

	/*
		_, err = ProxyAdminGet(false)
		if err != nil {
			DefaultLogger.Error("ProxyAdminGet error err: %s,config: %+v", err, GlobalConf.ProxyAdmin)
			return false
		}*/
	/*
		err, conn := grcache.RedisConnOPGet(GRAPPS_REDIS, &GlobalConf.RedisPool)

		if err != nil {
			DefaultLogger.Error("RedisConnOPGet error err: %s,config: %+v", err, GlobalConf.RedisPool)

			return false

		}
		defer conn.Close()*/

	/*

		DefaultMysqlPool.DBGetConn(MSP_ADMIN, GlobalConf.ArmChairAdmin.AdminDsn, GlobalConf.ArmChairAdmin.MaxIdleConn, false)

		DefaultMysqlPool.DBGetConn(MSP_ADMIN, GlobalConf.ArmChairAdmin.AdminDsn, GlobalConf.ArmChairAdmin.MaxIdleConn, true)

		DefaultMysqlPool.DBGetConn(MSP_USER, GlobalConf.ArmChairUser.AdminDsn, GlobalConf.ArmChairUser.MaxIdleConn, false)

		DefaultMysqlPool.DBGetConn(MSP_USER, GlobalConf.ArmChairUser.AdminDsn, GlobalConf.ArmChairUser.MaxIdleConn, true)*/

	//log.Println("conf load success")
	return true

}
func InitContents(contents []byte) bool {

	err := grfile.LoadXmlConfigWithContents(contents, &GlobalConf)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return BaseInit()
}

func Init(filename string) bool {

	_, err := grfile.LoadXmlConfig(filename, &GlobalConf)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return BaseInit()

}
