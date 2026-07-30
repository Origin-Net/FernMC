package block

import (
	"github.com/Origin-Net/FernMC/server/world"
)



type Dirt struct {
	solid

	
	
	Coarse bool
}


func (d Dirt) SoilFor(block world.Block) bool {
	switch block.(type) {
	case ShortGrass, Fern, DoubleTallGrass, DeadBush:
		return !d.Coarse
	case Flower, DoubleFlower, NetherSprouts, PinkPetals, SugarCane, BambooSapling, Bamboo:
		return true
	}
	return false
}


func (d Dirt) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, alwaysHarvestable, shovelEffective, oneOf(d))
}


func (d Dirt) Till() (world.Block, bool) {
	if d.Coarse {
		return Dirt{Coarse: false}, true
	}
	return Farmland{}, true
}


func (d Dirt) Shovel() (world.Block, bool) {
	return DirtPath{}, true
}


func (d Dirt) EncodeItem() (name string, meta int16) {
	if d.Coarse {
		return "minecraft:coarse_dirt", 0
	}
	return "minecraft:dirt", 0
}


func (d Dirt) EncodeBlock() (string, map[string]any) {
	if d.Coarse {
		return "minecraft:coarse_dirt", nil
	}
	return "minecraft:dirt", nil
}


func supportsVegetation(vegetation, block world.Block) bool {
	soil, ok := block.(Soil)
	return ok && soil.SoilFor(vegetation)
}


type Soil interface {
	
	SoilFor(world.Block) bool
}
