package secret

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"github.com/golang/glog"
)

// 私钥生成
//openssl genrsa -out rsa_private_key.pem 1024

// 公钥: 根据私钥生成
//openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem

var privateKey, publicKey []byte

func init() {
	/** 后续应该在这里生成 public.key 和 private_key比较好 */
	var err error
	publicKey, err = ioutil.ReadFile("./rsa/rsa_public_key.pem")
	if err != nil {
		glog.Errorf("init publicKey error, msg:[%s]", err.Error())
		return
	}

	privateKey, err = ioutil.ReadFile("./rsa/rsa_private_key.pem")
	if err != nil {
		glog.Errorf("init privateKey error, msg:[%s]", err.Error())
		return
	}
	glog.Infof("global key success")
}

// rsaEncrypt 加密
func rsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// rsaDecrypt 解密
func rsaDecrypt(ciphertext []byte) ([]byte, error) {
	var err error
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func Encrypt(str string) string {
	encrypt, err := rsaEncrypt([]byte(str))
	if err != nil {
		glog.Errorf("Encrypt error, msg:[%s]", err.Error())
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(encrypt)
}

func Decrypt(str string) []byte {
	strBytes, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		glog.Errorf("base64 Decode error, msg:[%s]", err.Error())
		return nil
	}
	decrypt, err := rsaDecrypt(strBytes)
	if err != nil {
		glog.Errorf("Decrypt error, msg:[%s]", err.Error())
		return nil
	}
	return decrypt
}
