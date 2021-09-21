package config

import (
	"database/sql"

	"github.com/gdgrc/grutils/grcache"
	"github.com/gdgrc/grutils/grdatabase"
	"github.com/gdgrc/grutils/grfile"
	"github.com/gdgrc/grutils/grthird"
	"gorm.io/gorm"
)

const (
	DATA_ADMIN         = "dataadmin"
	PROXY_ADMIN        = "proxyadmin"
	MATERIAL_LOG_ADMIN = "materiallogadmin"
	GRAPPS_REDIS       = "grapps_redis"

	DEFAULT_DATETIME = "1970-01-01 00:00:00"
)

type XmlServer struct {
	BindAddr string `xml:"bindaddr"`
	Host     string `xml:"host"`
	Debug    bool   `xml:"debug"`
}

type RotateSwitch struct {
	FrontUserStat           int `xml:"fronuserstat"`
	ProxyUserGobSessionStat int `xml:"proxyusergodsessionstat"`
	ProxyStatusStat         int `xml:"proxystatusstat"`
}
type XXWeiXin struct {
	Url string `xml:"url"`
}

type ConsulConfig struct {
	Host string `xml:"host"`
	Port int    `xml:"port"`
}

type AesConfig struct {
	AesKey string `xml:"aeskey"`
	AesIV  string `xml:"aesiv"`
}
type Rc4Config struct {
	Rc4Key string `xml:"rc4key"`
}

type xmlConfig struct {
	//XMLName    xml.Name  `xml:"config"`
	Server              XmlServer              `xml:"server"`
	DefaultLog          grfile.XmlLogger       `xml:"defaultlog"`
	DataAdmin           grdatabase.XmlMysqlDsn `xml:"dataadmin"`
	ProxyAdmin          grdatabase.XmlMysqlDsn `xml:"proxyadmin"`
	MaterialLogAdmin    grdatabase.XmlMysqlDsn `xml:"materiallogadmin"`
	RedisPool           grcache.XmlRedis       `xml:"redispool"`
	AlertOver           grthird.XmlAlertOver   `xml:"alertover"`
	Wechat              grthird.XmlWechat      `xml:"wechat"`
	WechatMp            grthird.XmlWechat      `xml:"wechat_mp"`
	RTSwitch            RotateSwitch           `xml:"rotate_switch"`
	XxWx                XXWeiXin               `xml:"xxweixin"`
	DefaultConsulConfig ConsulConfig           `xml:"defaultconsulconfig"`
	Aes                 AesConfig              `xml:"aes"`
	Rc4                 Rc4Config              `xml:"rc4"`
}

func MaterialLogAdminGet(slave bool) (db *sql.DB, err error) {

	dbname := MATERIAL_LOG_ADMIN
	dsn := GlobalConf.MaterialLogAdmin.AdminDsn
	maxIdelConn := GlobalConf.MaterialLogAdmin.MaxIdleConn
	if slave {
		dbname = dbname + "_slave"
	}
	return DefaultMysqlPool.DBGetConn(dbname, dsn, maxIdelConn)

}

func DataAdminGet(slave bool) (db *sql.DB, err error) {

	dbname := DATA_ADMIN
	dsn := GlobalConf.DataAdmin.AdminDsn
	maxIdelConn := GlobalConf.DataAdmin.MaxIdleConn
	if slave {
		dbname = dbname + "_slave"
	}
	return DefaultMysqlPool.DBGetConn(dbname, dsn, maxIdelConn)

}

func DataAdminOrmGet(slave bool) (db *gorm.DB, err error) {

	dbname := DATA_ADMIN
	dsn := GlobalConf.DataAdmin.AdminDsn
	maxIdelConn := GlobalConf.DataAdmin.MaxIdleConn
	if slave {
		dbname = dbname + "_slave"
	}
	return grdatabase.OrmGetConn(dbname, dsn, maxIdelConn)

}

func ProxyAdminGet(slave bool) (db *sql.DB, err error) {

	dbname := PROXY_ADMIN
	dsn := GlobalConf.ProxyAdmin.AdminDsn
	maxIdelConn := GlobalConf.ProxyAdmin.MaxIdleConn
	if slave {
		dbname = dbname + "_slave"
	}
	return DefaultMysqlPool.DBGetConn(dbname, dsn, maxIdelConn)

}
