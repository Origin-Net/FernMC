package item

import (
	"github.com/Origin-Net/FernMC/server/item/potion"
	"github.com/Origin-Net/FernMC/server/world"
	"time"
)


type Potion struct {
	
	Type potion.Potion
}


func (p Potion) MaxCount() int {
	return 1
}


func (p Potion) AlwaysConsumable() bool {
	return true
}


func (p Potion) ConsumeDuration() time.Duration {
	return DefaultConsumeDuration
}


func (p Potion) Consume(_ *world.Tx, c Consumer) Stack {
	for _, effect := range p.Type.Effects() {
		c.AddEffect(effect)
	}
	return NewStack(GlassBottle{}, 1)
}


func (p Potion) EncodeItem() (name string, meta int16) {
	return "minecraft:potion", int16(p.Type.Uint8())
}
