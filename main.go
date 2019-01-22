package main

import (
	"fmt"
	"os"

	"github.com/treeder/dj/cmds"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Docker Jockey"
	app.Usage = "Spinning sick containers"
	app.Version = "0.1.0"
	rcmd := cmds.RunCmd()
	app.Action = rcmd.Action

	app.Commands = []cli.Command{
		cmds.InstallCmd(),
		rcmd,
		cmds.UninstallCmd(),
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}
