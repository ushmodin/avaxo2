package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/ushmodin/avaxo2/internal/settings"
)

// CmdAgent CLI Command run avaxo agent
var CmdAgent = &cli.Command{
	Name:   "agent",
	Usage:  "run agent",
	Action: runAgent,
}

func runAgent(ctx *cli.Context) error {
	fmt.Println("Agent work")
	if err := settings.InitSettings(); err != nil {
		return err
	}
	return nil
}
