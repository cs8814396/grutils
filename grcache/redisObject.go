package grcache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"

	"sync"
)

var ORedisPoolMap sync.Map

const (
	REDIS_NIL = "REDIS::OBJECT::NIL"
)

func init() {

}

type RedisConn struct {
	Conn       redis.Conn
	PoolConfig XmlRedis
}

func RedisConnOPGet(poolName string, poolConfig *XmlRedis) (err error, rc RedisConn) {
	if rp, tok := ORedisPoolMap.Load(poolName); tok {
		err = nil

		redisMap := rp.(map[string]interface{})

		rc.PoolConfig = redisMap["poolConfig"].(XmlRedis)
		rc.Conn = (redisMap["redisPool"].(*redis.Pool)).Get()

		if rc.Conn == nil {
			err = errors.New("redis conn get from pool is nil, poolName: " + poolName)
		}

		return
	}

	if poolConfig != nil {

		var newRp *redis.Pool

		err, newRp = InitRedisCacheConn(*poolConfig)
		if err != nil {

			rc.Conn = nil

			return

		} else {

			rc.Conn = newRp.Get()

			newMap := make(map[string]interface{})
			newMap["redisPool"] = newRp
			newMap["poolConfig"] = *poolConfig

			ORedisPoolMap.Store(poolName, newMap)

			return

		}
	} else {
		err = errors.New("no found poolName: " + poolName + " and no found config")
		return
	}

}

func RedisConnOGet(poolName string, poolConfig XmlRedis) (err error, rc RedisConn) {
	if rp, tok := ORedisPoolMap.Load(poolName); tok {
		err = nil

		redisMap := rp.(map[string]interface{})

		rc.PoolConfig = redisMap["poolConfig"].(XmlRedis)
		rc.Conn = (redisMap["redisPool"].(*redis.Pool)).Get()

		if rc.Conn == nil {
			err = errors.New("redis conn get from pool is nil, poolName: " + poolName)
		}

		return
	}

	err, newRp := InitRedisCacheConn(poolConfig)
	if err != nil {

		rc.Conn = nil

		return

	} else {

		rc.Conn = newRp.Get()

		newMap := make(map[string]interface{})
		newMap["redisPool"] = newRp
		newMap["poolConfig"] = poolConfig

		ORedisPoolMap.Store(poolName, newMap)

		return

	}

}

func (rc *RedisConn) Close() {
	if rc.Conn != nil {
		rc.Conn.Close()
	}

}

// this is used for interface !!!! no a normal struct
func (rc *RedisConn) SetExRegistered(key string, value interface{}, expireSecond int) (err error) {

	err = nil

	var buffer bytes.Buffer

	enc := gob.NewEncoder(&buffer)

	err0 := enc.Encode(&value)
	if err0 != nil {
		err = errors.New("Encode error key=" + key + " err=" + err0.Error())
		return
	}

	if expireSecond == -1 {
		_, err = rc.Conn.Do("SET", key, buffer.Bytes())
		if err != nil {
			err = errors.New("redis do set error" + key + " err=" + err.Error())

			return
		}

	} else {
		_, err = rc.Conn.Do("SET", key, buffer.Bytes(), "EX", expireSecond)
		if err != nil {
			err = errors.New("redis do set error" + key + " err=" + err.Error())

			return
		}

	}

	//newItem := &memcache.Item{Key: md5Key, Value: buffer.Bytes(), Expiration: int32(et)}

	//startTime := time.Now().UnixNano()

	// globalGoMetric.AddValue("cache_write", float32(time.Now().UnixNano()-startTime)/1000000)
	//l4g.Debug("MCSet ok key=%s md5key=%s value=%+v ", key, md5Key, value)

	return
}

func (rc *RedisConn) SetNxEx(key string, value interface{}, expireSecond int) (success bool, err error) {

	err = nil

	var buffer bytes.Buffer

	enc := gob.NewEncoder(&buffer)

	err0 := enc.Encode(value)
	if err0 != nil {
		err = errors.New("Encode error key=" + key + " err=" + err0.Error())
		return
	}

	ret, err := rc.Conn.Do("SET", key, buffer.Bytes(), "NX", "EX", expireSecond)
	if err != nil {
		err = errors.New("redis do set error" + key + " err=" + err.Error())

		return
	}
	//fmt.Println(ret)

	if ret != nil {
		var retString string
		retString, err = redis.String(ret, err)
		if err != nil {
			err = fmt.Errorf("redis string err: %s", err)
			return
		}
		if retString == "OK" {
			success = true
		}

		//fmt.Println(retString)
	}

	//newItem := &memcache.Item{Key: md5Key, Value: buffer.Bytes(), Expiration: int32(et)}

	//startTime := time.Now().UnixNano()

	// globalGoMetric.AddValue("cache_write", float32(time.Now().UnixNano()-startTime)/1000000)
	//l4g.Debug("MCSet ok key=%s md5key=%s value=%+v ", key, md5Key, value)

	return
}
func (rc *RedisConn) Do(cmd string, args ...interface{}) (result interface{}, exist bool, err error) {
	exist = false
	err = nil
	result, err = rc.Conn.Do(cmd, args...)
	if err != nil {

		err = errors.New("redis do get error cmd: " + cmd + " err=" + err.Error())

		//config.DefaultLogger.Error("redis get error err: %s, key: %s", err, WECHAT_ACCESSTOKEN_REDIS_KEY_PREFIX+wechatAccessToken)
		return
	}
	if result != nil {
		exist = true
	}

	return

}

func (rc *RedisConn) DoForBytes(cmd string, args ...interface{}) (result interface{}, exist bool, err error) {
	exist = false
	err = nil
	result, err = redis.Bytes(rc.Conn.Do(cmd, args...))
	if err != nil {

		err = errors.New("redis do get error cmd: " + cmd + " err=" + err.Error())

		//config.DefaultLogger.Error("redis get error err: %s, key: %s", err, WECHAT_ACCESSTOKEN_REDIS_KEY_PREFIX+wechatAccessToken)
		return
	}
	if result != nil {
		exist = true
	}

	return

}
func (rc *RedisConn) SetHash(hash string, key string, value interface{}) (err error) {

	err = nil

	var buffer bytes.Buffer

	enc := gob.NewEncoder(&buffer)

	err0 := enc.Encode(value)
	if err0 != nil {
		err = errors.New("Encode error key=" + key + " err=" + err0.Error())
		return
	}

	_, err = rc.Conn.Do("HSET", hash, key, buffer.Bytes())
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
func (rc *RedisConn) SetEx(key string, value interface{}, expireSecond int) (err error) {

	err = nil

	var buffer bytes.Buffer

	enc := gob.NewEncoder(&buffer)

	err0 := enc.Encode(value)
	if err0 != nil {
		err = errors.New("Encode error key=" + key + " err=" + err0.Error())
		return
	}

	_, err = rc.Conn.Do("SET", key, buffer.Bytes(), "EX", expireSecond)
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
func (rc *RedisConn) GetHash(hash string, key string, value interface{}) (exist bool, err error) {
	exist = false
	err = nil
	result, err0 := rc.Conn.Do("HGET", hash, key)
	if err0 != nil {

		err = errors.New("redis do get error key: " + key + " err=" + err0.Error())

		//config.DefaultLogger.Error("redis get error err: %s, key: %s", err, WECHAT_ACCESSTOKEN_REDIS_KEY_PREFIX+wechatAccessToken)
		return
	}

	if rc.PoolConfig.ReadOff == 1 {
		result = nil
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

func (rc *RedisConn) GetWithNil(key string, value interface{}) (exist bool, isNil bool, err error) {

	exist = false
	err = nil
	result, err0 := rc.Conn.Do("GET", key)
	if err0 != nil {

		err = errors.New("redis do get error key: " + key + " err=" + err0.Error())

		//config.DefaultLogger.Error("redis get error err: %s, key: %s", err, WECHAT_ACCESSTOKEN_REDIS_KEY_PREFIX+wechatAccessToken)
		return
	}

	if rc.PoolConfig.ReadOff == 1 {
		result = nil
	}

	if result != nil {
		exist = true

		buf := bytes.NewReader(result.([]byte))
		dec := gob.NewDecoder(buf)

		var tmpString string
		err0 = dec.Decode(&tmpString)
		if err0 == nil && tmpString == REDIS_NIL {

			isNil = true
			return
		}

		//dec.Reset()
		//buf.Seek(io.SeekStart, io.SeekStart)
		buf = bytes.NewReader(result.([]byte))
		dec = gob.NewDecoder(buf)

		err0 = dec.Decode(value)
		if err0 != nil {
			err = errors.New("Decode error key=" + key + " err=" + err0.Error())
			return
		}

	}

	return

}
func (rc *RedisConn) Get(key string, value interface{}) (exist bool, err error) {
	exist = false
	err = nil
	result, err0 := rc.Conn.Do("GET", key)
	if err0 != nil {

		err = errors.New("redis do get error key: " + key + " err=" + err0.Error())

		//config.DefaultLogger.Error("redis get error err: %s, key: %s", err, WECHAT_ACCESSTOKEN_REDIS_KEY_PREFIX+wechatAccessToken)
		return
	}

	if rc.PoolConfig.ReadOff == 1 {
		result = nil
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

func (rc *RedisConn) GetInterface(key string, value *interface{}) (exist bool, err error) {
	exist = false
	err = nil
	result, err0 := rc.Conn.Do("GET", key)
	if err0 != nil {

		err = errors.New("redis do get error key: " + key + " err=" + err.Error())

		//config.DefaultLogger.Error("redis get error err: %s, key: %s", err, WECHAT_ACCESSTOKEN_REDIS_KEY_PREFIX+wechatAccessToken)
		return
	}

	if rc.PoolConfig.ReadOff == 1 {
		result = nil
	}

	if result != nil {

		buf := bytes.NewBuffer(result.([]byte))
		dec := gob.NewDecoder(buf)

		err0 = dec.Decode(*value)
		if err0 != nil {
			err = errors.New("Decode error key=" + key + " err=" + err0.Error())
			return
		}
		exist = true

	}

	return

}
