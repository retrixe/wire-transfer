package main

import (
	"crypto/ed25519"
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

type KeyPair struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

var encryptionKeys *KeyPair

func LoadEncryptionKeys() {
	var keyFile KeyFile
	keyFileData, err := os.ReadFile("ed25519_keys.toml")

	if err != nil && os.IsNotExist(err) {
		// Generate new keys.
		publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			log.Panicln("Failed to generate server key pair:", err)
		}
		encryptionKeys = &KeyPair{PublicKey: publicKey, PrivateKey: privateKey}

		// Save them to file.
		keyFile.PublicKey, err = MarshalPEMPublicKey(publicKey)
		if err != nil {
			log.Panicln("Failed to marshal server public key:", err)
		}
		keyFile.PrivateKey, err = MarshalPEMPrivateKey(privateKey)
		if err != nil {
			log.Panicln("Failed to marshal server private key:", err)
		}
		keyFileData, err = toml.Marshal(keyFile)
		if err != nil {
			log.Panicln("Failed to marshal encryption_keys.toml:", err)
		}
		err = os.WriteFile("ed25519_keys.toml", keyFileData, 0644)
		if err != nil {
			log.Panicln("Failed to create ed25519_keys.toml:", err)
		}
	} else if err != nil {
		log.Panicln("Failed to read ed25519_keys.toml:", err)
	}

	err = toml.Unmarshal(keyFileData, &keyFile)
	if err != nil {
		log.Panicln("Failed to parse ed25519_keys.toml:", err)
	}

	publicKey, err := ParsePEMPublicKey(keyFile.PublicKey)
	if err != nil {
		log.Panicln("Failed to parse server public key:", err)
	}
	privateKey, err := ParsePEMPrivateKey(keyFile.PrivateKey)
	if err != nil {
		log.Panicln("Failed to parse server private key:", err)
	}
	encryptionKeys = &KeyPair{PublicKey: publicKey, PrivateKey: privateKey}
}

func MarshalPEMPublicKey(key ed25519.PublicKey) (string, error) {
	pkixKey, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pkixKey,
	})), nil
}

func MarshalPEMPrivateKey(key ed25519.PrivateKey) (string, error) {
	pkcs8Key, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Key,
	})), nil
}

func ParsePEMPublicKey(data string) (ed25519.PublicKey, error) {
	block, _ := pem.Decode([]byte(data))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("no valid PEM block found")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	return key.(ed25519.PublicKey), err
}

func ParsePEMPrivateKey(data string) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode([]byte(data))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("no valid PEM block found")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	return key.(ed25519.PrivateKey), err
}
