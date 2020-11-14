package grcommon

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"sort"

	"golang.org/x/text/encoding/charmap"
)

func GetStringMapMd5(resultMap map[string]string) (mapMd5 string, err error) {
	uniqString := ""
	sorted_keys := make([]string, 0)
	for k, _ := range resultMap {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)

	for _, k := range sorted_keys {
		v := resultMap[k]

		uniqString += v

	}
	uniqStringBytes := []byte(uniqString)
	md5Hash := md5.Sum(uniqStringBytes)
	mapMd5 = fmt.Sprintf("%x", md5Hash)

	return
}
func GetMapMd5(resultMap map[string]interface{}) (mapMd5 string, err error) {
	uniqString := ""
	sorted_keys := make([]string, 0)
	for k, _ := range resultMap {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)

	for _, k := range sorted_keys {
		v := resultMap[k]
		switch vType := v.(type) {
		case string:
			uniqString += v.(string)
		case int:
			uniqString += fmt.Sprintf("%d", v.(int))
		case int64:
			uniqString += fmt.Sprintf("%d", v.(int64))
		case map[string]interface{}, []interface{}:
			var tmpBytes []byte
			tmpBytes, err = json.Marshal(vType)
			if err != nil {
				return
			}
			uniqString += fmt.Sprintf("%x", md5.Sum(tmpBytes))
		default:
			err = fmt.Errorf("getUniqCode Can't get type: %s", vType)
			return

		}
	}
	uniqStringBytes := []byte(uniqString)
	md5Hash := md5.Sum(uniqStringBytes)
	mapMd5 = fmt.Sprintf("%x", md5Hash)

	return
}
func SubStrOnBytes(str string, start, count int) string {

	ByteArr := []byte(str)
	len := len(ByteArr)

	min := start + count
	if min > len {
		min = len
	}
	return string(ByteArr[start:min])
}
func JsonMarshalWithoutEscape(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func ISO88598ToUtf8(s []byte) []byte {
	e := charmap.ISO8859_8.NewDecoder()
	es, _, err := transform.Bytes(e, s)
	if err != nil {
		return s
	}
	return es
}

func LatinToUTF8(iso8859_1_buf []byte) string {
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)
}
func GBKToUTF8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func UTF8ToGBK(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
func IfExistInStringSlice(stringList []string, ele string) (index int, ok bool) {
	ok = false
	for i, v := range stringList {
		if v == ele {
			ok = true
			index = i

		}
	}
	return
}
func SliceToMySqlString(name string, s interface{}) (str string) {
	if s == nil {
		str = " 1=2"
		return
	}

	var err error

	var opt string
	switch vType := s.(type) {

	case []int64, []string, []int:
		var b []byte
		b, err = json.Marshal(s)
		str = string(b)

		opt = "in"

		length := len(str)

		if length > 2 {

			str = fmt.Sprintf("%s %s (%s)", name, opt, str[1:length-1])
		} else {
			str = " 1=2"
		}
	default:
		err = fmt.Errorf("type not good %s", vType)
	}

	if err != nil {
		panic(fmt.Sprintf("SliceToSqlString type is not good %s,err: %s", s, err))
	}

	return
}
