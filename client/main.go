package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/retrixe/wire-transfer/core"
	"github.com/urfave/cli/v2"
)

func main() {
	log.SetFlags(0)
	app := &cli.App{
		Name:  "wire-transfer",
		Usage: "Transfer files over the internet efficiently, using direct peer-to-peer connections where possible.",
		Commands: []*cli.Command{
			{
				Name:      "ping-server",
				Usage:     "Check if a server is online.",
				ArgsUsage: "[server address]",
				Action: func(c *cli.Context) error {
					if c.Args().Len() != 1 {
						return cli.Exit("No argument provided for server address!", 1)
					}
					conn, err := net.Dial("udp", c.Args().First())
					conn.SetDeadline(time.Now().Add(30 * time.Second))
					if err != nil {
						return cli.Exit("Failed to connect to server!", 1)
					}
					defer conn.Close()
					_, err = conn.Write(core.CreateInfoPacket(core.ProtocolVersion, nil).Serialize())
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
					log.Println("Server is online with public key: " + string(data.PublicKey))
					return nil
				},
			},
		},
		Suggest:                true,
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
