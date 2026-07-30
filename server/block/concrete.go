package block

import (
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/world"
)



type Concrete struct {
	solid
	bassDrum

	
	Colour item.Colour
}


func (c Concrete) BreakInfo() BreakInfo {
	return newBreakInfo(1.8, pickaxeHarvestable, pickaxeEffective, oneOf(c))
}


func (c Concrete) EncodeItem() (name string, meta int16) {
	return "minecraft:" + c.Colour.String() + "_concrete", 0
}


func (c Concrete) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + c.Colour.String() + "_concrete", nil
}


func allConcrete() []world.Block {
	b := make([]world.Block, 0, 16)
	for _, c := range item.Colours() {
		b = append(b, Concrete{Colour: c})
	}
	return b
}
