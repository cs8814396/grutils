package grmath

import (
	"errors"
	"fmt"
	"testing"
)

var AES_KEY []byte = []byte("cf4d96a6755d4cc2")
var AES_IV []byte = []byte("93d8c38908caagga")
var RC4_KEY []byte = []byte("26689cf23f9923f74ca30aa9a75css11")

func TestCrypt(test *testing.T) {
	var err error
	if err = testCryptString("{\"username\": \"hello\", \"password\": \"world111111\"}"); err != nil {
		test.Fatal(err)
		return
	}
	if err = testCryptString("{\"username\": \"hello\", \"password\": \"world1111111\"}"); err != nil {
		test.Fatal(err)
		return
	}
}

func testCryptString(str string) (err error) {
	//go test -v -test.run Test_Sign
	res, err := Aes256Rc5Base64Encrypt([]byte(str), AES_KEY, AES_IV, RC4_KEY)
	if err != nil {
		return err
	}
	back, err := Aes256Rc4Base64Decrypt(res, AES_KEY, AES_IV, RC4_KEY)
	if err != nil {
		return err
	}

	if string(back) == str {
		fmt.Println("decode back equal")
	} else {
		return errors.New("decode back not equal")
	}

	return nil
}
