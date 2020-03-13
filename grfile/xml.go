package grfile

import (
	"encoding/xml"
	"fmt"
	//"grutils/grdatabase"
	//""
)

type XmlMemcache struct {
	Addr      string `xml:"addr"`
	RWTimeout int    `xml:"rwtimeout"`
}

func LoadXmlConfigWithContents(contents []byte, xmlStruct interface{}) (err error) {

	if err = xml.Unmarshal(contents, xmlStruct); err != nil {

		err = fmt.Errorf("LoadConfig: Error: Could not parse XML configuration in %s: %s\n", string(contents), err)

		return
	}
	//fmt.Printf("Global Conf: \n%+v \n", xmlStruct)

	return nil

}
func LoadXmlConfig(filename string, xmlStruct interface{}) (contents []byte, err error) {
	/*
			defer func() {
				if r := recover(); r != nil {
					if _, ok := r.(runtime.Error); ok {
						panic(r)
					}
					err = r.(error)
				}
			}()

		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			return &InvalidUnmarshalError{reflect.TypeOf(v)}
		}*/

	contents, err = ReadFile(filename)
	if err != nil {
		err = fmt.Errorf("LoadConfig: Error: Could not open %q: %s \n", filename, err)

		return
	}

	if err = xml.Unmarshal(contents, xmlStruct); err != nil {

		err = fmt.Errorf("LoadConfig: Error: Could not parse XML configuration in %q: %s\n", filename, err)

		return
	}
	//fmt.Printf("Global Conf: \n%+v \n", xmlStruct)
	err = nil
	return

}

type XmlServer struct {
	Bindaddr string `xml:"bindaddr"`
	Host     string `xml:"host"`
}

/*
type xmlConfig struct {
	//XMLName    xml.Name  `xml:"config"`
	Server     XmlServer              `xml:"server"`
	DefaultLog grfile.XmlLogger       `xml:"defaultlog"`
	DataAdmin  grdatabase.XmlMysqlDsn `xml:"dataadmin"`

	RedisPool grcache.XmlRedis     `xml:"redispool"`
	AlertOver grthird.XmlAlertOver `xml:"alertover"`
	Wechat    grthird.XmlWechat    `xml:"wechat"`
}
*/
