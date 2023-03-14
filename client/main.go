package main

import (
	"log"
	"os"

	"github.com/retrixe/wire-transfer/client/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	log.SetFlags(0)
	app := &cli.App{
		Name:  "wire-transfer",
		Usage: "Transfer files over the internet efficiently, using direct peer-to-peer connections where possible.",
		Commands: []*cli.Command{
			cmd.TrustInfo,
			// wire-transfer file-info [server address] [file ID]
			// wire-transfer upload [server address] [file ID]
			// wire-transfer download [server address] [file ID]
		},
		Suggest:                true,
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
