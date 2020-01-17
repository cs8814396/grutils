package grmath

import (
	"crypto/md5"
	//"crypto/sha1"
	"encoding/hex"
	"fmt"
	//"io"
	"sort"
	"strings"
)

/*
func SortOrderUpperSha1Sign(mReq map[string]interface{}, secretKey string) string {

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
	for _, k := range sorted_keys {
		//fmt.Printf("k=%v, v=%v\n", k, mReq[k])

		switch mReq[k].(type) {
		case string:
			tString := mReq[k].(string)
			if tString != "" {
				//signStrings = signStrings + k + "=" + value + "&"

				signStrings = fmt.Sprintf("%s%s=%s&", signStrings, k, tString)
			}
		case int:

			tInt := mReq[k].(int)
			if tInt != 0 {
				signStrings = fmt.Sprintf("%s%s=%d&", signStrings, k, tInt)
			}

		}
		//fmt.Println(k, reflect.TypeOf(mReq[k]), reflect.TypeOf(mReq[k]) == reflect.Int)

	}

	//fmt.Println("========STEP3, 在键值对的最后加上key=API_KEY========")
	//STEP3, 在键值对的最后加上key=API_KEY

	signStrings = signStrings + secretKey

	//fmt.Println("========STEP4, 进行MD5签名并且将所有字符转为大写.========")
	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := sha1.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))

	return upperSign
}*/

func SortOrderUpperMd5Sign(mReq map[string]interface{}, secretKey string) string {

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
	for _, k := range sorted_keys {
		//fmt.Printf("k=%v, v=%v\n", k, mReq[k])

		switch mReq[k].(type) {
		case string:
			tString := mReq[k].(string)
			if tString != "" {
				//signStrings = signStrings + k + "=" + value + "&"

				signStrings = fmt.Sprintf("%s%s=%s&", signStrings, k, tString)
			}
		case int:

			tInt := mReq[k].(int)
			if tInt != 0 {
				signStrings = fmt.Sprintf("%s%s=%d&", signStrings, k, tInt)
			}

		}
		//fmt.Println(k, reflect.TypeOf(mReq[k]), reflect.TypeOf(mReq[k]) == reflect.Int)

	}

	//fmt.Println("========STEP3, 在键值对的最后加上key=API_KEY========")
	//STEP3, 在键值对的最后加上key=API_KEY

	signStrings = signStrings + "key=" + secretKey

	//fmt.Println("========STEP4, 进行MD5签名并且将所有字符转为大写.========")
	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))

	return upperSign
}
