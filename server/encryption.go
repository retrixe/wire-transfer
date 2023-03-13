package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type KeyFile struct {
	PublicKey  string `toml:"public_key"`
	PrivateKey string `toml:"private_key"`
}

var ecdhKey *ecdh.PrivateKey

func LoadEncryptionKeys() {
	var keyFile KeyFile
	keyFileData, err := os.ReadFile("ecdh_keys.toml")

	if err != nil && os.IsNotExist(err) {
		// Generate new keys.
		key, err := ecdh.X25519().GenerateKey(rand.Reader)
		if err != nil {
			log.Panicln("Failed to generate server key pair:", err)
		}
		ecdhKey = key

		// Save them to file.
		keyFile.PublicKey, err = MarshalPEMPublicKey(key.PublicKey())
		if err != nil {
			log.Panicln("Failed to marshal server public key:", err)
		}
		keyFile.PrivateKey, err = MarshalPEMPrivateKey(key)
		if err != nil {
			log.Panicln("Failed to marshal server private key:", err)
		}
		keyFileData, err = toml.Marshal(keyFile)
		if err != nil {
			log.Panicln("Failed to marshal encryption_keys.toml:", err)
		}
		err = os.WriteFile("ecdh_keys.toml", keyFileData, 0600)
		if err != nil {
			log.Panicln("Failed to create ecdh_keys.toml:", err)
		}
	} else if err != nil {
		log.Panicln("Failed to read ecdh_keys.toml:", err)
	}

	err = toml.Unmarshal(keyFileData, &keyFile)
	if err != nil {
		log.Panicln("Failed to parse ecdh_keys.toml:", err)
	}

	publicKey, err := ParsePEMPublicKey(keyFile.PublicKey)
	if err != nil {
		log.Panicln("Failed to parse server public key:", err)
	}
	privateKey, err := ParsePEMPrivateKey(keyFile.PrivateKey)
	if err != nil {
		log.Panicln("Failed to parse server private key:", err)
	}
	if !publicKey.Equal(privateKey.PublicKey()) {
		log.Panicln("Server public and private keys do not match!")
	}
	ecdhKey = privateKey
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
