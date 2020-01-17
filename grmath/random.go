package grmath

import (
	"crypto/rand"
	"errors"
	"fmt"
	"grutils/grcommon"

	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var orderLock sync.Mutex
var payIDCounter uint64
var ProcessCode string

const (
	MaxUint32 = 1<<32 - 1
)

func GenRandString(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	if randType == "salphanum" {
		dictionary = "0123456789abcdefghijklmnopqrstuvwxyz"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func OrderInt64Gen() (orderId int64, create_time string) {
	orderLock.Lock()
	defer orderLock.Unlock()

	create_time = GetFormatCurTime()

	b := []byte(create_time)
	timestr := string(b[0:4]) + string(b[5:7]) + string(b[8:10]) + string(b[11:13]) + string(b[14:16]) + string(b[17:19])

	//fmt.Println(timestr)

	payIDCounter = (payIDCounter + 1) % 100000

	host, _ := os.Hostname()

	temp := strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(os.Getpid()) + host //strconv.Itoa(payIDCounter) + strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(os.Getpid()) + host

	uniqnum, _ := strconv.ParseUint(grcommon.SubStrOnBytes(Md5(temp), 0, 5), 16, 32)

	tradeNO := grcommon.SubStrOnBytes(fmt.Sprintf("%s%07d", timestr, uniqnum), 0, 17)

	//fmt.Println(tradeNO)

	orderId, _ = strconv.ParseInt(tradeNO, 10, 64)

	return
}

func OrderCodeGen(preFix string) (orderId string, create_time string) {

	orderLock.Lock()
	defer orderLock.Unlock()

	create_time = GetFormatCurTime()

	b := []byte(create_time)
	timestr := string(b[0:4]) + string(b[5:7]) + string(b[8:10]) + string(b[11:13]) + string(b[14:16]) + string(b[17:19])

	payIDCounter = (payIDCounter + 1) % 100000

	host, _ := os.Hostname()

	temp := strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(os.Getpid()) + host //strconv.Itoa(payIDCounter) + strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(os.Getpid()) + host

	uniqnum, _ := strconv.ParseUint(grcommon.SubStrOnBytes(Md5(temp), 0, 5), 16, 32)

	tradeNO := grcommon.SubStrOnBytes(fmt.Sprintf("%s%s%07d", preFix, timestr, uniqnum), 0, 21)

	return tradeNO, create_time

}

func DistributedUniqueCodeGen(preFix string) (orderId string, create_time string) {

	orderLock.Lock()
	defer orderLock.Unlock()

	create_time = GetFormatCurTime()

	b := []byte(create_time)
	timestr := string(b[0:4]) + string(b[5:7]) + string(b[8:10]) + string(b[11:13]) + string(b[14:16]) + string(b[17:19])

	payIDCounter = (payIDCounter + 1) % 100000

	host, _ := os.Hostname()

	millisecond := fmt.Sprintf("%d", time.Now().UnixNano()/1000000)[10:]

	temp := strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(os.Getpid()) + host //strconv.Itoa(payIDCounter) + strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(os.Getpid()) + host

	uniqnum, _ := strconv.ParseUint(grcommon.SubStrOnBytes(Md5(temp), 0, 5), 16, 32)

	//fmt.Println(uniqnum)

	//fmt.Println(millisecond)

	tradeNO := grcommon.SubStrOnBytes(fmt.Sprintf("%s%s%s%07d", preFix, timestr, millisecond, uniqnum), 0, 24)

	return tradeNO, create_time

}

func GetCreateTimeFromOrder(preFix string, orderId string) (err error, create_time string) {

	arr := strings.Split(orderId, preFix)
	if len(arr) != 2 || len(arr[1]) <= 14 {
		msg := fmt.Sprintf("get createTime fail from order: %s, prefix: %s", orderId, preFix)
		err = errors.New(msg)
	}

	create_time = fmt.Sprintf("%s-%s-%s %s:%s:%s", arr[1][0:4], arr[1][4:6], arr[1][6:8], arr[1][8:10], arr[1][10:12], arr[1][12:14])

	return

}

func GetCreateTimeFromOrderInt64(orderId int64) (err error, create_time string) {

	orderIdString := strconv.FormatInt(orderId, 10)

	create_time = fmt.Sprintf("%s-%s-%s %s:%s:%s", orderIdString[0:4], orderIdString[4:6], orderIdString[6:8], orderIdString[8:10], orderIdString[10:12], orderIdString[12:14])

	return

}
func init() {
	ProcessCode = GenRandString(64, "alphanum")
}
