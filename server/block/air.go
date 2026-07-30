package block

import "github.com/Origin-Net/FernMC/server/world"


type Air struct {
	empty
	replaceable
	transparent
}


func (Air) HasLiquidDrops() bool {
	return false
}


func (Air) PortalInterior(target world.Dimension) bool {
	return target == world.Nether
}


func (Air) EncodeItem() (name string, meta int16) {
	return "minecraft:air", 0
}


func (Air) EncodeBlock() (string, map[string]any) {
	return "minecraft:air", nil
}
