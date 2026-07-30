package skywars

import (
	"time"

	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/go-gl/mathgl/mgl64"
)


type handler struct {
	player.NopHandler
	match *Match
}



func (h *handler) HandleHurt(ctx *player.Context, damage *float64, _ bool, _ *time.Duration, _ world.DamageSource) {
	p := ctx.Player()
	if !h.match.pvpEnabled {
		ctx.Cancel()
		return
	}
	if p.Health()-*damage <= mgl64.Epsilon {
		ctx.Cancel()
		h.match.handleDeath(p)
	}
}


func (h *handler) HandleQuit(p *player.Player) {
	h.match.handleQuit(p)
}


func (h *handler) HandleBlockBreak(ctx *player.Context, _ cube.Pos, _ *[]item.Stack, _ *int) {
	if !h.match.pvpEnabled {
		ctx.Cancel()
	}
}


func (h *handler) HandleBlockPlace(ctx *player.Context, _ cube.Pos, _ world.Block) {
	if !h.match.pvpEnabled {
		ctx.Cancel()
	}
}


func (h *handler) HandleAttackEntity(ctx *player.Context, _ world.Entity, _, _ *float64, _ *bool) {
	if !h.match.pvpEnabled {
		ctx.Cancel()
	}
}


func (h *handler) HandleItemUse(ctx *player.Context) {
	p := ctx.Player()
	if p.GameMode().AllowsInteraction() {
		return 
	}
	held, _ := p.HeldItems()
	if _, ok := held.Item().(item.Compass); !ok {
		return
	}
	ctx.Cancel()
	h.match.openDeathMenu(p)
}


func (h *handler) HandleItemUseOnBlock(ctx *player.Context, _ cube.Pos, _ cube.Face, _ mgl64.Vec3) {
	if !ctx.Player().GameMode().AllowsInteraction() {
		ctx.Cancel()
	}
}


func (h *handler) HandleItemUseOnEntity(ctx *player.Context, _ world.Entity) {
	if !ctx.Player().GameMode().AllowsInteraction() {
		ctx.Cancel()
	}
}
