package xtools

import (
	"crypto/rand"
	"fmt"
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
