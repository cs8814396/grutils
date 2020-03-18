package grthird

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gdgrc/grutils/grcache"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

const (
	WECHAT_SESSIONKEY_REDIS_KEY_PREFIX          = "WechatSessionKey:"
	WECHAT_ACCESSTOKEN_REDIS_KEY_PREFIX         = "WechatAccessToken:"
	WECHAT_REFRESHTOKEN_REDIS_KEY_PREFIX        = "WechatRefreshToken:"
	WECHAT_COMPANY_ACCESSTOKEN_REDIS_KEY_PREFIX = "WechatCompanyAccessToken:"
	WECHAT_JSTICKET_REDIS_KEY_PREFIX            = "WechatJsTicket:"
)

type XmlWechat struct {
	Appid           string `xml:"appid"`
	AppSecret       string `xml:"appsecret"`
	MChnId          string `xml:"mchn_id"`
	NotifyUrl       string `xml:"notify_url"`
	ReturnUrl       string `xml:"return_url"`
	RefundNotifyUrl string `xml:"refund_notify_url"`

	CaCertPath string `xml:"ca_cert_path"`

	ClientCertPath string `xml:"client_cert_path"`
	ClientKeyPath  string `xml:"client_key_path"`

	Key      string `xml:"key"`
	OauthUrl string `xml:"oauth_url"`
}

type WechatOpenidRsp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	ReFreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionid"`

	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
type WechatUserInfoRsp struct {
	OpenId     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgUrl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege,omitempty"`
	UnionId    string   `json:"unionid"`
	ErrCode    int      `json:"errcode"`
	ErrMsg     string   `json:"errmsg"`
}

func GetWechatOpenid(code, appid, secret string) (wechatRsp WechatOpenidRsp, err error) {

	requestLine := strings.Join([]string{"https://api.weixin.qq.com/sns/oauth2/access_token",
		"?appid=", appid,
		"&secret=", secret,
		"&code=", code,
		"&grant_type=authorization_code"}, "")

	resp, err := http.Get(requestLine)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("http return code is not 200, code: %d", resp.StatusCode)
		err = errors.New(msg)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	atr := WechatOpenidRsp{}
	err = json.Unmarshal(body, &atr)
	if err != nil {
		return
	}
	//atr.ErrCode = 10004
	if atr.ErrCode != 0 || atr.OpenId == "" {
		msg := fmt.Sprintf("wechat openid get fail! body: %s", body)
		if atr.ErrCode == 10004 {
			msg = fmt.Sprintf("this appid is banned! body: %s", body)
		}
		err = errors.New(msg)
		return
	}

	wechatRsp = atr

	return

}
func GetWechatUserInfo(accessToken, openid string) (wechatRsp WechatUserInfoRsp, err error) {

	requestLine := strings.Join([]string{"https://api.weixin.qq.com/sns/userinfo",
		"?access_token=", accessToken,
		"&openid=", openid,
		"&lang=", "zh_CN"}, "")

	resp, err := http.Get(requestLine)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("http return code is not 200, code: %d", resp.StatusCode)
		err = errors.New(msg)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	atr := WechatUserInfoRsp{}
	err = json.Unmarshal(body, &atr)
	if err != nil {
		return
	}
	if atr.ErrCode != 0 {
		msg := fmt.Sprintf("wechat openid get fail! body: %s", body)
		err = errors.New(msg)
		return
	}

	wechatRsp = atr

	return

}

type AccessTokenRsp struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
}

//type AccessTokenErrorResponse struct {
//  Errcode float64
//  Errmsg  string
//}
//获取wx_AccessToken 拼接get请求 解析返回json结果 返回 AccessToken和err

//https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_4
func GetWechatAccessToken(appID, appSecret string, slave bool, rc *grcache.RedisConn) (accessToken string, err error) {

	accessToken = ""

	if slave {
		var exist bool
		exist, err = rc.Get(WECHAT_ACCESSTOKEN_REDIS_KEY_PREFIX+appID, &accessToken)
		if err != nil {
			err = fmt.Errorf("redis get error err: %s, key: %s", err, WECHAT_ACCESSTOKEN_REDIS_KEY_PREFIX+appID)
			return
		}

		if exist {
			return
		}

	}

	requestLine := strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/token",
		"?grant_type=client_credential&appid=",
		appID,
		"&secret=",
		appSecret}, "")

	resp, err := http.Get(requestLine)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("http return code is not 200, code: %d", resp.StatusCode)
		err = errors.New(msg)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	atr := AccessTokenRsp{}
	err = json.Unmarshal(body, &atr)
	if err != nil || atr.AccessToken == "" {
		msg := fmt.Sprintf("rsp json unmarshal fail! err: %s, body: %s", err, body)
		err = errors.New(msg)
		return
	}

	err = rc.SetEx(WECHAT_ACCESSTOKEN_REDIS_KEY_PREFIX+appID, atr.AccessToken, int(atr.ExpiresIn))
	if err != nil {
		err = fmt.Errorf("redis do error err: %s,config: %+v", err, rc.PoolConfig)
		return
	}

	accessToken = atr.AccessToken

	return

}

func GetWechatJsTicket(appId string, appSecret string, rc *grcache.RedisConn) (jsTicket string, err error) {

	jsTicket = ""

	exist, err := rc.Get(WECHAT_JSTICKET_REDIS_KEY_PREFIX+appId, &jsTicket)
	if err != nil {
		err = fmt.Errorf("redis get error err: %s, key: %s", err, WECHAT_JSTICKET_REDIS_KEY_PREFIX+jsTicket)
		return
	}

	if exist {
		return
	}

	accessToken, err := GetWechatAccessToken(appId, appSecret, true, rc)
	if err != nil {
		return
	}

	requestLine := strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=",
		accessToken,
		"&type=jsapi",
	}, "")

	resp, err := http.Get(requestLine)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("http return code is not 200, code: %d", resp.StatusCode)
		err = errors.New(msg)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	atr := WechatJsTicketRsp{}
	err = json.Unmarshal(body, &atr)
	if err != nil || atr.Ticket == "" {
		msg := fmt.Sprintf("rsp json unmarshal fail! err: %s, body: %s", err, body)
		err = errors.New(msg)
		return
	}

	err = rc.SetEx(WECHAT_JSTICKET_REDIS_KEY_PREFIX+appId, atr.Ticket, 4800)
	if err != nil {
		fmt.Errorf("redis do error err: %s,config: %+v", err, rc.PoolConfig)
		return
	}

	jsTicket = atr.Ticket

	return

}

type WechatJsTicketRsp struct {
	ErrorCode int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

//微信支付页面config签名
func WechatJsTicketSign(mReq map[string]interface{}) string {

	//fmt.Println("========STEP 1, 对key进行升序排序.========")
	//fmt.Println("微信支付签名计算, API KEY:", key)
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)

	//fmt.Println("========STEP2, 对key=value的键值对用&连接起来，略过空值========")
	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for i, k := range sorted_keys {
		//fmt.Printf("k=%v, v=%v\n", k, mReq[k])
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			if i != (len(sorted_keys) - 1) {
				signStrings = signStrings + k + "=" + value + "&"
			} else {
				signStrings = signStrings + k + "=" + value //最后一个不加此符号
			}
		}
	}

	//对字符串进行SHA1哈希
	t := sha1.New()
	io.WriteString(t, signStrings)
	upperSign := fmt.Sprintf("%x", t.Sum(nil))

	return upperSign
}
