package grmath

import (
	"encoding/base64"
	//"encoding/json"
	//"bytes"
	"bytes"
	"crypto/rc4"

	//"grutils/grcommon"
	//"net/url"
	//"strings"
	//"encoding/hex"
	"testing"
)

/*
ARMCHAIR_CONSUME_PAY_TYPE = (
    (0x01, '赠送余额'),
    (0x02, '余额'),
    (0x04, '微信'),
    (0x08, '支付宝'),
    (0x10, '工银e支付'),*/

type ElifeProxyCallbackrReq struct {
	Charge       int    `json:"charge"`
	SubsId       string `json:"subs_id"`
	UserId       string `json:"user_id"`
	UserPhone    string `json:"user_phone"`
	UserDeviceId string `json:"user_deviceID"`
	OrderId      string `json:"order_id"`
}

func Test_AES_RC4_BASE64(t *testing.T) {
	// go test -v

	key := []byte("64sfe99aec1b0d57")
	iv := []byte("c709e1ccb45362hg")

	testData := `{"charge":1,"orderid":"C20180512015536035127","subs_name":"aaa"}`

	encrytData, err := AesCBCZeroPaddingEncrypt([]byte(testData), key, iv)
	if err != nil {
		t.Fatal(err)
	}

	rc4key := []byte("c709e1ccb4gggggg64sfe99aec1b0d57")

	rc4obj, _ := rc4.NewCipher(rc4key)

	rc4Result := make([]byte, len(encrytData))

	rc4obj.XORKeyStream(rc4Result, encrytData)

	encryptBase64 := base64.StdEncoding.EncodeToString(rc4Result)

	//encryptBase64 = "je1oWg5A9v0vANBnJbTrNahW9xi4mirLfLSt6CTymdQ3TlVfodOtX4nYKtPkkByycsSd44mli75fS2ZW/9w/sA=="

	/*
		decryptBase64, _ := base64.StdEncoding.DecodeString(encryptBase64)

		decryptRc4 := make([]byte, len(decryptBase64))

		rc4obj, _ = rc4.NewCipher(rc4key)

		rc4obj.XORKeyStream(decryptRc4, decryptBase64)
		originData, err := AesCBCZeroPaddingDecrypt(decryptRc4, key, iv)

		t.Log(string(originData))*/

	decryptBase64, _ := base64.StdEncoding.DecodeString(encryptBase64)

	if !bytes.Equal(decryptBase64, rc4Result) {
		t.Fatal("not equal ", decryptBase64, rc4Result)
	}

	decryptRc4 := make([]byte, len(decryptBase64))

	rc4obj, _ = rc4.NewCipher(rc4key)

	rc4obj.XORKeyStream(decryptRc4, decryptBase64)

	if !bytes.Equal(decryptRc4, encrytData) {
		t.Fatal("not equal")
	}

	originData, err := AesCBCZeroPaddingDecrypt(decryptRc4, key, iv)
	if string(originData) != testData {
		t.Fatal("not equal", string(originData), testData)
	}

}

func Test_AES(t *testing.T) {
	// go test -v

	/*
		flatMap := make(map[string]interface{})
		err := json.Unmarshal([]byte(data), &json_data)

		t.Log(json_data)
		if err == nil {
			// InterfaceToString(json_data["errorres"]["detailinfo"]["FlightInfo"]["PSA_FlightInfo"]["FlightDetail"]["PSA_FlightDetail"])
			t.Log(Md5("123456"))
			getKeyFromMap(JSONDATA_KEY_TRANSFORM_MAP, json_data, flatMap)

			t.Log(flatMap)

		} else {
			t.Error("err: %s", err)
		}
	*/

	//orderId := "qR1NqSfSH3q2zycjra4GFA=="

	//ciphertext, err := base64.StdEncoding.DecodeString("kKXRcSiMtgmo53L1Nj0f+TKRi51uvqJwMfC8Suvtp4I8Wqs2Hqt3mfMSaEuUcaV0qX8VzlqJpaMM6k0dQdQz8tVwSfUfL2MPv4jaXDhRPVI=")
	//t.Log(err)

	key := []byte("64sfe99aec1b0d57")
	iv := []byte("c709e1ccb45362hg")

	/*

		c, err := AesCBCZeroPaddingEncrypt([]byte("10001360"), key, iv)

		d := hex.EncodeToString(c)

		t.Log(string(c), d, err)

		c, err = AesCBCPKCS5Encrypt([]byte("10001360"), iv, key)

		d = hex.EncodeToString(c)

		t.Log(string(c), d, err)

		c, err = AesCBCPKCS5Encrypt([]byte(`{"charge":1,"orderid":"C20180512015536035127","subs_name":"aaa"}`), key, iv)

		d = base64.StdEncoding.EncodeToString(c)

		t.Log(d)*/

	testData := `{"charge":1,"orderid":"C20180512015536035127","subs_name":"aaa"}`

	encrytData, err := AesCBCZeroPaddingEncrypt([]byte(testData), key, iv)
	if err != nil {
		t.Fatal(err)
	}

	originData, err := AesCBCZeroPaddingDecrypt(encrytData, key, iv)
	if string(originData) != testData {
		t.Fatal("not equal", string(originData), testData)
	}

}
