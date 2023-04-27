package authkeys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func GenerateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func MarshalPrivateKey(key *rsa.PrivateKey) []byte {
	data := x509.MarshalPKCS1PrivateKey(key)
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: data,
	})
}

func ParsePrivateKey(pemData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}
	if block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to parse PEM block type: %s", block.Type)
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
