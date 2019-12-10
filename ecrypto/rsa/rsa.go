package rsa

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/gob"

	"github.com/azd1997/ego/ecrypto"
	"github.com/azd1997/ego/utils"
)

// RSA 签名算法使用PSS； 加密使用OAEP；详见crypto/rsa
type RSA struct {}


func (r RSA) GenerateKeyPair() (ecrypto.PrivKey, ecrypto.PubKey, error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, nil, utils.WrapError("RSA_GenerateKeyPair", err)
	}
	privKey, err := utils.GobEncode(privateKey)
	if err != nil {
		return nil, nil, utils.WrapError("RSA_GenerateKeyPair", err)
	}
	pubKey, err := utils.GobEncode(privateKey.PublicKey)
	if err != nil {
		return nil, nil, utils.WrapError("RSA_GenerateKeyPair", err)
	}
	return privKey, pubKey, nil
}

func (r RSA) Sign(hash ecrypto.Hash, priv ecrypto.PrivKey) (ecrypto.Signature, error) {

	opts := rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	}

	// 注意opts内Hash函数必须和计算hash用的哈希函数一致。这里都采用了SHA256

	sig, err := rsa.SignPSS(rand.Reader, r.privateKey(priv), crypto.SHA256, hash, &opts)
	if err != nil {
		return nil, utils.WrapError("RSA_Sign", err)
	}

	return sig, nil
}

func (r RSA)  VerifySign(hash ecrypto.Hash, sig ecrypto.Signature, pub ecrypto.PubKey) bool {

	opts := rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	}

	err := rsa.VerifyPSS(r.publicKey(pub), crypto.SHA256, hash, sig, &opts)
	if err != nil {
		return false
	}

	return true
}

func (r RSA)  Encrypt(raw []byte, pub ecrypto.PubKey) (encrypted []byte, err error) {
	// raw 长度有限制。crypto原话是：
	// The message must be no longer than the length of the public modulus minus
	// twice the hash length, minus a further 2.
	// 消息的长度不得超过公用模数的长度减去哈希长度的两倍，再减去2。

	// label设置为空，详见 crypto/rsa EncryptOAEP()

	encrypted, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, r.publicKey(pub), raw, nil)
	if err != nil {
		return nil, utils.WrapError("RSA_Encrypt", err)
	}

	return encrypted, nil
}

func (r RSA)  Decrypt(encrypted []byte, priv ecrypto.PrivKey) (raw []byte, err error) {

	raw, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, r.privateKey(priv), encrypted, nil)
	if err != nil {
		return nil, utils.WrapError("RSA_Decrypt", err)
	}

	return raw, nil
}



func (r RSA) privateKey(priv ecrypto.PrivKey) *rsa.PrivateKey {
	privateKey := &rsa.PrivateKey{}
	_ = gob.NewDecoder(bytes.NewReader(priv)).Decode(privateKey)
	return privateKey
}

func (r RSA) publicKey(pub ecrypto.PubKey) *rsa.PublicKey {
	publicKey := &rsa.PublicKey{}
	_ = gob.NewDecoder(bytes.NewReader(pub)).Decode(publicKey)
	return publicKey
}