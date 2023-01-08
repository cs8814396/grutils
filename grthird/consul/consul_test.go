package consul

import (
	//"encoding/json"
	"github.com/gdgrc/grutils/grapps/config"
	"github.com/gdgrc/grutils/grmath"

	"testing"
)

/*
ARMCHAIR_CONSUME_PAY_TYPE = (
    (0x01, '赠送余额'),
    (0x02, '余额'),
    (0x04, '微信'),
    (0x08, '支付宝'),
    (0x10, '工银e支付'),*/

func Test_Consul(t *testing.T) {
	// go test -v -test.run Test_WechatMp
	config.Init("/home/.../coding/massive_data_server/bin/massive_data_query_server/config.xml")
	c := NewConsul(config.GlobalConf.DefaultConsulConfig.Host, config.GlobalConf.DefaultConsulConfig.Port)
	c.IsMetaEncrypt = true
	var arb grmath.Aes256Rc5Base64Cryptor

	arb.AesKey = []byte(config.GlobalConf.Aes.AesKey)
	arb.AesIv = []byte(config.GlobalConf.Aes.AesIV)
	arb.Rc4Key = []byte(config.GlobalConf.Rc4.Rc4Key)

	encryptData, err := arb.Encrypt([]byte(`{"username": "admin", "password": "xxx"}`))
	if err != nil {
		t.Fatalf("encrypt err: %s", err.Error())
	}

	dd, err := arb.Decrypt(string(encryptData))
	if err != nil {
		t.Fatalf("decrypt err: %s", err.Error())
	}
	t.Logf("encrypt data: %s dd: %s", string(encryptData), string(dd))

	c.MetaCryptor = &arb

	/*
		c.IsMetaEncrypt = true

		var arb grmath.Aes256Rc5Base64Cryptor
		totalKey := grmath.Md5("key for aes333444444")

		rc4Key := grmath.Md5("key for rc455555")

		aesKey := totalKey[:16]
		aesIv := totalKey[16:32]

		t.Log(aesKey, aesIv, rc4Key)
		arb.AesKey = []byte(aesKey)
		arb.AesIv = []byte(aesIv)
		arb.Rc4Key = []byte(rc4Key)

		c.MetaCryptor = &arb*/

	ch, err := c.GetHealthServices("bdb_0")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ch, err)
	return
	//MCENCRYPT_KEY_ACI = 7799
	//MCENCRYPT_KEY_PT  = 6712
	//MCENCRYPT_KEY_PSI = 6782

}
