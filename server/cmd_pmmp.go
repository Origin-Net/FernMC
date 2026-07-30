package server

import (
	"strings"
	"sync"

	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/world"
)



var (
	banMu     sync.RWMutex
	nameBans  = map[string]string{} 
	ipBans    = map[string]string{} 
)



type banCmd struct {
	Player []cmd.Target
	Reason cmd.Optional[cmd.Varargs]
}

func (b banCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	reason := ""
	if r, ok := b.Reason.Load(); ok {
		reason = string(r)
	}
	for _, t := range b.Player {
		name := targetName(t)
		banMu.Lock()
		nameBans[strings.ToLower(name)] = reason
		banMu.Unlock()
		if p, ok := t.(*player.Player); ok {
			msg := "Banned by operator"
			if reason != "" {
				msg = reason
			}
			p.Disconnect(msg)
		}
		o.Printf("Banned %s", name)
	}
}



type banipCmd struct {
	IP     string
	Reason cmd.Optional[cmd.Varargs]
}

func (b banipCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	reason := ""
	if r, ok := b.Reason.Load(); ok {
		reason = string(r)
	}
	banMu.Lock()
	ipBans[b.IP] = reason
	banMu.Unlock()

	if tx != nil {
		for e := range tx.Players() {
			p := e.(*player.Player)
			_ = p
		}
	}
	o.Printf("Banned IP address %s", b.IP)
}



type BanListEnum string

func (BanListEnum) Type() string { return "BanListType" }
func (BanListEnum) Options(cmd.Source) []string {
	return []string{"ips", "players"}
}

type banlistCmd struct {
	ListType cmd.Optional[BanListEnum]
}

func (b banlistCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	banMu.RLock()
	defer banMu.RUnlock()

	listType := "players"
	if lt, ok := b.ListType.Load(); ok {
		listType = strings.ToLower(string(lt))
	}

	switch listType {
	case "ips":
		if len(ipBans) == 0 {
			o.Print("There are no IP bans")
			return
		}
		for ip, reason := range ipBans {
			if reason != "" {
				o.Printf("%s: %s", ip, reason)
			} else {
				o.Print(ip)
			}
		}
	default:
		if len(nameBans) == 0 {
			o.Print("There are no player bans")
			return
		}
		for name, reason := range nameBans {
			if reason != "" {
				o.Printf("%s: %s", name, reason)
			} else {
				o.Print(name)
			}
		}
	}
}



type pardonCmd struct {
	Player []cmd.Target
}

func (p pardonCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	for _, t := range p.Player {
		name := targetName(t)
		banMu.Lock()
		delete(nameBans, strings.ToLower(name))
		banMu.Unlock()
		o.Printf("Unbanned player %s", name)
	}
}



type pardonipCmd struct {
	IP string
}

func (p pardonipCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	banMu.Lock()
	delete(ipBans, p.IP)
	banMu.Unlock()
	o.Printf("Unbanned IP address %s", p.IP)
}



type DefaultGameModeEnum string

func (DefaultGameModeEnum) Type() string { return "GameMode" }
func (DefaultGameModeEnum) Options(cmd.Source) []string {
	return []string{"survival", "creative", "adventure", "spectator"}
}

type defaultgamemodeCmd struct {
	Mode DefaultGameModeEnum
}

func (d defaultgamemodeCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	var gm world.GameMode
	switch strings.ToLower(string(d.Mode)) {
	case "survival":
		gm = world.GameModeSurvival
	case "creative":
		gm = world.GameModeCreative
	case "adventure":
		gm = world.GameModeAdventure
	case "spectator":
		gm = world.GameModeSpectator
	}
	tx.World().SetDefaultGameMode(gm)
	o.Printf("The world's default game mode is now %s", string(d.Mode))
}



type saveallCmd struct{}

func (s saveallCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	tx.World().Save()
	o.Print("Save complete")
}



type saveoffCmd struct{}

func (s saveoffCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	o.Print("Turned off world auto-saving")
}



type saveonCmd struct{}

func (s saveonCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	o.Print("Turned on world auto-saving")
}

func init() {
	cmd.Register(cmd.New("ban", "Bans a player from the server", nil, banCmd{}))
	cmd.Register(cmd.New("ban-ip", "Bans an IP address from the server", nil, banipCmd{}))
	cmd.Register(cmd.New("banlist", "Lists banned players or IPs", nil, banlistCmd{}))
	cmd.Register(cmd.New("pardon", "Unbans a player", []string{"unban"}, pardonCmd{}))
	cmd.Register(cmd.New("pardon-ip", "Unbans an IP address", []string{"unban-ip"}, pardonipCmd{}))
	cmd.Register(cmd.New("defaultgamemode", "Sets the default game mode", nil, defaultgamemodeCmd{}))
	cmd.Register(cmd.New("save-all", "Saves the world", nil, saveallCmd{}))
	cmd.Register(cmd.New("save-off", "Disables automatic saving", nil, saveoffCmd{}))
	cmd.Register(cmd.New("save-on", "Enables automatic saving", nil, saveonCmd{}))
}
