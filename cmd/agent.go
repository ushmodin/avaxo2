package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// CmdAgent CLI Command run avaxo agent
var CmdAgent = &cli.Command{
	Name:   "agent",
	Usage:  "run agent",
	Action: runAgent,
}

func runAgent(ctx *cli.Context) error {
	fmt.Println("Agent work")
	return nil
}
