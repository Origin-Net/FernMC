package model

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/world"
)


type Carpet struct{}


func (Carpet) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{cube.Box(0, 0, 0, 1, 0.0625, 1)}
}


func (Carpet) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
