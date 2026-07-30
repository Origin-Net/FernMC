package model

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/world"
)


type Stonecutter struct{}


func (Stonecutter) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{cube.Box(0, 0, 0, 1, 0.5625, 1)}
}


func (Stonecutter) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
