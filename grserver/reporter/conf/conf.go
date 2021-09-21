package conf

import (
	"fmt"
	"github.com/gdgrc/grutils/grfile"
)

type xmlConfig struct {
	//XMLName    xml.Name  `xml:"config"`
	//Server XmlServer `xml:"server"`
}

type Conf struct {
	//TODO: define you config here
	LogDir string   `toml:"log_dir"`
	Events []string `toml:"events"`
}

var GlobalReporterConf Conf

var GlobalConf xmlConfig

func Init(filename string, xmlFilename string) bool {

	_, err := grfile.LoadTomlFile(xmlFilename, &GlobalReporterConf)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
