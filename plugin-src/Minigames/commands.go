package main

import (
	"fmt"

	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/Origin-Net/FernMC/plugin-src/Minigames/framework"
)


type joinCmd struct {
	fw *framework.Framework
	Game string `cmd:"game"`
}

func (j joinCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Only players can use this command.")
		return
	}

	if err := j.fw.Games().JoinGame(p, j.Game); err != nil {
		o.Error(fmt.Sprintf("§c%s", err.Error()))
		return
	}
	o.Print(fmt.Sprintf("§6[Minigames]§r Joined %s!", j.Game))
}


type leaveCmd struct {
	fw *framework.Framework
}

func (l leaveCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Only players can use this command.")
		return
	}

	match, ok := l.fw.Matches().MatchByPlayer(p)
	if !ok {
		o.Error("§cYou are not in a match.")
		return
	}
	l.fw.Players().ReturnToLobby(p)
	o.Print(fmt.Sprintf("§6[Minigames]§r Left %s.", match.GameID()))
}


type gamesCmd struct {
	fw *framework.Framework
}

func (g gamesCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("§6Available Minigames:§r")
	for _, game := range g.fw.Games().Registry().All() {
		o.Print(fmt.Sprintf("  §a/join %s§r - %s", game.ID(), game.Description()))
	}
}


type setSpawnCmd struct {
	fw *framework.Framework
}

func (s setSpawnCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Error("Only players can use this command.")
		return
	}
	pos := p.Position()
	yaw, pitch := p.Rotation().Elem()
	s.fw.Players().SetSpawn(pos, yaw, pitch)
	o.Print(fmt.Sprintf("§a[Minigames]§r Lobby spawn set to (%.1f, %.1f, %.1f) with yaw=%.1f pitch=%.1f", pos.X(), pos.Y(), pos.Z(), yaw, pitch))
}
