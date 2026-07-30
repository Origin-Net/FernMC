package item

import (
	"github.com/Origin-Net/FernMC/server/entity/effect"
	"github.com/Origin-Net/FernMC/server/world"
	"time"
)


type GoldenApple struct{}


func (e GoldenApple) AlwaysConsumable() bool {
	return true
}


func (e GoldenApple) ConsumeDuration() time.Duration {
	return DefaultConsumeDuration
}


func (e GoldenApple) Consume(_ *world.Tx, c Consumer) Stack {
	c.Saturate(4, 9.6)
	prev := c.Absorption()
	c.AddEffect(effect.New(effect.Absorption, 1, 2*time.Minute))
	c.SetAbsorption(max(prev, min(prev+4, 16)))
	c.AddEffect(effect.New(effect.Regeneration, 2, 5*time.Second))
	return Stack{}
}


func (e GoldenApple) EncodeItem() (name string, meta int16) {
	return "minecraft:golden_apple", 0
}
