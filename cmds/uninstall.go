// TODO:
// * Keep track of what we've installed so we don't uninstall something that wasn't installed via this.

package cmds

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func UninstallCmd() cli.Command {
	return cli.Command{
		Name:      "uninstall",
		Usage:     "remove a program",
		Flags:     []cli.Flag{},
		ArgsUsage: "NAME",
		Action: func(c *cli.Context) error {
			name := c.Args().First()
			if name == "" {
				return errors.New("Must provide name of program to uninstall")
			}
			err := os.Remove(linkFileLocation(name))
			if err != nil {
				return err
			}
			fmt.Printf("%s removed", name)
			return nil
		},
	}
}
