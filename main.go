package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/ushmodin/avaxo2/cmd"
	"github.com/ushmodin/avaxo2/internal/settings"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config, c",
				Usage:    "Load configuration from `FILE`",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			cmd.CmdMinion,
			cmd.CmdGru,
		},
	}

	for _, c := range app.Commands {
		c.Before = beforeAllCommand
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func beforeAllCommand(ctx *cli.Context) error {
	settings.ConfigPath = ctx.String("config")
	return nil
}
