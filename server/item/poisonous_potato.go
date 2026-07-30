package item

import (
	"github.com/Origin-Net/FernMC/server/entity/effect"
	"github.com/Origin-Net/FernMC/server/world"
	"math/rand/v2"
	"time"
)


type PoisonousPotato struct {
	defaultFood
}


func (p PoisonousPotato) Consume(_ *world.Tx, c Consumer) Stack {
	c.Saturate(2, 1.2)
	if rand.Float64() < 0.6 {
		c.AddEffect(effect.New(effect.Poison, 1, 5*time.Second))
	}
	return Stack{}
}


func (p PoisonousPotato) EncodeItem() (name string, meta int16) {
	return "minecraft:poisonous_potato", 0
}
