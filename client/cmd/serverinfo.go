package cmd

import (
	"crypto/ecdh"
	"crypto/x509"
	"log"
	"net"
	"strings"
	"time"

	"github.com/retrixe/wire-transfer/core"
	"github.com/urfave/cli/v2"
)

var TrustInfo = &cli.Command{
	Name:      "serverinfo",
	Aliases:   []string{"server-info", "si"},
	Usage:     "Get info about a file transfer server.",
	ArgsUsage: "[server address]",
	Action: func(c *cli.Context) error {
		if c.Args().Len() != 1 {
			return cli.Exit("No argument provided for server address!", 1)
		}
		serverAddress := c.Args().First()
		if !strings.Contains(serverAddress, ":") {
			serverAddress += ":14776"
		}
		conn, err := net.Dial("udp", serverAddress)
		conn.SetDeadline(time.Now().Add(30 * time.Second))
		if err != nil {
			return cli.Exit("Failed to connect to server!", 1)
		}
		defer conn.Close()
		_, err = conn.Write(core.CreateInfoPacket(core.ProtocolVersion, nil, nil, nil, "").Serialize())
		if err != nil {
			return cli.Exit("Failed to send info packet to server!", 1)
		}
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			return cli.Exit("Failed to read response from server!", 1)
		}
		packet, err := core.ParsePacket(buf[:n])
		if err != nil {
			return cli.Exit("Failed to parse response from server!", 1)
		} else if packet.ID != core.InfoPacketId {
			return cli.Exit("Unexpected response from server!", 1)
		}
		data, err := core.ParseInfoPacket(packet)
		if err != nil {
			return cli.Exit("Failed to parse info packet from server!", 1)
		}
		derBytes := data.PublicKey
		if derBytes == nil {
			return cli.Exit("Server did not provide an ECDH public key! Servers which don't support encryption cannot be trusted.", 1)
		}
		pKeyUncasted, err := x509.ParsePKIXPublicKey(derBytes)
		if err != nil {
			return cli.Exit("Failed to parse ECDH public key from server!", 1)
		}
		pKey, ok := pKeyUncasted.(*ecdh.PublicKey)
		if !ok {
			return cli.Exit("Server's public key is not of the correct type!", 1)
		}
		pKeyPemEncoded, err := core.MarshalPEMPublicKey(pKey)
		if err != nil {
			return cli.Exit("Failed to marshal ECDH public key from server!", 1)
		}
		pKeyEncoded := strings.Split(string(pKeyPemEncoded), "\n")[1]
		log.Println("ECDH public key:", pKeyEncoded)
		log.Println("Protocol version:", data.Version)
		log.Println("Maximum acceptable file size:", data.MaxFileSize)
		log.Println("Maximum acceptable expiry time:", data.MaxExpiryTime)
		log.Println("Extra information:", data.Info)
		return nil
	},
}
