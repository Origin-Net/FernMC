package item

import (
	"github.com/Origin-Net/FernMC/server/world"
)


type MushroomStew struct {
	defaultFood
}


func (MushroomStew) MaxCount() int {
	return 1
}


func (MushroomStew) Consume(_ *world.Tx, c Consumer) Stack {
	c.Saturate(6, 7.2)
	return NewStack(Bowl{}, 1)
}


func (MushroomStew) EncodeItem() (name string, meta int16) {
	return "minecraft:mushroom_stew", 0
}
