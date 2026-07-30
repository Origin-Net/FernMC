package block

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/world"
)


type Barrier struct {
	sourceWaterDisplacer
	transparent
	solid
}


func (Barrier) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}


func (Barrier) EncodeItem() (name string, meta int16) {
	return "minecraft:barrier", 0
}


func (Barrier) EncodeBlock() (string, map[string]any) {
	return "minecraft:barrier", nil
}
