package grfile

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

//parseConf
func LoadTomlFile(filename string, tomlStruct interface{}) (contents []byte, err error) {

	contents, err = ReadFile(filename)
	if err != nil {
		dir, dirErr := os.Getwd()
		err = fmt.Errorf("LoadConfig: Error: Could not open %q: %s.pwd: %s, pwderr: %s \n", filename, err, dir, dirErr)

		return
	}

	if _, err = toml.Decode(string(contents), tomlStruct); err != nil {

		err = fmt.Errorf("LoadConfig: Error: Could not parse XML configuration in %q: %s\n", filename, err)

		return
	}
	//fmt.Printf("Global Conf: \n%+v \n", xmlStruct)
	err = nil
	return

}
