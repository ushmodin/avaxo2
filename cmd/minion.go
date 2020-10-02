package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/ushmodin/avaxo2/internal/minion"
	"github.com/ushmodin/avaxo2/internal/settings"
)

var (
	listen string
	port   int
)

// CmdMinion CLI Command run avaxo minion
var CmdMinion = &cli.Command{
	Name:   "minion",
	Usage:  "run minion",
	Action: runAgent,
	Subcommands: []*cli.Command{
		&cli.Command{
			Name:      "service",
			Usage:     "Run as Windows service",
			ArgsUsage: "<Service Name>",
			Action:    runWinSrvAgent,
		},
	},
}

func runAgent(ctx *cli.Context) error {
	if err := settings.InitSettings(); err != nil {
		return err
	}
	server, err := minion.NewServer(
		settings.MinionSettings.Listen,
		settings.MinionSettings.Keyfile,
		settings.MinionSettings.Certfile,
		settings.MinionSettings.Cafile,
	)
	if err != nil {
		return err
	}
	return server.Run()
}

func runWinSrvAgent(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		cli.ShowCommandHelp(ctx, "service")
		return nil
	}
	return runWinSrv(ctx.Args().First())
}
