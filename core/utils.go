package core

import (
	"crypto/ecdh"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func boolToInt(b bool) uint8 {
	if b {
		return 0x01
	}
	return 0x00
}

func MarshalPEMPublicKey(key *ecdh.PublicKey) (string, error) {
	pkixKey, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pkixKey,
	})), nil
}

func MarshalPEMPrivateKey(key *ecdh.PrivateKey) (string, error) {
	pkcs8Key, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Key,
	})), nil
}

func ParsePEMPublicKey(data string) (*ecdh.PublicKey, error) {
	block, _ := pem.Decode([]byte(data))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("no valid PEM block found")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	return key.(*ecdh.PublicKey), err
}

func ParsePEMPrivateKey(data string) (*ecdh.PrivateKey, error) {
	block, _ := pem.Decode([]byte(data))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("no valid PEM block found")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	return key.(*ecdh.PrivateKey), err
}
