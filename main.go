package main

import (
	"fmt"
	"os"

	"github.com/devo/dj/cmds"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Docker Jockey"
	app.Usage = "Spinning sick containers"
	app.Version = "0.0.5"
	app.Action = cli.ShowAppHelp

	app.Commands = []cli.Command{
		cmds.InstallCmd(),
		cmds.RunCmd(),
		cmds.UninstallCmd(),
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}
