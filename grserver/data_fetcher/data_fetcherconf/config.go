package data_fetcherconf

import (
	"fmt"
	"github.com/gdgrc/grutils/grfile"
)

type XmlServer struct {
	Bindaddr string `xml:"bindaddr"`
	Host     string `xml:"host"`
}

type xmlConfig struct {
	//XMLName    xml.Name  `xml:"config"`
	Server XmlServer `xml:"server"`
}

var GlobalConf xmlConfig

func Init(filename string, xmlFilename string) bool {
	_, err := grfile.LoadXmlConfig(filename, &GlobalConf)
	if err != nil {
		fmt.Println(err)
		return false
	}

	_, err = grfile.LoadTomlFile(xmlFilename, &GlobalDataFetcherConf)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
