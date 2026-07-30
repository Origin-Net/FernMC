package block

import (
	"github.com/Origin-Net/FernMC/server/item"
)


type RawCopper struct {
	solid
	bassDrum
}


func (r RawCopper) BreakInfo() BreakInfo {
	return newBreakInfo(5, func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierStone.HarvestLevel
	}, pickaxeEffective, oneOf(r)).withBlastResistance(30)
}


func (RawCopper) EncodeItem() (name string, meta int16) {
	return "minecraft:raw_copper_block", 0
}


func (RawCopper) EncodeBlock() (string, map[string]any) {
	return "minecraft:raw_copper_block", nil
}
