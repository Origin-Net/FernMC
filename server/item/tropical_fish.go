package item

import "github.com/Origin-Net/FernMC/server/world"


type TropicalFish struct {
	defaultFood
}


func (TropicalFish) Consume(_ *world.Tx, c Consumer) Stack {
	c.Saturate(1, 0.2)
	return Stack{}
}


func (TropicalFish) EncodeItem() (name string, meta int16) {
	return "minecraft:tropical_fish", 0
}
