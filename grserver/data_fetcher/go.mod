module github.com/gdgrc/grutils/grserver/data_fetcher

go 1.13

require (
	github.com/gdgrc/grutils v0.0.0-20200829134842-c4bed968fc2e
	github.com/gin-gonic/gin v1.5.0
	github.com/valyala/fasthttp v1.9.0
)

replace github.com/gdgrc/grutils => ../../../grutils

replace github.com/gdgrc/grutils/grserver/data_fetcher => ../../../grutils/grserver/data_fetcher
