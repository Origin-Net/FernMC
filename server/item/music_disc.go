package item

import (
	"github.com/Origin-Net/FernMC/server/world/sound"
)


type MusicDisc struct {
	
	DiscType sound.DiscType
}


func (MusicDisc) MaxCount() int {
	return 1
}


func (m MusicDisc) EncodeItem() (name string, meta int16) {
	return "minecraft:music_disc_" + m.DiscType.String(), 0
}
