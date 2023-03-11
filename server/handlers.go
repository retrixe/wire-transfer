package main

import (
	"crypto/x509"
	"log"

	"github.com/retrixe/wire-transfer/core"
)

func HandleInfoPacket(packet *core.Packet, respond func([]byte)) {
	var publicKey []byte = nil
	if encryptionKeys != nil {
		pKey, err := x509.MarshalPKIXPublicKey(encryptionKeys.PublicKey)
		if err != nil {
			log.Println("Failed to marshal public key in response to info request!", err)
			return
		}
		publicKey = pKey
	}
	respond(core.CreateInfoPacket(core.ProtocolVersion, publicKey).Serialize())
}
