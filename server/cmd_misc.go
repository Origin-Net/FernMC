package server

import (
	"sort"
	"strings"

	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/world"
)

type helpCmd struct {
	Command cmd.Optional[CmdName]
}

func (h helpCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if name, ok := h.Command.Load(); ok {
		c, ok := cmd.ByAlias(strings.ToLower(string(name)))
		if !ok {
			o.Error("Unknown command: " + name)
			return
		}
		o.Printf("§6%s§r - %s", c.Name(), c.Description())
		o.Printf("Aliases: %s", strings.Join(c.Aliases(), ", "))
		o.Print(c.Usage())
		return
	}
	cmds := cmd.Commands()
	names := make([]string, 0, len(cmds))
	for n := range cmds {
		names = append(names, n)
	}
	sort.Strings(names)
	o.Print("§l--- Showing help page 1 of 1 ---§r")
	for _, n := range names {
		c := cmds[n]
		if c.Name() == n {
			o.Printf("§6/%s§r - %s", c.Name(), c.Description())
		}
	}
}

func init() {
	cmd.Register(cmd.New("help", "Shows command help", []string{"?"}, helpCmd{}))
}
