package cmds

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
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
			cli.BoolFlag{
				Name:  "bin",
				Usage: "Installs any ol' binary onto your system. Just use this flag and the URL to the binary. No more lame bash install scripts.",
			},
			cli.StringFlag{
				Name:  "to",
				Usage: "Write the bin to this directory.",
			},
		},
		Action: func(c *cli.Context) error {
			image := c.Args().First()
			if image == "" {
				return errors.New("Must provide image name to install")
			}
			var contents io.Reader
			name := ""
			if c.Bool("bin") {
				name = path.Base(image)
				binURL := image
				// Check for {LATEST} and swap it in with latest release
				// https://github.com/devo/dj/releases/download/{LATEST}/dj_mac
				// https://api.github.com/repos/fnproject/cli/releases/latest
				if strings.Contains(binURL, "{LATEST}") {
					split := strings.Split(binURL, "/")
					si := 0
					for i, v := range split {
						// fmt.Println(i, "=", v)
						if v == "github.com" {
							si = i
							break
						}
					}
					url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", split[si+1], split[si+2])
					// fmt.Println("release url", url)
					r, err := http.Get(url)
					if err != nil {
						return err
					}
					defer r.Body.Close()
					if r.StatusCode >= 400 {
						return fmt.Errorf("Could not get bin from %v", url)
					}
					ghRelease := &GitHubRelease{}
					err = json.NewDecoder(r.Body).Decode(ghRelease)
					if err != nil {
						return err
					}
					binURL = strings.Replace(binURL, "{LATEST}", ghRelease.TagName, 1)
					image = binURL
				}
				// fmt.Println("Getting bin from", binURL)
				response, e := http.Get(binURL)
				if e != nil {
					log.Fatal(e)
				}
				defer response.Body.Close()
				contents = response.Body
			} else {
				tag := ""
				split := strings.Split(image, ":")
				if len(split) == 2 {
					image = split[0]
					tag = split[1]
				}
				if tag != "" {
					//  nothing, just don't want to remove it, will be used later
				}
				name = image
				split = strings.Split(image, "/")
				if len(split) > 2 {
					return fmt.Errorf("Invalid image name to install")
				}
				if len(split) == 2 {
					name = split[1]
				}
				// prepare script template
				tmpl, err := template.New("linkage").Parse(linkTmpl)
				if err != nil {
					panic(err)
				}
				proggie := Proggie{Image: image}
				var b bytes.Buffer
				err = tmpl.Execute(&b, proggie)
				if err != nil {
					panic(err)
				}
				contents = &b
			}
			// user can override the name
			if c.String("name") != "" {
				name = c.String("name")
			}
			dest := linkFileLocation(c, name)
			// if !c.Bool("force") {
			// instead of using force, we'll just do a backup (less annoying). need to add a rollback function to get it back
			if _, err := os.Stat(dest); !os.IsNotExist(err) {
				// return fmt.Errorf("File already exists at %v", dest)
			} else {
				os.Rename(dest, dest+".bak")
			}
			// }
			fmt.Printf("installing %v to %v\n", image, dest)

			f, err := os.Create(dest)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			_, err = io.Copy(f, contents)
			if err != nil {
				log.Fatal(err)
			}

			if err := os.Chmod(dest, 0755); err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}
}

func linkFileLocation(c *cli.Context, name string) string {
	if c.IsSet("to") {
		return fmt.Sprintf("%s/%s", c.String("to"), name)
	}
	switch os := runtime.GOOS; os {
	case "darwin", "linux":
		return "/usr/local/bin/" + name
	case "windows":
		log.Fatalf("Windows not fully supported yet. Where should we install it??")
	default:
		log.Fatalf("Unsupported operating system %v. Please make an issue on GitHub.", os)
	}
	return ""
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

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}
