package xtools

import (
	"crypto/rand"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"io"
	"strconv"
	"strings"
)

//HashID Hash
func HashID(id string, m int64) int64 {
	k := MD5(id)
	hash := int64(0)
	for v := range k {
		va := int64(k[v])
		fh := strconv.FormatInt(va, 16)
		fo, _ := strconv.ParseInt(fh, 10, 64)
		hash += fo
	}
	hash = (hash * 1) % m
	return hash + 1

}

//RandID 随机ID
func RandID(width int) string {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	ustr := ""
	for _, val := range b {
		ustr = ustr + fmt.Sprintf("%d", val)
	}
	if len(ustr) >= width {
		return strings.ToUpper(ustr)[0:width]
	}
	return ustr
}

//Guid Guid
func Guid() string {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return strings.ToUpper(uuid)
}

/*
//OrderID OrderID
func OrderID(userid string) string {

	// 20190718  99  7724592251254372
	//   20060102  12  26-12=14
	// [时间8位][99 2位][用户id hash 2位][机器位 2 位][顺序号 4位][随机数 8] = 8+2+2+2+4+8=26
	datetime := time.Now()
	datetime.Second()
	hostname, _ := os.Hostname()

	hashUser := fmt.Sprintf("%02s", xtools.ToStr(HashID(userid, int64(orderModVal))))

	//datetimestr := time.Now().Format("20060102150405") // 14
	hashHost := fmt.Sprintf("%0.1d", HashID(hostname, 9))
	randstr := datetime.Format("150405") + hashHost + randid(7)

	return datetime.Format("20060102") + "99" + hashUser + randstr
}
*/

func ObjectIDCounter(w int) string {
	serial := fmt.Sprintf("%d", bson.NewObjectId().Counter())
	if w < 1 {
		return serial
	}
	l := len(serial)
	if l > w {
		return serial[l-w:]
	}
	if l == w {
		return serial
	}
	return serial + RandID(w-l)
}
