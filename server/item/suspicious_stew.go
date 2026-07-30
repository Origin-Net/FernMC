package item

import (
	"github.com/Origin-Net/FernMC/server/world"
)


type SuspiciousStew struct {
	defaultFood

	
	Type StewType
}


func (SuspiciousStew) MaxCount() int {
	return 1
}


func (SuspiciousStew) AlwaysConsumable() bool {
	return true
}


func (s SuspiciousStew) EncodeItem() (name string, meta int16) {
	return "minecraft:suspicious_stew", int16(s.Type.Uint8())
}


func (s SuspiciousStew) Consume(_ *world.Tx, c Consumer) Stack {
	for _, effect := range s.Type.Effects() {
		c.AddEffect(effect)
	}
	c.Saturate(6, 7.2)

	return NewStack(Bowl{}, 1)
}
