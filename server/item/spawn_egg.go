package item

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type SpawnEgg struct {
	ItemName string
}

func (e SpawnEgg) EncodeItem() (name string, meta int16) {
	return e.ItemName, 0
}

func (e SpawnEgg) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user User, ctx *UseContext) bool {
	mainHand, _ := user.HeldItems()
	v, ok := mainHand.Value("mob_type")
	if !ok {
		return false
	}
	name, ok := v.(string)
	if !ok {
		return false
	}
	handler, ok := world.SpawnEggHandlerByName(name)
	if !ok {
		return false
	}
	if handler(tx, pos.Side(face)) {
		ctx.SubtractFromCount(1)
		return true
	}
	return false
}
