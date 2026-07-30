package item

import (
	"github.com/Origin-Net/FernMC/server/entity/effect"
	"github.com/Origin-Net/FernMC/server/world"
	"math/rand/v2"
	"time"
)


type RottenFlesh struct {
	defaultFood
}


func (RottenFlesh) Consume(_ *world.Tx, c Consumer) Stack {
	c.Saturate(4, 0.8)
	if rand.Float64() < 0.8 {
		c.AddEffect(effect.New(effect.Hunger, 1, 30*time.Second))
	}
	return Stack{}
}


func (RottenFlesh) EncodeItem() (name string, meta int16) {
	return "minecraft:rotten_flesh", 0
}
