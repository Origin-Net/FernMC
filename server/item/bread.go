package item

import "github.com/Origin-Net/FernMC/server/world"


type Bread struct {
	defaultFood
}


func (Bread) Consume(_ *world.Tx, c Consumer) Stack {
	c.Saturate(5, 6)
	return Stack{}
}


func (Bread) CompostChance() float64 {
	return 0.85
}


func (Bread) EncodeItem() (name string, meta int16) {
	return "minecraft:bread", 0
}
