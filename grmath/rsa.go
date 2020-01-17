package grmath

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

/*
privateKey := "-----BEGIN RSA PRIVATE KEY-----\n"
	privateKeyAccessinfo := config.AccessInfo["pri_key"]
	if privateKeyAccessinfo == "" {
		l4g.Error("[SDK][iapppayh5]not declare the PrivateKey in AccessInfo")
		ret = sdkServerArgsError
		return
	}

	privateKey += privateKeyAccessinfo
	privateKey += "\n-----END RSA PRIVATE KEY-----"

	l4g.Debug("[SDK][iapppayh5]private_key:%s", privateKey)
*/

func GetPKCS1RSASignByMD5(origData []byte, priKey string) ([]byte, error) {
	block, _ := pem.Decode([]byte(priKey))
	if block == nil {

		return nil, errors.New("pem Decode error interface{} trun fail")
	}

	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {

		return nil, err
	}

	privateKey := pri

	h := md5.New()
	h.Write(origData)
	sum := h.Sum(nil)
	str, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.MD5, sum)
	if err != nil {

		return nil, err

	}

	signs := []byte(base64.StdEncoding.EncodeToString([]byte(str)))
	return signs, nil

}

/*
privateKey := "-----BEGIN RSA PRIVATE KEY-----\n"
	privateKeyAccessinfo := config.AccessInfo["pri_key"]
	if privateKeyAccessinfo == "" {
		l4g.Error("[SDK][QingYuan]not declare the PrivateKey in AccessInfo")
		ret = sdkServerArgsError
		return
	}

	privateKey += privateKeyAccessinfo
	privateKey += "\n-----END RSA PRIVATE KEY-----"

	l4g.Debug("[SDK][QingYuan]private_key:%s", privateKey)
	l4g.Debug("[SDK][QingYuan]unsigned_str:%s", unsigned_str)
	ret, sign := ymutil.RsaEncrypt(unsigned_str, privateKey, "pkcs8")
	tmp_map["Sign"] = sign
*/

// ============================
/*
sign, err = getRSASignByMD5([]byte(origData), config.AccessInfo["rsa_key"])
	if err != nil {
		ret = sdkServerArgsError
		l4g.Error("[SDK][Qbao]getSign error %s", err)
		return
	}
*/

func GetRSASignByMD5(origData []byte, priKey string) ([]byte, error) {

	data, err := base64.StdEncoding.DecodeString(priKey)
	if err != nil {

		return nil, err
	}

	pri, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {

		return nil, err
	}

	privateKey, ok := pri.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("interface{} trun fail")
	}

	h := md5.New()
	h.Write(origData)
	sum := h.Sum(nil)
	str, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.MD5, sum)
	if err != nil {

		return nil, err

	}

	signs := []byte(base64.StdEncoding.EncodeToString([]byte(str)))
	return signs, nil

}

func FormatPrivateKey(rawPrivateKey string) (privateKey string) {
	privateKey = "-----BEGIN RSA PRIVATE KEY-----\n"

	length := len(rawPrivateKey)

	for i := 64; ; i = i + 64 {
		if i > length {
			privateKey += rawPrivateKey[i-64 : length]
			break
		} else {
			privateKey += rawPrivateKey[i-64 : i]
		}

	}

	privateKey += "\n-----END RSA PRIVATE KEY-----"

	return
}

func GetRSASignBySha(origData string, privateKeyString string) (sign []byte, err error) {

	block, _ := pem.Decode([]byte(privateKeyString))

	if block == nil {

		return nil, errors.New(" pem Decode error interface{} trun fail")
	}

	tprivateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	privateKey, ok := tprivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("interface{} trun fail")
	}

	h := sha1.New()
	h.Write([]byte(origData))
	digest := h.Sum(nil)

	s, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA1, digest)
	if err != nil {

		return nil, err
	}
	data := base64.StdEncoding.EncodeToString(s)
	return []byte(data), nil
}
