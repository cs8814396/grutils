package grcommon

import (
	"testing"
)

func Test_CheckStringIfFormatCell(t *testing.T) {
	// go test -v -test.run Test_Sign
	testString := "13535133321"
	CheckStringIfFormatCell(testString)
}
