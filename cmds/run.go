package cmds

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/urfave/cli"
	contextold "golang.org/x/net/context" // need this for docker, they haven't upgraded their libs
)

func RunCmd() cli.Command {
	return cli.Command{
		Name:            "run",
		Usage:           "complete a task on the list",
		SkipFlagParsing: true,
		Action: func(c *cli.Context) error {
			ctx := contextold.Background()
			if len(c.Args()) == 0 {
				fmt.Println("Invalid args")
				os.Exit(1)
			}
			c1 := c.Args().First()
			fmt.Printf("RUNNING! -%v-\n", c1)

			foo2 := map[string]*CommandMeta{} // make(Commands)
			getJson("https://raw.githubusercontent.com/treeder/dj/784bdf804a330f7a98433d4140ab236d7ba0c4a4/commands.json", &foo2)
			fmt.Printf("COMMANDS: %+v\n", foo2)

			command := foo2[c1]
			if command == nil {
				fmt.Println("Unknown command")
				os.Exit(1)
			}

			image := command.Image
			fmt.Println("Image", image)
			image = "docker.io/" + image

			cli, err := client.NewEnvClient()
			if err != nil {
				panic(err)
			}
			cmd := c.Args().Tail()
			if len(cmd) == 0 {

			} else {
				fmt.Println("CMD:", cmd, cmd[0])
				if len(cmd) > 0 {
					cmd = strings.Fields(cmd[0])
				}
				fmt.Println("CMD:", cmd)
			}

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
					Stdin:  true,
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
	}
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
