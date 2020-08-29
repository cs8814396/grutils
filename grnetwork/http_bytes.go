package grnetwork

import (
	"bytes"
	"crypto/tls"
	//"crypto/x509"
	"fmt"
	"io/ioutil"
	//"net"
	"net/http"
	//"strconv"
	//"strings"
	"time"
)

func HttpPostRetry(url string, body []byte, heads map[string]string, msec int, tryTimes int) (rsp []byte, statusCode int, err error) {

	for i := 0; i < tryTimes; i++ {
		rsp, statusCode, err = HttpPost(url, body, heads, msec)
		if err != nil {
			continue
		} else {
			break
		}

	}
	if err != nil {
		return
	}
	return
}
func HttpPost(url string, body []byte, heads map[string]string, msec int) (rsp []byte, statusCode int, err error) {

	c := http.Client{Transport: &http.Transport{
		Dial:              timeoutDialler(time.Now().Add(time.Duration(msec) * time.Millisecond)),
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	},
	}

	//golang那个设置的意思是设置了http连接的timeout时间之后，下次发起请求的时候，复用同一个连接，但是没有重新设置timeout.所以keepalive disable
	rsp = nil
	err = nil

	var err0 error
	req, err0 := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		err = err0
		return
	}
	// req.Header.Set("Content-Type", contentType)

	//Set Header
	if heads != nil {
		for k, v := range heads {
			req.Header.Add(k, v)
		}
	}

	resp, err0 := c.Do(req)
	if err0 != nil {
		err = err0
		return
	}

	defer resp.Body.Close()

	rsp, err0 = ioutil.ReadAll(resp.Body)
	if err0 != nil {
		err = err0
		return
	}

	statusCode = resp.StatusCode

	if resp.StatusCode == 200 {

	} else {

		err = fmt.Errorf("resp Code is %d not 200", resp.StatusCode)
	}
	return
}
