package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	initConfig("../etc/config.debug.xml", "../etc/data_fetcher_conf.toml")
	code := m.Run()
	os.Exit(code)
}
