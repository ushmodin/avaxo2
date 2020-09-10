package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/ushmodin/avaxo2/internal/gru"
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
			},
		},
	},
}

func runGru(ctx *cli.Context) error {
	return nil
}

func ls(ctx *cli.Context) error {
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

	files, err := g.Ls(minion, path)
	if err != nil {
		return err
	}
	printFiles(files)
	return nil
}

func printFiles(files []interface{}) {

}
