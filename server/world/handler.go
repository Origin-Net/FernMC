package world

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
)




type Handler interface {
	
	
	
	
	
	HandleLiquidFlow(ctx *Context, from, into cube.Pos, liquid Liquid, replaced Block)
	
	
	
	
	
	HandleLiquidDecay(ctx *Context, pos cube.Pos, before, after Liquid)
	
	
	
	
	HandleLiquidHarden(ctx *Context, hardenedPos cube.Pos, liquidHardened, otherLiquid, newBlock Block)
	
	
	
	HandleSound(ctx *Context, s Sound, pos mgl64.Vec3)
	
	
	
	
	
	
	HandleFireSpread(ctx *Context, from, to cube.Pos)
	
	
	
	
	
	HandleBlockBurn(ctx *Context, pos cube.Pos)
	
	HandleCropTrample(ctx *Context, pos cube.Pos)
	
	
	
	HandleLeavesDecay(ctx *Context, pos cube.Pos)
	
	
	HandleEntitySpawn(tx *Tx, e Entity)
	
	
	HandleEntityDespawn(tx *Tx, e Entity)
	
	
	
	
	HandleExplosion(ctx *Context, position mgl64.Vec3, entities *[]Entity, blocks *[]cube.Pos, itemDropChance *float64, spawnFire *bool)
	
	
	HandleRedstoneUpdate(ctx *Context, update RedstoneUpdate)
	
	
	
	
	HandleClose(tx *Tx)
}


var _ Handler = (*NopHandler)(nil)




type NopHandler struct{}

func (NopHandler) HandleLiquidFlow(*Context, cube.Pos, cube.Pos, Liquid, Block)                  {}
func (NopHandler) HandleLiquidDecay(*Context, cube.Pos, Liquid, Liquid)                          {}
func (NopHandler) HandleLiquidHarden(*Context, cube.Pos, Block, Block, Block)                    {}
func (NopHandler) HandleSound(*Context, Sound, mgl64.Vec3)                                       {}
func (NopHandler) HandleFireSpread(*Context, cube.Pos, cube.Pos)                                 {}
func (NopHandler) HandleBlockBurn(*Context, cube.Pos)                                            {}
func (NopHandler) HandleCropTrample(*Context, cube.Pos)                                          {}
func (NopHandler) HandleLeavesDecay(*Context, cube.Pos)                                          {}
func (NopHandler) HandleEntitySpawn(*Tx, Entity)                                                 {}
func (NopHandler) HandleEntityDespawn(*Tx, Entity)                                               {}
func (NopHandler) HandleExplosion(*Context, mgl64.Vec3, *[]Entity, *[]cube.Pos, *float64, *bool) {}
func (NopHandler) HandleRedstoneUpdate(*Context, RedstoneUpdate)                                 {}
func (NopHandler) HandleClose(*Tx)                                                               {}
