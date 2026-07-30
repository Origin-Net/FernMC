package model

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/world"
)


type Cactus struct{}


func (Cactus) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{cube.Box(0.025, 0, 0.025, 0.975, 1, 0.975)}
}


func (Cactus) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
