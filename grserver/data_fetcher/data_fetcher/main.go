package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"syscall"
	"time"

	"github.com/gdgrc/grutils/grapps/config"
	"github.com/gdgrc/grutils/grapps/config/log"
	"github.com/gdgrc/grutils/grframework"
	"github.com/gdgrc/grutils/grframework/fasthttp"
	econfig "github.com/gdgrc/grutils/grserver/data_fetcher/data_fetcherconf"

	"encoding/base64"
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

func AuthMiddleware(c *grframework.Context) (err *grframework.Error) {
	if econfig.GlobalDataFetcherConf.Auth.ClientId != "" {

		auth := c.Headers["Authorization"]
		target := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(econfig.GlobalDataFetcherConf.Auth.ClientId+":"+ econfig.GlobalDataFetcherConf.Auth.ClientSecret)))
		log.Debug("headers: ", c.Headers, "auth: ", auth, " target: ", target)

		if auth != target {
			err = grframework.NewError(-1, "auth fail")
			return
		}
	}
	return
}
func main() {
	if !Init() {
		time.Sleep(1e9)
		return
	}
	fasthttp.Register("/fetch_data", FetchData, AuthMiddleware)
	fasthttp.Register("/insert_data", InsertData, AuthMiddleware)
	fasthttp.ListenAndBlock(config.GlobalConf.Server.BindAddr)
}
