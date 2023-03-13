package main

import (
	"crypto/ecdh"
	"crypto/x509"
	"log"
	"net"
	"time"

	"github.com/retrixe/wire-transfer/core"
)

func HandleInfoPacket(packet *core.Packet, server *net.UDPConn, addr *net.UDPAddr) {
	var publicKey []byte = nil
	if ecdhKey != nil {
		pKey, err := x509.MarshalPKIXPublicKey(ecdhKey.PublicKey())
		if err != nil {
			log.Println("Failed to marshal public key in response to info request!", err)
			return
		}
		publicKey = pKey
	}
	server.WriteToUDP(core.CreateInfoPacket(
		core.ProtocolVersion,
		publicKey,
		// TODO: Add config option for this
		nil,
		&config.MaxFileExpiryTime,
		"Reference implementation",
	).Serialize(), addr)
}

func HandleHandshakePacket(packet *core.Packet, server *net.UDPConn, addr *net.UDPAddr) {
	data, err := core.ParseHandshakeRequestPacket(packet)
	if err != nil {
		log.Println("Failed to parse handshake packet from", addr, "!", err)
		return
	}
	if data.Version != core.ProtocolVersion {
		server.WriteToUDP(core.CreateClosePacket("Unsupported protocol version!").Serialize(), addr)
		return
	}
	connection := &Connection{
		Addr:        addr,
		LastSeen:    time.Now(),
		TimeoutChan: make(chan bool),
	}
	if data.PublicKey != nil {
		pKeyUncasted, err := x509.ParsePKIXPublicKey(data.PublicKey)
		if err != nil {
			server.WriteToUDP(core.CreateClosePacket("Invalid public key!").Serialize(), addr)
			return
		} else if pKey, ok := pKeyUncasted.(*ecdh.PublicKey); !ok {
			server.WriteToUDP(core.CreateClosePacket("Invalid public key!").Serialize(), addr)
			return
		} else {
			connection.SharedSecret, err = ecdhKey.ECDH(pKey)
			if err != nil {
				server.WriteToUDP(core.CreateClosePacket("Invalid public key!").Serialize(), addr)
				return
			}
		}
	} else if config.RequireEncryption {
		server.WriteToUDP(core.CreateClosePacket("Encryption is required!").Serialize(), addr)
		return
	}
	connections.Store(addr.String(), connection)

	// Setup timeout.
	go func() {
		for {
			var res interface{}
			switch res {
			case <-connection.TimeoutChan:
				connection.LastSeen = time.Now()
				continue
			case <-time.After(time.Millisecond * time.Duration(config.UDPTimeoutDuration)):
				connections.Delete(connection.Addr.String())
				server.WriteToUDP(core.CreateClosePacket("Connection timed out!").Serialize(), addr)
				return
			}
			if res == true {
				return
			}
		}
	}()
	var publicKey []byte = nil
	if ecdhKey != nil {
		pKey, err := x509.MarshalPKIXPublicKey(ecdhKey.PublicKey())
		if err != nil {
			log.Println("Failed to marshal public key in response to info request!", err)
			return
		}
		publicKey = pKey
	}
	server.WriteToUDP(core.CreateHandshakeResponsePacket(
		core.ProtocolVersion,
		publicKey,
		// TODO: Add config option for this
		nil,
		&config.MaxFileExpiryTime,
		"Reference implementation",
	).Serialize(), addr)
}

func HandleClosePacket(packet *core.Packet, server *net.UDPConn, addr *net.UDPAddr) {
	conn, ok := connections.LoadAndDelete(addr.String())
	if !ok {
		return
	}
	conn.TimeoutChan <- true
}
