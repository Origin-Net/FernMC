package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/world"
)

type pluginHelp struct {
	Sub cmd.SubCommand `cmd:"help"`
}

func (pluginHelp) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("§6/plugin help§r - Show this help")
	o.Print("§6/plugin list§r - List loaded plugins")
	o.Print("§6/plugin load <name>§r - Load a plugin from plugins/<name>.pl")
	o.Print("§6/plugin reload <name>§r - Reload a loaded plugin")
	o.Print("§6/plugin unload <name>§r - Unload a loaded plugin")
}

type pluginList struct {
	Sub cmd.SubCommand `cmd:"list"`
}

func (pluginList) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	all := AllPlugins()
	if len(all) == 0 {
		o.Print("§7No plugins loaded.")
		return
	}
	o.Print(fmt.Sprintf("§6Plugins (%d):§r", len(all)))
	for _, p := range all {
		o.Print(fmt.Sprintf("  §a%s§r v%s", p.Meta().Name, p.Meta().Version))
	}
}

type pluginLoad struct {
	Sub  cmd.SubCommand `cmd:"load"`
	Name string
}

func (p pluginLoad) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	name := strings.ToLower(p.Name)
	if !safePluginName(name) {
		o.Error("Invalid plugin name.")
		return
	}
	path := func() string {
		for _, ext := range []string{".fpl", ".so", ".pl"} {
			p := "plugins/" + name + ext
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
		return "plugins/" + name + ".fpl"
	}()
	plug, err := LoadPlugin(path)
	if err != nil {
		o.Error("Failed to load plugin: " + err.Error())
		return
	}
	o.Print(fmt.Sprintf("Loaded plugin §a%s§r v%s", plug.Meta().Name, plug.Meta().Version))
	if srv := GetServer(); srv != nil {
		srv.RefreshPlayersCommands()
	}
}

type pluginReload struct {
	Sub  cmd.SubCommand `cmd:"reload"`
	Name string
}

func (p pluginReload) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	name := strings.ToLower(p.Name)
	if !safePluginName(name) {
		o.Error("Invalid plugin name.")
		return
	}
	srv := GetServer()
	if srv == nil {
		o.Error("Server not available for reload.")
		return
	}
	newPlug, err := ReloadPlugin(name, srv)
	if err != nil {
		o.Error(fmt.Sprintf("Failed to reload plugin §e%s§r: %v", name, err))
		return
	}
	o.Print(fmt.Sprintf("Reloaded plugin §a%s§r v%s", newPlug.Meta().Name, newPlug.Meta().Version))
	srv.RefreshPlayersCommands()
}

type pluginUnload struct {
	Sub  cmd.SubCommand `cmd:"unload"`
	Name string
}

func (p pluginUnload) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	name := strings.ToLower(p.Name)
	if !safePluginName(name) {
		o.Error("Invalid plugin name.")
		return
	}
	plug, ok := UnloadPlugin(name)
	if !ok {
		o.Error(fmt.Sprintf("Plugin §e%s§r is not loaded.", name))
		return
	}
	o.Print(fmt.Sprintf("Unloaded plugin §a%s§r v%s", plug.Meta().Name, plug.Meta().Version))
	if srv := GetServer(); srv != nil {
		srv.RefreshPlayersCommands()
	}
}

func safePluginName(name string) bool {
	if len(name) == 0 || len(name) > 64 {
		return false
	}
	for _, r := range name {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '_' || r == '-' {
			continue
		}
		return false
	}
	return true
}

func init() {
	cmd.Register(cmd.New("plugin", "Manage plugins", []string{"pl"},
		pluginHelp{},
		pluginList{},
		pluginLoad{},
		pluginReload{},
		pluginUnload{},
	))
}
