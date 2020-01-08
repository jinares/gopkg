package xcrypt

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
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
func Sha256(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}
