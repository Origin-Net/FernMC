package framework

import (
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/go-gl/mathgl/mgl64"
)


type SpectatorManager struct {
	logger *LogWrapper
}


func NewSpectatorManager(logger *LogWrapper) *SpectatorManager {
	return &SpectatorManager{logger: logger}
}


func (sm *SpectatorManager) EnableSpectator(p *player.Player) {
	p.SetInvisible()
	p.SetGameMode(world.GameModeSpectator)
	p.SetNameTag("§7[SPECTATOR] " + p.Name())
}


func (sm *SpectatorManager) DisableSpectator(p *player.Player) {
	p.SetGameMode(world.GameModeSurvival)
	p.SetNameTag(p.Name())
}


func (sm *SpectatorManager) CycleSpectatorTarget(p *player.Player, alive []*player.Player) {
	for _, target := range alive {
		if target != nil && target != p {
			sm.spectatePlayer(p, target)
			return
		}
	}
}

func (sm *SpectatorManager) spectatePlayer(viewer, target *player.Player) {
	pos := target.Position()
	viewer.Teleport(pos.Add(mgl64.Vec3{0, 5, 0}))
	viewer.Message("Now spectating: " + target.Name())
}
