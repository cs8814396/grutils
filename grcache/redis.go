package grcache

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

var MasterRedisPool *redis.Pool

//var RedisPoolMap map[string]*redis.Pool

var RedisPoolMap sync.Map

var (
	ErrCacheMiss = errors.New("CacheMiss")
)

type XmlRedis struct {
	Addr             string        `xml:"addr"`
	Db               int           `xml:"db"`
	Pass             string        `xml:"pass"`
	ConnectTimeoutMs time.Duration `xml:"connecttimeoutms"`
	ReadTimeoutMs    time.Duration `xml:"readtimeoutms"`
	WriteTimeoutMs   time.Duration `xml:"writetimeoutms"`
	IdleTimeoutMs    time.Duration `xml:"idletimeoutms"`
	MaxIdle          int           `xml:"maxidle"`
	KeyExpireSec     int           `xml:"keyexpiresec",omitempty` //must be omitempty or you will get error from where it was used before
	DbWrite          int           `xml:"dbwrite",omitempty`
	ReadOff          int           `xml:"readoff,omitempty"`
}

func init() {
	//RedisPoolMap = make(map[string]*redis.Pool)
}

func RedisConnGet(poolName string, poolConfig XmlRedis) (err error, conn redis.Conn) {
	if rp, tok := RedisPoolMap.Load(poolName); tok {
		err = nil
		conn = (rp.(*redis.Pool)).Get()
		if conn == nil {
			err = errors.New("redis conn get from pool is nil, poolName: " + poolName)
		}
		return
	}
	err, newRp := InitRedisCacheConn(poolConfig)
	if err != nil {

		conn = nil

		return

	} else {
		conn = newRp.Get()
		RedisPoolMap.Store(poolName, newRp)

		return

	}

}

func RedisSetEx(conn redis.Conn, key string, value interface{}, expireSecond int) (err error) {

	err = nil

	var buffer bytes.Buffer

	enc := gob.NewEncoder(&buffer)

	err0 := enc.Encode(value)
	if err0 != nil {
		err = errors.New("Encode error key=" + key + " err=" + err0.Error())
		return
	}

	_, err = conn.Do("SET", key, buffer.Bytes(), "EX", expireSecond)
	if err != nil {
		err = errors.New("redis do set error" + key + " err=" + err.Error())

		return
	}

	//newItem := &memcache.Item{Key: md5Key, Value: buffer.Bytes(), Expiration: int32(et)}

	//startTime := time.Now().UnixNano()

	// globalGoMetric.AddValue("cache_write", float32(time.Now().UnixNano()-startTime)/1000000)
	//l4g.Debug("MCSet ok key=%s md5key=%s value=%+v ", key, md5Key, value)

	return
}

func RedisGet(conn redis.Conn, key string, value interface{}) (exist bool, err error) {
	exist = false
	err = nil
	result, err0 := conn.Do("GET", key)
	if err0 != nil {

		err = errors.New("redis do get error key: " + key + " err=" + err.Error())

		//config.DefaultLogger.Error("redis get error err: %s, key: %s", err, WECHAT_ACCESSTOKEN_REDIS_KEY_PREFIX+wechatAccessToken)
		return
	}

	if result != nil {

		buf := bytes.NewBuffer(result.([]byte))
		dec := gob.NewDecoder(buf)

		err0 = dec.Decode(value)
		if err0 != nil {
			err = errors.New("Decode error key=" + key + " err=" + err0.Error())
			return
		}
		exist = true

	}

	return

}

func InitRedisCacheConn(poolConfig XmlRedis) (err error, newRp *redis.Pool) {

	newRp = redisPool(poolConfig.Addr, poolConfig.Db, poolConfig.Pass, poolConfig.ConnectTimeoutMs,
		poolConfig.ReadTimeoutMs, poolConfig.WriteTimeoutMs, poolConfig.IdleTimeoutMs, poolConfig.MaxIdle)

	err = testRedisPool(newRp)

	//var mmerr error
	//MasterMemCache, mmerr = NewMemcache(config.GlobalConf.MasterMemCache.Addr, time.Duration(config.GlobalConf.MasterMemCache.RWTimeout)*time.Millisecond)

	return err, newRp //&& mmerr == nil
}

func redisPool(addr string, db int, pass string, connectTimeoutMs time.Duration, readTimeoutMs time.Duration, writeTimeoutMs time.Duration, idleTimeoutMs time.Duration, maxIdle int) *redis.Pool {

	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeoutMs * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", addr, time.Millisecond*connectTimeoutMs, time.Millisecond*readTimeoutMs, time.Millisecond*writeTimeoutMs)
			if err != nil {
				//l4g.Error("redis connect %s fail", addr)
				return nil, err
			}
			if pass != "" {
				_, err = c.Do("AUTH", pass)
				if err != nil {
					//l4g.Error("redis password error: %s.", err.Error())
					return nil, err
				}
			}

			err = c.Send("SELECT", db)
			if err != nil {
				//l4g.Error("redis select database %d fail, error: %s", db, err.Error())
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func testRedisPool(pool *redis.Pool) (err error) {
	c := pool.Get()
	defer c.Close()

	err = pool.TestOnBorrow(c, time.Now())
	if err != nil {
		return
	}

	return nil
}
