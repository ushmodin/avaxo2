package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

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
		&cli.Command{
			Name:            "get",
			Usage:           "Get file from minion",
			ArgsUsage:       "<remote path> [local path]",
			Action:          get,
			SkipFlagParsing: true,
		},
		&cli.Command{
			Name:            "put",
			Usage:           "Put file to minion",
			ArgsUsage:       "<local path> <remote path>",
			Action:          put,
			SkipFlagParsing: true,
		},
		&cli.Command{
			Name:      "exec",
			Usage:     "Execute command on minion",
			ArgsUsage: "<command> [args...]",
			Action:    exec,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "nowait",
					Usage: "Do't wait for the end of execution",
				},
				&cli.IntFlag{
					Name:        "timeout",
					Usage:       "Maximum waiting time (seconds)",
					DefaultText: "60",
				},
			},
		},
	},
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

func get(ctx *cli.Context) error {
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

	if ctx.NArg() < 1 {
		cli.ShowCommandHelp(ctx, "get")
		return nil
	}

	minion := ctx.String("minion")
	remote := ctx.Args().Get(0)
	var dest io.WriteCloser
	if ctx.NArg() > 1 {
		dest, err = os.OpenFile(ctx.Args().Get(1), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer dest.Close()
	} else {
		dest = os.Stdout
	}

	reader, err := g.GetFile(minion, remote)
	if err != nil {
		return err
	}
	defer reader.Close()

	if _, err := io.Copy(dest, reader); err != nil {
		return err
	}

	return nil
}

func put(ctx *cli.Context) error {
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

	if ctx.NArg() < 2 {
		cli.ShowCommandHelp(ctx, "put")
		return nil
	}

	minion := ctx.String("minion")
	remote := ctx.Args().Get(1)
	src, err := os.Open(ctx.Args().Get(0))
	if err != nil {
		return err
	}
	defer src.Close()

	return g.PutFile(minion, remote, src)
}

func exec(ctx *cli.Context) error {
	nowait := ctx.Bool("nowait")
	timeout := ctx.Int("timeout")
	if !ctx.IsSet("timeout") {
		timeout = 60
	}

	if ctx.NArg() < 1 {
		cli.ShowCommandHelp(ctx, "exec")
		return nil
	}
	cmd := ctx.Args().First()
	args := ctx.Args().Tail()

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

	return g.Exec(cmd, args, nowait, timeout)
}
