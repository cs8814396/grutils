package grcommon

import (

	//"unicode"
	"regexp"
)

func CheckStringIfFormatCell(cell string) bool {

	result, err := regexp.MatchString("[\\d]{11}", cell)
	if err != nil {
		return false
	}

	return result

}
