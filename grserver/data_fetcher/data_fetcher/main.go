package main

import (
	"flag"
	"fmt"
	"github.com/gdgrc/grutils/grapps/config"
	"github.com/gdgrc/grutils/grframework/fasthttp"

	econfig "data_fetcher/data_fetcherconf"

	"os"
	"path"
	"syscall"
	"time"
)

var configFile = flag.String("c", "", "公共配置文件地址（绝对路径或者bin目录为基准的相对路径）")
var specificConfigFile = flag.String("sc", "", "配置文件地址（绝对路径或者bin目录为基准的相对路径）")

var displayHelp = flag.Bool("help", false, "显示此帮助信息")

func Init() bool {
	flag.Parse()
	fmt.Printf("help:[ %t ] c:[ %s ]\n", *displayHelp, *configFile)
	if *displayHelp || *configFile == "" {
		flag.PrintDefaults()
		return false
	}
	syscall.Umask(0)
	os.Chdir(path.Dir(os.Args[0]))

	return initConfig(*configFile, *specificConfigFile)
}

func initConfig(configFilePath string, specificConfigFilePath string) bool {
	return econfig.Init(configFilePath, specificConfigFilePath) && config.Init(configFilePath)
}

func main() {
	if !Init() {
		time.Sleep(1e9)
		return
	}
	fasthttp.Register("/fetch_data", FetchData)
	fasthttp.Register("/insert_data", InsertData)
	fasthttp.ListenAndBlock(config.GlobalConf.Server.BindAddr)
}
