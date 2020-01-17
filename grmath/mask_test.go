package grmath

import (
	"testing"
)

func Test_MaskWithPartialMd5(t *testing.T) {
	// go test -v -test.run Test_Sign
	testString := "442111199001240241"
	markString, md5, err := MaskWithPartialMd5(testString, 9, 14)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(markString, md5)

	unMark, err := UnMaskWithPartialMd5(markString, 9, 14, md5, 1)

	if err != nil {
		t.Fatal(err)
	}

	if unMark != testString {
		t.Fatal("mark and unmark is not equal", unMark, testString)
	}

	testString = "12739135470"
	markString, md5, err = MaskWithWholeMd5(testString, 7, 10)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(markString, md5)
	unMark, err = UnMaskWithWholeMd5(markString, 7, 10, md5, 1)

	if err != nil {
		t.Fatal(err)
	}

	if unMark != testString {
		t.Fatal("mark and unmark is not equal", unMark, testString)
	}

}
