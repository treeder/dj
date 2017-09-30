package cmds

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func InstallCmd() cli.Command {
	return cli.Command{
		Name:  "install",
		Usage: "install a program",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name",
				Usage: "Alternate name instead of default",
			},
			cli.BoolFlag{
				Name:  "force,f",
				Usage: "Force overwrite",
			},
		},
		Action: func(c *cli.Context) error {
			image := c.Args().First()
			if image == "" {
				return errors.New("Must provide image name to install")
			}
			tag := ""
			split := strings.Split(image, ":")
			if len(split) == 2 {
				image = split[0]
				tag = split[1]
			}
			if tag != "" {
				//  nothing, just don't want to remove it, will be used later
			}
			name := image
			split = strings.Split(image, "/")
			if len(split) > 2 {
				return fmt.Errorf("Invalid image name to install")
			}
			if len(split) == 2 {
				name = split[1]
			}
			// user can override the name
			if c.String("name") != "" {
				name = c.String("name")
			}
			dest := linkFileLocation(name)
			if !c.Bool("force") {
				if _, err := os.Stat(dest); !os.IsNotExist(err) {
					return fmt.Errorf("File already exists at %v", dest)
				}
			}
			fmt.Printf("installing %v to %v\n", image, dest)
			tmpl, err := template.New("linkage").Parse(linkTmpl)
			if err != nil {
				panic(err)
			}
			f, err := os.Create(dest)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			proggie := Proggie{Image: image}
			err = tmpl.Execute(f, proggie)
			if err != nil {
				panic(err)
			}
			if err := os.Chmod(dest, 0700); err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}
}

func linkFileLocation(name string) string {
	return "/usr/local/bin/" + name
}

type Proggie struct {
	Image string
}

const (
	linkTmpl = `
x="$@"
dj run {{ .Image }} "$x"
`
)
