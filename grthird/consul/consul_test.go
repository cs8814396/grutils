package consul

import (
	//"encoding/json"
	"grutils/grapps/config"

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
	config.Init("../../config.xml")
	c := NewConsul(config.GlobalConf.DefaultConsulConfig.Host, config.GlobalConf.DefaultConsulConfig.Port)

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
