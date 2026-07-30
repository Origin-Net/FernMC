package model

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/world"
)


type Cake struct {
	
	
	Bites int
}


func (c Cake) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return []cube.BBox{cube.Box(0.0625, 0, 0.0625, 0.9375, 0.5, 0.9375).
		ExtendTowards(cube.FaceWest, -(float64(c.Bites) / 8))}
}


func (c Cake) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
