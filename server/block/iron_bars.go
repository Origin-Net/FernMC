package block

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/world"
)


type IronBars struct {
	transparent
	thin
	sourceWaterDisplacer
}


func (i IronBars) BreakInfo() BreakInfo {
	return newBreakInfo(5, pickaxeHarvestable, pickaxeEffective, oneOf(i)).withBlastResistance(30)
}


func (i IronBars) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}


func (IronBars) EncodeItem() (name string, meta int16) {
	return "minecraft:iron_bars", 0
}


func (IronBars) EncodeBlock() (string, map[string]any) {
	return "minecraft:iron_bars", nil
}
