package block

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/world"
)


type GlassPane struct {
	transparent
	thin
	clicksAndSticks
	sourceWaterDisplacer
}


func (p GlassPane) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}


func (p GlassPane) BreakInfo() BreakInfo {
	return newBreakInfo(0.3, alwaysHarvestable, nothingEffective, silkTouchOnlyDrop(p))
}


func (GlassPane) EncodeItem() (name string, meta int16) {
	return "minecraft:glass_pane", meta
}


func (GlassPane) EncodeBlock() (string, map[string]any) {
	return "minecraft:glass_pane", nil
}
