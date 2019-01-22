package cmds

type Commands map[string]*CommandMeta

type CommandMeta struct {
	Image string `json:"image"`
}
