package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/retrixe/wire-transfer/core"
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
		keyFile.PublicKey, err = core.MarshalPEMPublicKey(key.PublicKey())
		if err != nil {
			log.Panicln("Failed to marshal server public key:", err)
		}
		keyFile.PrivateKey, err = core.MarshalPEMPrivateKey(key)
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

	publicKey, err := core.ParsePEMPublicKey(keyFile.PublicKey)
	if err != nil {
		log.Panicln("Failed to parse server public key:", err)
	}
	privateKey, err := core.ParsePEMPrivateKey(keyFile.PrivateKey)
	if err != nil {
		log.Panicln("Failed to parse server private key:", err)
	}
	if !publicKey.Equal(privateKey.PublicKey()) {
		log.Panicln("Server public and private keys do not match!")
	}
	ecdhKey = privateKey
}
