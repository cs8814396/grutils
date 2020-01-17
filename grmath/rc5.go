package grmath

import (
	"bytes"

	//"encoding/base64"
	//"encoding/hex"
	"unicode/utf8"
)

//变种rc4
func Rc4Encode(data []byte, rc4Key []byte) (res []byte) {
	box := make([]int, 256)
	randKey := make([]int, 256)
	key_length := len(rc4Key)
	res = make([]byte, len(data))

	var i int = 0
	var j int = 0
	var a int = 0
	var tmp int
	for i = 0; i < 255; i++ {
		randKey[i] = int(rc4Key[i%key_length])
	}
	for i = 0; i < 256; i++ {
		box[i] = i
	}
	for i = 0; i < 255; i++ {
		j = 0
		j = int(j+box[i]+randKey[i]) % 256
		tmp = box[i]
		box[i] = box[j]
		box[j] = tmp
	}
	for i = 0; i < len(data); i++ {
		j = 0
		a = 0
		a = (a + 1) % 256
		j = (j + box[a]) % 256
		tmp = box[a]
		box[a] = box[j]
		box[j] = tmp
		res[i] = (data[i] ^ byte(box[(box[a]+box[j])%256]))
	}

	return res
}

func Rc4Decode(data []byte, rc4Key []byte) (res []byte) {
	return Rc4Encode(data, rc4Key)
}

//unicode字节数组转utf8
func EncodeUnicodeToUtf8(data []byte) []byte {
	s := make([][]byte, len(data))
	for i := 0; i < len(data); i++ {
		b := make([]byte, utf8.UTFMax)
		n := utf8.EncodeRune(b, rune(data[i]))
		s[i] = b[:n]
	}

	return bytes.Join(s, []byte(""))
}

//utf8字节数组转unicode
func DecodeUtf8ToUnicode(data []byte) (res []byte) {
	res = make([]byte, 0)
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		data = data[size:]
		res = append(res, byte(r))
	}
	return res
}

/*
//加密字节数组接口
func EncodeData(data []byte, aesKey []byte, aesIV []byte, rc4Key []byte) (res []byte, err error) {
	var tmpData []byte

	if tmpData, err = AesCBCZeroPaddingEncrypt(data, aesKey, aesIV); err != nil {
		return []byte(""), err
	}
	tmpDataStr := hex.EncodeToString(tmpData)
	tmpData = []byte(tmpDataStr)
	tmpData = Rc4Encode(tmpData, rc4Key)
	tmpData = EncodeUnicodeToUtf8(tmpData)
	tmpStr := base64.StdEncoding.EncodeToString(tmpData)
	res = []byte(tmpStr)
	return res, nil
}

//解密字节数组接口
func DecodeData(data []byte, aesKey []byte, aesIV []byte, rc4Key []byte) (res []byte, err error) {
	var tmpData []byte

	if tmpData, err = base64.StdEncoding.DecodeString(string(data)); err != nil {
		return []byte(""), err
	}
	tmpData = DecodeUtf8ToUnicode(tmpData)
	tmpData = Rc4Decode(tmpData, rc4Key)
	tmpData, _ = hex.DecodeString(string(tmpData))
	if res, err = AesCBCZeroPaddingDecrypt(tmpData, aesKey, aesIV); err != nil {
		return []byte(""), err
	}
	return res, nil
}
*/
