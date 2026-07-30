package model

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/world"
)


type Portal struct {
	
	Axis cube.Axis
}


func (Portal) BBox(cube.Pos, world.BlockSource) []cube.BBox {
	return nil
}


func (Portal) FaceSolid(cube.Pos, cube.Face, world.BlockSource) bool {
	return false
}
