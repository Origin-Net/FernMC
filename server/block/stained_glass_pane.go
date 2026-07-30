package block

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/world"
)


type StainedGlassPane struct {
	transparent
	thin
	clicksAndSticks
	sourceWaterDisplacer

	
	Colour item.Colour
}


func (p StainedGlassPane) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}


func (p StainedGlassPane) BreakInfo() BreakInfo {
	return newBreakInfo(0.3, alwaysHarvestable, nothingEffective, silkTouchOnlyDrop(p))
}


func (p StainedGlassPane) EncodeItem() (name string, meta int16) {
	return "minecraft:" + p.Colour.String() + "_stained_glass_pane", 0
}


func (p StainedGlassPane) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + p.Colour.String() + "_stained_glass_pane", nil
}


func allStainedGlassPane() []world.Block {
	b := make([]world.Block, 0, 16)
	for _, c := range item.Colours() {
		b = append(b, StainedGlassPane{Colour: c})
	}
	return b
}
