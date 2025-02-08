package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/ttak0422/hostage/commands"
)

const version = "0.1.0"

var revision = "HEAD"

func main() {
	cmd := &cli.Command{
		Name:    "hostage",
		Usage:   "hosts file manager for cli",
		Version: fmt.Sprintf("%s (rev:%s)", version, revision),
		Commands: []*cli.Command{
			{
				Name:   "setup",
				Usage:  "Create a new hostage config file (if not exists)",
				Action: commands.SetupConfig,
			},
			{
				Name:   "keys",
				Usage:  "Get hosts group keys",
				Action: commands.ShowConfigKey,
			},
			{
				Name:   "get",
				Usage:  "Get hosts group config",
				Action: commands.ShowConfig,
			},
			{
				Name:   "edit",
				Usage:  "edit config file",
				Action: commands.OpenWithEditor,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
