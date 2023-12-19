package grcommon

import "regexp"

func CheckMail(mail string) bool {
	if len(mail) <= 7 {
		return false
	}
	result, err := regexp.MatchString("^.+\\@(\\[?)[a-zA-Z0-9\\-\\.]+\\.([a-zA-Z]{2,3}|[0-9]{1,3})(\\]?)$", mail)
	if err != nil {
		return false
	}
	return result
}
