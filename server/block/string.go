package block

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/go-gl/mathgl/mgl64"
)




type String struct {
	empty
	transparent

	
	Attached bool
	
	Disarmed bool
	
	Powered bool
	
	Suspended bool
}


func (s String) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, _, canPlace := firstReplaceable(tx, pos, face, s)
	if !canPlace {
		return false
	}
	below := pos.Side(cube.FaceDown)
	s.Suspended = !tx.Block(below).Model().FaceSolid(below, cube.FaceUp, tx)
	place(tx, pos, s, user, ctx)
	return placed(ctx)
}


func (s String) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	below := pos.Side(cube.FaceDown)
	suspended := !tx.Block(below).Model().FaceSolid(below, cube.FaceUp, tx)
	if suspended != s.Suspended {
		s.Suspended = suspended
		tx.SetBlock(pos, s, nil)
	}
}


func (s String) HasLiquidDrops() bool {
	return true
}


func (s String) BreakInfo() BreakInfo {
	return newBreakInfo(0, alwaysHarvestable, nothingEffective, oneOf(String{}))
}


func (s String) EncodeItem() (name string, meta int16) {
	return "minecraft:string", 0
}


func (s String) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:trip_wire", map[string]any{
		"attached_bit":  boolByte(s.Attached),
		"disarmed_bit":  boolByte(s.Disarmed),
		"powered_bit":   boolByte(s.Powered),
		"suspended_bit": boolByte(s.Suspended),
	}
}


func allString() (blocks []world.Block) {
	for meta := 0; meta < 16; meta++ {
		blocks = append(blocks, String{
			Powered:   meta&0x1 != 0,
			Suspended: meta&0x2 != 0,
			Attached:  meta&0x4 != 0,
			Disarmed:  meta&0x8 != 0,
		})
	}
	return blocks
}
