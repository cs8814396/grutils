package main

import (
	econfig "device_filter/reporter/conf"
	"flag"
	"fmt"
	"github.com/gdgrc/grutils/grapps/config"
	"github.com/gdgrc/grutils/grframework/fasthttp"

	"device_filter/reporter/handler"
	"os"
	"path"
	"syscall"
	"time"
)

var configFile = flag.String("c", "../conf/config.xml", "公共配置文件地址（绝对路径或者bin目录为基准的相对路径）")
var specificConfigFile = flag.String("sc", "../conf/conf_dev.toml", "配置文件地址（绝对路径或者bin目录为基准的相对路径）")

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
	return econfig.Init(configFilePath, specificConfigFilePath) && config.Init(configFilePath) && handler.Init()
}

func main() {
	if !Init() {
		time.Sleep(1e9)
		return
	}
	fasthttp.Register("/commReport", handler.CommReport)
	fasthttp.Register("/post_web_env_v2", handler.Report)
	fasthttp.Register("/getWebEnvResults_v2", handler.GetResult)
	fasthttp.ListenAndBlock(config.GlobalConf.Server.BindAddr)
}
