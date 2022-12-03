package grmonitor

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func StartPProfWithNewHttp(host string) {
	go func() {
		err := http.ListenAndServe(host, nil)
		if err != nil {
			log.Fatalf("StartPProfWithNewHttp err: " + err.Error())
		}
	}()
}
