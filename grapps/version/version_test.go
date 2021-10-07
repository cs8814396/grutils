package version

import (
	//"encoding/json"

	"github.com/gdgrc/grutils/grapps/config"

	//"grutils/grmath"
	"testing"
)

func Test_Order(t *testing.T) {
	// go test -v -test.run Test_Order

	config.Init("../../config.xml")

	version, exist, err := GetVersionFromDBorCache("grtaskwall", 10001, true)
	t.Log(version, exist, err)

	//appsecret := "70c6ce9f0a7de960673be6009e8d9caf"
	//ConsumeOrderCreate

	//rspJson, _ := json.Marshal(rsp)
	//t.Logf("%+v, %s", rspt, err)

	return
	//MCENCRYPT_KEY_ACI = 7799
	//MCENCRYPT_KEY_PT  = 6712
	//MCENCRYPT_KEY_PSI = 6782

}
