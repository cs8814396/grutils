package grthird

import (
	"testing"
)

func Test_Alterover(t *testing.T) {
	// go test -v

	source := "xxxxxxxx"
	user := "xxxxxxxxxxx"
	title := "xxxxx"
	content := "xxxx"

	err := AlertOverNotify(source, user, title, content)
	t.Log(err)

	return
	//MCENCRYPT_KEY_ACI = 7799
	//MCENCRYPT_KEY_PT  = 6712
	//MCENCRYPT_KEY_PSI = 6782

}
