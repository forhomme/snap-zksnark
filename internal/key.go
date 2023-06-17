package internal

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"smart-contract-service/configuration"
)

func GeneratePrivateKey(cfg configuration.ConfigApp) (*rsa.PrivateKey, error) {
	priv, err := os.ReadFile(cfg.PrivateKeyLocation)
	if err != nil {
		return nil, err
	}

	privPem, _ := pem.Decode(priv)
	var privPemBytes []byte
	if privPem.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("Not RSA private key")
	}
	privPemBytes = privPem.Bytes

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			return nil, err
		}
	}

	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, err
	}
	return privateKey, nil
}
