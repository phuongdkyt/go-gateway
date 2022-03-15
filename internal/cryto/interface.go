package cryto

import (
	"crypto/rsa"
)

type CrytoInterface interface {
	GetRsaPrivateKey() *rsa.PrivateKey
	GetRsaPublicKey() *rsa.PublicKey
	RsaDecrypt(base4Data string) (string, error)
	AeadEncrypt(raw, associated, key string) (string, error)
	AeadDecrypt(cipher, associated, key string) (string, error)
}
