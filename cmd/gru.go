package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/ushmodin/avaxo2/internal/gru"
	"github.com/ushmodin/avaxo2/internal/model"
	"github.com/ushmodin/avaxo2/internal/settings"
)

// CmdGru CLI Command run avaxo minion
var CmdGru = &cli.Command{
	Name:  "gru",
	Usage: "run gru (minion manager)",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "minion, m",
			Usage:    "target minion",
			Required: true,
		},
	},
	Subcommands: []*cli.Command{
		&cli.Command{
			Name:   "ls",
			Usage:  "list directory contents",
			Action: ls,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "path",
					Usage:    "target path",
					Required: true,
				},
				&cli.BoolFlag{
					Name:  "json",
					Usage: "print in json",
				},
			},
		},
	},
}

func runGru(ctx *cli.Context) error {
	return nil
}

func ls(ctx *cli.Context) error {
	if err := settings.InitSettings(); err != nil {
		return err
	}

	g, err := gru.NewGru(
		settings.GruSettings.Certfile,
		settings.GruSettings.Keyfile,
		settings.GruSettings.Cafile,
	)
	if err != nil {
		return err
	}
	minion := ctx.String("minion")
	path := ctx.String("path")
	isJSON := ctx.IsSet("json")

	files, err := g.Ls(minion, path)
	if err != nil {
		return err
	}

	if isJSON {
		out, _ := json.Marshal(files)
		fmt.Println(string(out))
	} else {
		out := model.PrintFiles(files)
		fmt.Println(string(out))
	}

	return nil
}
