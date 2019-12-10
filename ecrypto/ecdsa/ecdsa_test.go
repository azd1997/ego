package ecdsa

import (
	"crypto/sha256"
	"testing"

	"github.com/azd1997/ego/ecrypto"
)

func TestECDSA_Sign(t *testing.T) {
	ac := ECDSA{}
	priv, pub, err := ac.GenerateKeyPair()
	if err != nil {
		t.Error(err)
	}

	str := []byte("Hello Eiger")
	hash256 := sha256.Sum256(str)
	hash := ecrypto.Hash(hash256[:])

	sig, err := ac.Sign(hash, priv)
	if err != nil {
		t.Error(err)
	}

	valid := ac.VerifySign(hash, sig, pub)
	if !valid {
		t.Error("invalid signature")
	}
}
