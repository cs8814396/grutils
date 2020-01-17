package grmath

import (
	"encoding/json"
	"github.com/ugorji/go/codec"
	"reflect"
	"strings"
	"time"
)

//TimeChange 转化为标准模式输出
func TimeChange(timestamp int64) string {
	const layout = "2006-01-02 15:04:05"
	timeChange := time.Unix(timestamp, 0)
	return timeChange.Format(layout)
}

//Marshal 转化为json返回
func Marshal(res interface{}) (bytes []byte, err error) {
	bytes, err = json.Marshal(res)
	if err != nil {

		return
	}
	return
}

var (
	mapStrIntfTyp = reflect.TypeOf(map[string]interface{}(nil))
)

func UnSerialized(buffer []byte, v interface{}) (err error) {

	var msgpack codec.MsgpackHandle
	msgpack.MapType = mapStrIntfTyp
	dec := codec.NewDecoderBytes(buffer, &msgpack)
	err = dec.Decode(&v)
	return err
}

func Serialized(buffer *[]byte, v interface{}) (err error) {

	var msgpack codec.MsgpackHandle
	msgpack.MapType = mapStrIntfTyp
	enc := codec.NewEncoderBytes(buffer, &msgpack)
	err = enc.Encode(v)
	return err
}

// InArray 查找是否在数组里面
func InArray(result string, array []string) bool {
	for _, values := range array {
		if strings.Contains(result, values) {
			return true
		}
	}
	return false
}
