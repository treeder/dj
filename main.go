package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/urfave/cli"
	contextold "golang.org/x/net/context"
)

func main() {
	app := cli.NewApp()
	app.Name = "Docker Jockey"
	app.Usage = "Spinning sick containers"
	app.Version = "0.0.1"
	app.Action = cli.ShowAppHelp

	app.Commands = []cli.Command{
		{
			Name:  "install",
			Usage: "add a task to the list",
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
				dest := "/usr/local/bin/" + name
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
		},
		{
			Name:            "run",
			Usage:           "complete a task on the list",
			SkipFlagParsing: true,
			Action: func(c *cli.Context) error {
				ctx := contextold.Background()
				image := c.Args().First()
				// fmt.Printf("RUNNING! -%v-", image)
				cli, err := client.NewEnvClient()
				if err != nil {
					panic(err)
				}
				cmd := c.Args().Tail()
				// fmt.Println("CMD:", cmd, cmd[0])
				if len(cmd) > 0 {
					cmd = strings.Fields(cmd[0])
				}
				// fmt.Println("CMD:", cmd)

				// see if we already have image, if not, pull it
				_, _, err = cli.ImageInspectWithRaw(ctx, image)
				if client.IsErrNotFound(err) {
					out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
					if err != nil {
						panic(err)
					}
					io.Copy(os.Stdout, out)
				}
				wd, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				mounts := []string{fmt.Sprintf("%s:%s", wd, "/wd")}
				cfg := &container.Config{
					Image:        image,
					AttachStdout: true,
					AttachStderr: true,
					OpenStdin:    true,
					AttachStdin:  true,
					Tty:          true,
					// Volumes:      mounts, // List of volumes (mounts) used for the container
					WorkingDir: "/wd", // Current directory (PWD) in the command will be launched
				}
				// if len(cmd) > 0 {
				cfg.Cmd = cmd
				// }
				resp, err := cli.ContainerCreate(ctx, cfg, &container.HostConfig{
					Binds: mounts,
				}, nil, "")
				if err != nil {
					panic(err)
				}

				go func() {
					resp2, err := cli.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
						Stream: true,
						// Stdin      bool
						Stdout: true,
						Stderr: true,
						// DetachKeys string
						// Logs       bool
					})
					if err != nil {
						panic(err)
					}
					defer resp2.Close()
					io.Copy(os.Stdout, resp2.Reader)
				}()

				if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
					panic(err)
				}

				if _, err = cli.ContainerWait(ctx, resp.ID); err != nil {
					panic(err)
				}

				// out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
				// if err != nil {
				// 	panic(err)
				// }
				// io.Copy(os.Stdout, out)
				return nil
			},
		},
		{
			Name:    "template",
			Aliases: []string{"t"},
			Usage:   "options for task templates",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new template",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
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
