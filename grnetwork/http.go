package grnetwork

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	IP_LOCAL       = 0
	IP_CDN_SRC     = 1
	IP_X_REAL      = 2
	IP_REMOTE_ADDR = 3
	IP_X_FORWARD   = 4
)
const CONNECT_TIMEOUT = time.Duration(10000) * time.Millisecond //1000 ms
func isLocalIp(ip string) bool {
	if strings.Contains(ip, "127.0.0.1") {
		return true
	}
	ips := strings.Split(ip, ".")

	ip1, err := strconv.Atoi(ips[1])
	if err != nil {
		return true
	}

	if "10" == ips[0] {
		return true
	}

	if "172" == ips[0] && ip1 > 15 && ip1 < 32 {
		return true
	}

	if "192" == ips[0] && 168 == ip1 {
		return true
	}

	return false
}

func GetUserIp(r *http.Request) (ip string, src int) {

	src = IP_LOCAL
	ip = r.Header.Get("Cdn-Src-Ip")
	if ip != "" {
		src = IP_CDN_SRC
		return
	}

	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		src = IP_X_REAL
		return
	}

	ip = r.RemoteAddr
	if index := strings.Index(ip, ":"); index != -1 {
		ip = ip[0:index]
	}

	if !isLocalIp(ip) {
		src = IP_REMOTE_ADDR
		return
	} else {
		ip = "127.0.0.1"
		return
	}

	/*
		//		ip = r.Header.Get("X-Forwarded-For")
	*/
	return
}

func timeoutDialler(t time.Time) func(net, addr string) (c net.Conn, err error) {

	return func(netw, addr string) (net.Conn, error) {
		c, err := net.DialTimeout(netw, addr, CONNECT_TIMEOUT)
		if err != nil {
			return nil, err
		}
		c.SetDeadline(t)
		c.SetReadDeadline(t)
		c.SetWriteDeadline(t)
		return c, nil
	}
}
func HttpGetExTimeout(url string, heads map[string]string, msec int) (rsp string, statusCode int, err error) {

	c := http.Client{
		Transport: &http.Transport{
			Dial:              timeoutDialler(time.Now().Add(time.Duration(msec) * time.Millisecond)),
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		},
	}
	rsp = ""
	err = nil

	var err0 error
	req, err0 := http.NewRequest("GET", url, nil)
	if err != nil {
		err = err0
		return
	}

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
	b := []byte{}
	b, err0 = ioutil.ReadAll(resp.Body)
	if err0 != nil {
		err = err0
		return
	}
	statusCode = resp.StatusCode
	if resp.StatusCode == 200 {
		rsp = string(b)
	} else {
		err = fmt.Errorf("resp Code is %d not 200", resp.StatusCode)
	}
	return
}
func HttpGetExTimeoutReTry(url string, heads map[string]string, msec int, tryTimes int) (rsp string, statusCode int, err error) {

	for i := 0; i < tryTimes; i++ {
		rsp, statusCode, err = HttpGetExTimeout(url, heads, msec)
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
func HttpPostExTimeoutReTry(url, body string, heads map[string]string, msec int, tryTimes int) (rsp string, statusCode int, err error) {

	for i := 0; i < tryTimes; i++ {
		rsp, statusCode, err = HttpPostExTimeout(url, body, heads, msec)
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
func HttpPostExTimeout(url, body string, heads map[string]string, msec int) (rsp string, statusCode int, err error) {

	c := http.Client{Transport: &http.Transport{
		Dial:              timeoutDialler(time.Now().Add(time.Duration(msec) * time.Millisecond)),
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	},
	}

	//golang那个设置的意思是设置了http连接的timeout时间之后，下次发起请求的时候，复用同一个连接，但是没有重新设置timeout.所以keepalive disable
	rsp = ""
	err = nil

	var err0 error
	req, err0 := http.NewRequest("POST", url, strings.NewReader(body))
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
	b := []byte{}
	b, err0 = ioutil.ReadAll(resp.Body)
	if err0 != nil {
		err = err0
		return
	}

	statusCode = resp.StatusCode

	if resp.StatusCode == 200 {
		rsp = string(b)
	} else {
		rsp = string(b)
		err = fmt.Errorf("resp Code is %d not 200", resp.StatusCode)
	}
	return
}

func HttpPostExTimeoutWithCert(url, body string, heads map[string]string, msec int, caCrtPath string, clientCrtPath string, clientKey string) (rsp string, statusCode int, err error) {
	pool := x509.NewCertPool()

	caCrt, err := ioutil.ReadFile(caCrtPath)
	if err != nil {
		return
	}
	cliCrt, err := tls.LoadX509KeyPair(clientCrtPath, clientKey)
	if err != nil {

		return
	}

	pool.AppendCertsFromPEM(caCrt)
	c := http.Client{Transport: &http.Transport{
		Dial: timeoutDialler(time.Now().Add(time.Duration(msec) * time.Millisecond)),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            pool,
			Certificates:       []tls.Certificate{cliCrt},
		},
		DisableKeepAlives: true,
	},
	}

	//golang那个设置的意思是设置了http连接的timeout时间之后，下次发起请求的时候，复用同一个连接，但是没有重新设置timeout.所以keepalive disable
	rsp = ""
	err = nil

	var err0 error
	req, err0 := http.NewRequest("POST", url, strings.NewReader(body))
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
	b := []byte{}
	b, err0 = ioutil.ReadAll(resp.Body)
	if err0 != nil {
		err = err0
		return
	}

	statusCode = resp.StatusCode

	if resp.StatusCode == 200 {
		rsp = string(b)
	} else {
		err = fmt.Errorf("resp Code is %d not 200", resp.StatusCode)
	}
	return
}
