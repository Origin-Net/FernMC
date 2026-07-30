package block

import (
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/world/sound"
)


type Emerald struct {
	solid
}


func (e Emerald) Instrument() sound.Instrument {
	return sound.Bit()
}


func (e Emerald) BreakInfo() BreakInfo {
	return newBreakInfo(5, func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierIron.HarvestLevel
	}, pickaxeEffective, oneOf(e)).withBlastResistance(30)
}


func (Emerald) PowersBeacon() bool {
	return true
}


func (Emerald) EncodeItem() (name string, meta int16) {
	return "minecraft:emerald_block", 0
}


func (Emerald) EncodeBlock() (string, map[string]any) {
	return "minecraft:emerald_block", nil
}
