package grcommon

import (
	"regexp"
	"strings"
)

func CheckName(name string) bool {
	//中文加这个点
	tmpName := strings.ReplaceAll(name, "·", "")

	result, err := regexp.MatchString("^[\u4e00-\u9fa5]$", tmpName)
	if err != nil {
		return false
	}

	return result

}
