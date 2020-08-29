package data_fetcherconf

import (
	"github.com/gdgrc/grutils/grfile"
	"testing"
)

func TestParseConf(t *testing.T) {

	_, err := grfile.LoadTomlFile("../etc/data_fetcher_conf.toml", &GlobalDataFetcherConf)
	if err != nil {
		t.Fatal(err)
	}

	config_data_name := "config_data_name"

	_, ok := GlobalDataFetcherConf.Querys[config_data_name]
	if !ok {
		t.Fatal("why not exist data: ", config_data_name)
	}
	t.Log("Pass")
}
