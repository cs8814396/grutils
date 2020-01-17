package grmath

import (
	"errors"
	"fmt"
	"strconv"
)

func IntegerMCEncrypt(code int, encrypt_key int, check_num int) (out int, err error) {
	err = nil
	out = 0

	cmbNumStr := fmt.Sprintf("%d%d", code, encrypt_key)

	cmbNum, err := strconv.ParseInt(cmbNumStr, 10, 32)
	if err != nil {
		return
	}
	cmbNum32 := cmbNum

	checkNum := 1
	for cmbNum32 > 0 {
		tn := int(cmbNum32 % 100)
		checkNum = checkNum * tn
		cmbNum32 = cmbNum32 / 100

	}

	formatString := fmt.Sprintf("%%d%%0%dd", check_num)

	formatMod := Pow(10, check_num)

	outStr := fmt.Sprintf(formatString, code, checkNum%(formatMod))

	tout, err := strconv.ParseInt(outStr, 10, 32)
	if err != nil {
		return
	}
	out = int(tout)
	// cmdNumStrLen := len(cmbNumStr)

	return
}

func IntegerMCDecrypt(encrypt_code int, encrypt_key int, check_num int) (out int, err error) {
	err = nil
	out = 0

	formatMod := Pow(10, check_num)
	code := encrypt_code / formatMod

	genCode, err := IntegerMCEncrypt(code, encrypt_key, check_num)
	if err != nil {
		return
	}

	if encrypt_code == genCode {
		out = code
		return
	} else {
		err = errors.New(fmt.Sprintf("IntegerMCDecrypt check is not match! gencode: %d", genCode))
		return
	}

	return
}
