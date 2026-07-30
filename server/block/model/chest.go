package model

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/world"
)



type Chest struct{}


func (Chest) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{cube.Box(0.025, 0, 0.025, 0.975, 0.95, 0.975)}
}


func (Chest) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
