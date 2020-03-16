package grcache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/bradfitz/gomemcache/memcache"

	"github.com/gdgrc/grutils/grmath"
	//"sync"
	"time"
)

var MasterMemCache *memcache.Client

func GetMemCacheMaster() *memcache.Client {
	return MasterMemCache
}

func NewMemcache(addr string, rwtimeout time.Duration) (client *memcache.Client, err error) {

	err = nil

	ss := new(memcache.ServerList)
	err = ss.SetServers(addr)

	if err != nil {
		return
	}

	client = memcache.NewFromSelector(ss)

	if rwtimeout != 0 {
		client.Timeout = rwtimeout
	}

	return
}

func DBCacheSet(key string, value interface{}, expire int) (err error) {

	err = nil
	mm := GetMemCacheMaster()
	if mm == nil {
		err = errors.New("GetMemcache error key=" + key)
		return
	}

	var buffer bytes.Buffer

	enc := gob.NewEncoder(&buffer)

	err0 := enc.Encode(value)
	if err0 != nil {
		err = errors.New("Encode error key=" + key + " err=" + err0.Error())
		return
	}

	et := expire

	md5Key := grmath.Md5(key)
	newItem := &memcache.Item{Key: md5Key, Value: buffer.Bytes(), Expiration: int32(et)}

	//startTime := time.Now().UnixNano()
	err0 = mm.Set(newItem)
	if err0 != nil {
		err = errors.New("MCSet key=" + key + " md5key=" + md5Key + " err=" + err0.Error())
		return
	}
	// globalGoMetric.AddValue("cache_write", float32(time.Now().UnixNano()-startTime)/1000000)
	//l4g.Debug("MCSet ok key=%s md5key=%s value=%+v ", key, md5Key, value)

	return
}

func DBCacheGet(key string, value interface{}) (err error) {

	err = nil
	mm := GetMemCacheMaster()
	if mm == nil {
		err = errors.New("GetMemcache error key=" + key)
		return
	}

	md5Key := grmath.Md5(key)

	//startTime := time.Now().UnixNano()
	item, err0 := mm.Get(md5Key)
	if err0 != nil && err0 != memcache.ErrCacheMiss {
		err = errors.New("MCGet key=" + key + " md5key=" + md5Key + " err=" + err0.Error())
		return
	}
	// globalGoMetric.AddValue("cache_write", float32(time.Now().UnixNano()-startTime)/1000000)

	if err0 == memcache.ErrCacheMiss {
		err = ErrCacheMiss
		return
	}

	if item.Key != md5Key {
		err = errors.New("MCGet key not match key=" + key + " md5key=" + md5Key + " rspkey=" + item.Key)
		return
	}

	buf := bytes.NewBuffer(item.Value)
	dec := gob.NewDecoder(buf)
	err0 = dec.Decode(value)
	if err0 != nil {
		err = errors.New("Decode error key=" + key + " err=" + err0.Error())
		return
	}

	return

}
