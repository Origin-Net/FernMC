package item

import (
	"github.com/Origin-Net/FernMC/server/entity/effect"
	"github.com/Origin-Net/FernMC/server/world"
	"time"
)


type SpiderEye struct {
	defaultFood
}


func (SpiderEye) Consume(_ *world.Tx, c Consumer) Stack {
	c.Saturate(2, 3.2)
	c.AddEffect(effect.New(effect.Poison, 1, time.Second*5))
	return Stack{}
}


func (SpiderEye) EncodeItem() (name string, meta int16) {
	return "minecraft:spider_eye", 0
}
