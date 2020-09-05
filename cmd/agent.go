package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/ushmodin/avaxo2/internal/agent"
	"github.com/ushmodin/avaxo2/internal/settings"
)

var (
	listen string
	port   int
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
	server, err := agent.NewServer(settings.AgentSettings.Listen, settings.AgentSettings.Keyfile, settings.AgentSettings.Certfile)
	if err != nil {
		return err
	}
	return server.Run()
}
