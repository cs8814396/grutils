package consul

import (
	"encoding/json"
	"fmt"
	"grutils/grmath"
	"grutils/grnetwork"
)

type Consul struct {
	ConsulHost    string
	ConsulPort    int
	IsMetaEncrypt bool

	ConsulServiceUrl string

	//=====

	MetaCryptor grmath.Cryptor
}

type ConsulHealth struct {
	Node    ConsulHealthNode    `json:"Node"`
	Service ConsulHealthService `json:"Service"`
}
type ConsulHealthNode struct {
	ID      string `json:"ID"`
	Address string `json:"Address"`
	Port    int    `json:"Port"`
}
type ConsulHealthService struct {
	Address string            `json:"Address"`
	Port    int               `json:"Port"`
	Meta    map[string]string `json:"Meta"`
}

func NewConsul(consulHost string, consulPort int) *Consul {
	c := Consul{}
	c.ConsulHost = consulHost
	c.ConsulPort = consulPort

	c.ConsulServiceUrl = fmt.Sprintf("http://%s:%d", c.ConsulHost, c.ConsulPort)

	return &c
}

func (this *Consul) GetHealthServices(serviceName string) (chSlice []ConsulHealth, err error) {

	url := fmt.Sprintf("%s/v1/health/service/%s?passing=true", this.ConsulServiceUrl, serviceName)

	heads := make(map[string]string)

	msec := 5000

	tryTimes := 2

	rsp, _, err := grnetwork.HttpGetExTimeoutReTry(url, heads, msec, tryTimes)

	err = json.Unmarshal([]byte(rsp), &chSlice)
	if err != nil {
		err = fmt.Errorf("json umarshal fail: %s,url: %s, data: %s", err, url, rsp)
		return
	}

	if this.IsMetaEncrypt {
		for _, ch := range chSlice {
			if metaConfig, ok := ch.Service.Meta["config"]; ok {
				var decryptBytes []byte

				decryptBytes, err = this.MetaCryptor.Decrypt(metaConfig)

				if err != nil {
					err = fmt.Errorf("err: %s, encrypt data: %s", err, metaConfig)
					return
				}

				err = json.Unmarshal(decryptBytes, &ch.Service.Meta)
				if err != nil {
					err = fmt.Errorf("err: %s, decrypt data unmarshal fail", err)
					return
				}

			}
		}
	}

	return

}
