package xcrypt

import (
	"encoding/hex"
)

/*
hex.DecodeString(s string)//解密
hex.EncodeToString(src []byte) string//加密
*/
//解密
func HexDecodeStr(s string) (string, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b), err
}

//加密
func HexEncodeStr(s string) string {
	return hex.EncodeToString([]byte(s))
}
