package grdatabase

import (
	"errors"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*type MysqlPoolMap struct {
	//mutex   *sync.RWMutex
	pool        sync.Map //map[string]*sql.DB
	dnsPool     sync.Map
	maxidlePool sync.Map
}*/

var ormPool sync.Map

func OrmGetConn(dbname string, dsn string, maxidleconn int) (g *gorm.DB, err error) {

	if idb, ok := ormPool.Load(dbname); ok {
		g = idb.(*gorm.DB)

		return
	}

	g, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		err = errors.New("OrmGetConn open fail err=" + err.Error())
		return
	}
	ormPool.Store(dbname, g)
	return

}
