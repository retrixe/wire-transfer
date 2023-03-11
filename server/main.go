package main

import (
	"log"
	"net"

	"github.com/retrixe/wire-transfer/core"
)

func main() {
	LoadConfig()
	LoadEncryptionKeys()

	server, err := net.ListenUDP("udp", &net.UDPAddr{Port: config.Port})
	if err != nil {
		log.Fatalln(err)
	}
	defer server.Close()
	log.Println("Listening on port", config.Port)

	for {
		buf := make([]byte, 1024)
		n, addr, err := server.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		data := make([]byte, n)
		copy(data, buf)

		packet, err := core.ParsePacket(data)
		if err != nil {
			continue
		}

		respond := func(b []byte) { server.WriteToUDP(b, addr) }
		switch packet.ID {
		case core.InfoPacketId:
			go HandleInfoPacket(packet, respond)
		}
	}
}
