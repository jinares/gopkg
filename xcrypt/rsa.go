package xcrypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
)

//验签：对采用sha256算法
//RsaVerfySignPKCS1v15WithSHA256
func RsaVerfySignPKCS1v15WithSHA256(originalData, signData []byte, pubKey string) error {

	pb, _ := Base64Decode(pubKey)
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey([]byte(pb))
	if err != nil {
		return err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	if err != nil {
		return err
	}
	hash := sha256.New()
	hash.Write(originalData)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash.Sum(nil), signData)
}

// 签名 对采用sha256算法
//RsaSignPKCS1v15WithSHA256
func RsaSignPKCS1v15WithSHA256(privateKey string, data []byte) ([]byte, error) {

	pb, _ := Base64Decode(privateKey)
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey([]byte(pb))
	if err != nil {
		return nil, err
	}
	hash := sha256.New()
	hash.Write(data)
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash.Sum(nil))
}

// 加密
func RsaEncryptPKCS1v15(publicKey string, origData []byte) ([]byte, error) {

	pb, _ := Base64Decode(publicKey)
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey([]byte(pb))
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecryptPKCS1v15(privateKey string, ciphertext []byte) ([]byte, error) {
	//解密
	pb, _ := Base64Decode(privateKey)
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey([]byte(pb))
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

//RSA公钥私钥产生 PKCS1   公钥私钥 进行base64返回
func GenRsaKeyWithPKCS1(bits int) (pubkey, prikey string, err error) {
	// 生成私钥文件

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)

	prikey = Base64Encode(string(derStream))
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return
	}
	pubkey = Base64Encode(string(derPkix))
	return
}
