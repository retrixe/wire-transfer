package main

import (
	"crypto/x509"
	"log"

	"github.com/retrixe/wire-transfer/core"
)

func HandleInfoPacket(packet *core.Packet, respond func([]byte)) {
	var publicKey []byte = nil
	if ecdhKey != nil {
		pKey, err := x509.MarshalPKIXPublicKey(ecdhKey.PublicKey())
		if err != nil {
			log.Println("Failed to marshal public key in response to info request!", err)
			return
		}
		publicKey = pKey
	}
	respond(core.CreateInfoPacket(core.ProtocolVersion,
		publicKey,
		// TODO: Add config option for this
		nil,
		&config.MaxFileExpiryTime,
		"Reference implementation").Serialize())
}
