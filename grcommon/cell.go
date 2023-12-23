package grcommon

import (

	//"unicode"
	"regexp"
)

var cellCodeList = []string{"13", "14", "15", "16", "17", "18", "19"}

func CheckStringIfFormatCell(cell string) bool {

	result, err := regexp.MatchString("[\\d]{11}", cell)
	if err != nil {
		return false
	}

	find := false
	for _, c := range cellCodeList {
		if c == cell[0:2] {
			find = true
			break
		}
	}
	if !find {
		return false
	}

	return result

}
