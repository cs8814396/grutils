package grcommon

import (
	"unicode"
)

func CheckName(name string) bool {
	//中文加这个点
	//tmpName := strings.ReplaceAll(name, "·", "")

	result := true
	for _, r := range name {
		if !unicode.Is(unicode.Han, r) && string(r) != "·" {
			result = false
		}
	}

	/*result, err := regexp.MatchString("^[\u4e00-\u9fa5]$", tmpName)
	if err != nil {
		return false
	}*/

	return result

}
