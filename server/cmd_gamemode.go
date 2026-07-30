package server

import (
	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/world"
)

type GamemodeEnum string

func (GamemodeEnum) Type() string { return "GameMode" }
func (GamemodeEnum) Options(cmd.Source) []string {
	return []string{"survival", "creative", "adventure", "spectator"}
}

type gamemodeCmd struct {
	Mode   GamemodeEnum
	Player cmd.Optional[[]cmd.Target]
}

func (g gamemodeCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	mode := modeFor(string(g.Mode))
	if mode == nil {
		o.Error("Invalid gamemode. Use: survival, creative, adventure, spectator")
		return
	}
	players, explicit := g.Player.Load()
	if explicit && len(players) > 0 {
		for _, t := range players {
			p, ok := t.(*player.Player)
			if !ok {
				o.Error("Target must be a player")
				return
			}
			p.SetGameMode(mode)
			o.Printf("Set %s's game mode to %s", p.Name(), string(g.Mode))
		}
		return
	}
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Only players can change their own game mode without specifying a target")
		return
	}
	p.SetGameMode(mode)
	o.Printf("Set own game mode to %s", string(g.Mode))
}

func modeFor(s string) world.GameMode {
	switch s {
	case "survival":
		return world.GameModeSurvival
	case "creative":
		return world.GameModeCreative
	case "adventure":
		return world.GameModeAdventure
	case "spectator":
		return world.GameModeSpectator
	}
	return nil
}

type Gmc struct{}

func (Gmc) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Only players can use this command")
		return
	}
	p.SetGameMode(world.GameModeCreative)
	o.Print("Set game mode to Creative")
}

type Gms struct{}

func (Gms) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Only players can use this command")
		return
	}
	p.SetGameMode(world.GameModeSurvival)
	o.Print("Set game mode to Survival")
}

type Gma struct{}

func (Gma) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Only players can use this command")
		return
	}
	p.SetGameMode(world.GameModeAdventure)
	o.Print("Set game mode to Adventure")
}

type Gmsp struct{}

func (Gmsp) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Only players can use this command")
		return
	}
	p.SetGameMode(world.GameModeSpectator)
	o.Print("Set game mode to Spectator")
}

func init() {
	cmd.Register(cmd.New("gamemode", "Change game mode", []string{"gm"},
		gamemodeCmd{},
	))
	cmd.Register(cmd.New("gmc", "Set game mode to Creative", nil, Gmc{}))
	cmd.Register(cmd.New("gms", "Set game mode to Survival", nil, Gms{}))
	cmd.Register(cmd.New("gma", "Set game mode to Adventure", nil, Gma{}))
	cmd.Register(cmd.New("gmsp", "Set game mode to Spectator", nil, Gmsp{}))
}
