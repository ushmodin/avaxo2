package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/ushmodin/avaxo2/cmd"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config, c",
				Usage: "Load configuration from `FILE`",
			},
		},
		Commands: []*cli.Command{
			cmd.CmdAgent,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
