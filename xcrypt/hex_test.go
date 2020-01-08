package xcrypt

import (
	"fmt"
	"testing"
)

func TestHexDecodeStr(t *testing.T) {
	data := `
当 Web 自信满满，步入移动时代之时，它还没有做好充足的准备。

`
	encodestr := HexEncodeStr(data)
	decodestr, err := HexDecodeStr(encodestr)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	if data != decodestr {
		t.Fatal("fail")
	}
	fmt.Println("加密:", encodestr, "解密:", decodestr)
	p, v, err := GenRsaKeyWithPKCS1(4096)
	fmt.Println(p, err)
	fmt.Println(v, err)
	ed, err := RsaEncryptPKCS1v15(p, []byte(data))
	if err != nil {
		fmt.Println("加密rsa:", err.Error(), len(data))
		return
	}
	sd, err := RsaDecryptPKCS1v15(v, ed)
	fmt.Println("解密:", string(sd), err)

	sign, err := RsaSignPKCS1v15WithSHA256(v, []byte(data))
	fmt.Println("sign", HexEncodeStr(string(sign)), err)
	err = RsaVerfySignPKCS1v15WithSHA256([]byte(data), sign, p)
	fmt.Println("rsa very sign:", err)

}

type (
	Pay struct {
		SignType  string          `json:"sign_type"`
		SignMsg   string          `json:"sign_msg"`
		UserID    string          `json:"user_id"`
		SessionID string          `json:"session_id"` //optional
		Data      PayParamPackage `json:"data"`       //user_id==data.user_id
	}
	PayParamPackage struct {
		UserID    string `json:"user_id"`
		ProductID string `json:"product_id"`
		Rmb       int    `json:"rmb"`
		Num       int    `json:"num"`
		Ext       string `json:"ext"`

		Expire     string `json:"expire"`      //time.RFC3339 过期时间
		Times      int    `json:"times"`       //使用次数
		SerialNO   string `json:"serial_no"`   //流水号
		MerchantNo string `json:"merchant_no"` //商户号
	}
)
