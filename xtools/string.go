package xtools

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//MD5  对一个字符串进行MD5加密,不可解密
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s)) //使用zhifeiya名字做散列值，设定后不要变
	return hex.EncodeToString(h.Sum(nil))
}

//Sha1 获取 SHA1 字符串
func Sha1(s string) string {
	t := sha1.New()
	t.Write([]byte(s))

	return hex.EncodeToString(t.Sum(nil))
}

//Base64Encode  base64加密
func Base64Encode(str string) string {
	var src []byte = []byte(str)
	return string([]byte(base64.StdEncoding.EncodeToString(src)))
}

//Base64Decode base64解密
func Base64Decode(str string) (string, error) {
	var src []byte = []byte(str)
	by, err := base64.StdEncoding.DecodeString(string(src))
	return string(by), err
}

//CRC32
func CRC32(data string) int64 {
	return int64(crc32.ChecksumIEEE([]byte(data)))
}

//Sleep Sleep
func Sleep(second time.Duration) {

	time.Sleep(second)
}

//Version 秒
func Version() int64 {
	return time.Now().Unix()
}

// MsTime 毫秒
func MsTime() int64 {

	return int64(time.Now().UnixNano()) / 1000000
}

//IntVal parse int
func IntVal(data string) (int64, bool) {
	//fmt.Println(regexp.MatchString("^[0-9]+$", "126645455"))
	//fmt.Println(regexp.MatchString("^[0-9]+.[0-9]*$", "126645455.2323"))
	//fmt.Println(strconv.ParseFloat("1656.455ssss5", 64))

	if ret, err := regexp.MatchString("^[0-9]+$", data); ret == true && err == nil {

		dataVal, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			return 0, false
		}
		return dataVal, true
	}
	return int64(0), false

}

//FloatVal parse float
func FloatVal(data string) (float64, bool) {
	if ret, err := regexp.MatchString("^[0-9]+\\.[0-9]+$", data); ret == true && err == nil {

		dataVal, err := strconv.ParseFloat(data, 64)

		if err != nil {
			return float64(0), false
		}
		return dataVal, true
	}
	dataVal, err := strconv.ParseFloat(data, 64)

	if err == nil {
		return dataVal, true
	}

	return float64(0), false
}

//ToStr ToStr
func ToStr(dval interface{}) string {
	switch vv := dval.(type) {
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", vv)
	case float32, float64:
		return fmt.Sprintf("%g", vv)
	case string:
		return string(vv)
	default:
		return ""
	}
	return ""

}

//UrlDecode UrlDecode
func UrlDecode(data string) map[string]string {

	ret := make(map[string]string)
	//var ret map[string]string
	vals, err := url.ParseQuery(data)
	if err != nil {
		return ret
	}
	for k, _ := range vals {
		ret[k] = vals.Get(k)
	}
	return ret

}

// UrlEncode UrlEncode
func UrlEncode(data map[string]string) string {
	vals, err := url.ParseQuery("")
	if err != nil {
		return ""
	}

	for k, v := range data {
		vals.Set(k, ToStr(v))
	}
	return vals.Encode()
}

//UTF8ToGBK UTF8ToGBK
func UTF8ToGBK(data string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(data)), simplifiedchinese.GBK.NewEncoder()))
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

//GbkToUtf8 GbkToUtf8
func GBKToUTF8(data string) (string, error) {

	ret, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(data)), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

// Has 数组中包含是否包含查找的对象
func Has(key interface{}, val ...interface{}) bool {
	skeyType := reflect.TypeOf(key).String()

	for _, itemVal := range val {
		if skeyType != reflect.TypeOf(itemVal).String() {
			continue
		}
		if key == itemVal {
			return true
		}
	}
	return false
}

//SubStr SubStr
func SubStr(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}
