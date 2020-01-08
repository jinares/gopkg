package xcrypt

import (
	"encoding/base64"
)

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
