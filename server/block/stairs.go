package block

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/block/model"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/Origin-Net/FernMC/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)


type Stairs struct {
	transparent
	sourceWaterDisplacer

	
	Block world.Block
	
	
	UpsideDown bool
	
	Facing cube.Direction
}



func (s Stairs) UseOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) (used bool) {
	pos, face, used = firstReplaceable(tx, pos, face, s)
	if !used {
		return
	}
	s.Facing = user.Rotation().Direction()
	if face == cube.FaceDown || (clickPos[1] > 0.5 && face != cube.FaceUp) {
		s.UpsideDown = true
	}

	place(tx, pos, s, user, ctx)
	return placed(ctx)
}


func (s Stairs) Model() world.BlockModel {
	return model.Stair{Facing: s.Facing, UpsideDown: s.UpsideDown}
}


func (s Stairs) BreakInfo() BreakInfo {
	breakInfo := s.Block.(Breakable).BreakInfo()
	return newBreakInfo(breakInfo.Hardness, breakInfo.Harvestable, breakInfo.Effective, oneOf(s)).withBlastResistance(breakInfo.BlastResistance)
}


func (s Stairs) Instrument() sound.Instrument {
	if _, ok := s.Block.(Planks); ok {
		return sound.Bass()
	}
	if _, ok := s.Block.(BambooMosaic); ok {
		return sound.Bass()
	}
	return sound.BassDrum()
}


func (s Stairs) FlammabilityInfo() FlammabilityInfo {
	if flammable, ok := s.Block.(Flammable); ok {
		return flammable.FlammabilityInfo()
	}
	return newFlammabilityInfo(0, 0, false)
}


func (s Stairs) FuelInfo() item.FuelInfo {
	if fuel, ok := s.Block.(item.Fuel); ok {
		return fuel.FuelInfo()
	}
	return item.FuelInfo{}
}


func (s Stairs) EncodeItem() (name string, meta int16) {
	return "minecraft:" + encodeStairsBlock(s.Block) + "_stairs", 0
}


func (s Stairs) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + encodeStairsBlock(s.Block) + "_stairs", map[string]any{"upside_down_bit": s.UpsideDown, "weirdo_direction": toStairsDirection(s.Facing)}
}


func toStairsDirection(v cube.Direction) int32 {
	return int32(3 - v)
}


func (s Stairs) SideClosed(pos, side cube.Pos, tx *world.Tx) bool {
	return s.Model().FaceSolid(pos, pos.Face(side), tx)
}


func (Stairs) CanRedstoneWireStepDown(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}


func allStairs() (stairs []world.Block) {
	f := func(facing cube.Direction, upsideDown bool) {
		for _, s := range StairsBlocks() {
			stairs = append(stairs, Stairs{Facing: facing, UpsideDown: upsideDown, Block: s})
		}
	}
	for i := cube.Direction(0); i <= 3; i++ {
		f(i, true)
		f(i, false)
	}
	return
}
