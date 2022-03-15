package cryto

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/insecurecleartextkeyset"
	"github.com/google/tink/go/keyset"
	"github.com/google/tink/go/tink"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
)

var ProviderCrytoSet = wire.NewSet(NewCrytoImpl)

type CrytoImpl struct {
	ctx        context.Context
	logger     *zap.Logger
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (c CrytoImpl) RsaDecrypt(base4Data string) (string, error) {
	encryptedKey, err := base64.StdEncoding.DecodeString(base4Data)
	if err != nil {
		return "", err
	}
	// operation.
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, c.privateKey, encryptedKey, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(plaintext), nil
}

func (c CrytoImpl) AeadEncrypt(raw, associated, key string) (string, error) {
	aeadInstance, err := constructAEADInstance(key)
	if err != nil {
		return "", err
	}
	cipherText, err := aeadInstance.Encrypt([]byte(raw), []byte(associated))
	if err != nil {
		return "", nil
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (c CrytoImpl) AeadDecrypt(cipher, associated, key string) (string, error) {
	aeadInstance, err := constructAEADInstance(key)
	if err != nil {
		return "", err
	}
	decodeString, err := base64.StdEncoding.DecodeString(cipher)
	if err != nil {
		return "", err
	}
	plain, err := aeadInstance.Decrypt(decodeString, []byte(associated))
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func (c CrytoImpl) GetRsaPrivateKey() *rsa.PrivateKey {
	return c.privateKey
}

func (c CrytoImpl) GetRsaPublicKey() *rsa.PublicKey {
	return c.publicKey
}

func NewCrytoImpl(ctx context.Context, logger *zap.Logger) CrytoInterface {
	key, err := parseRsaPrivateKey(viper.GetString("RSA_PRIVATE_KEY"))
	if err != nil {
		log.Fatal("error parse RSA private key: ", err)
	}
	pub, err := base64.StdEncoding.DecodeString(viper.GetString("AUTH_PUBLIC_KEY"))
	publicKey, err := ParsePublicKeyRSA(pub)
	if err != nil {
		log.Fatal("error parse RSA public key: ", err)
	}

	return &CrytoImpl{
		ctx:        ctx,
		logger:     logger,
		privateKey: key,
		publicKey:  publicKey,
	}
}

func parseRsaPrivateKey(base64Key string) (*rsa.PrivateKey, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}
	privPem, _ := pem.Decode(key)
	var privPemBytes []byte
	if privPem.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("RSA private key is of the wrong type")
	}
	privPemBytes = privPem.Bytes
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		return nil, err
	}
	return parsedKey.(*rsa.PrivateKey), nil
}

func ParsePublicKeyRSA(key []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, fmt.Errorf("public key: malformed or missing PEM format (RSA)")
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, err
		}
	}

	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key: expected a type of *rsa.PublicKey")
	}

	return publicKey, nil
}

func constructAEADInstance(key string) (tink.AEAD, error) {
	raw, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	kh, err := insecurecleartextkeyset.Read(keyset.NewBinaryReader(bytes.NewReader(raw)))
	if err != nil {
		return nil, err
	}
	return aead.New(kh)
}
