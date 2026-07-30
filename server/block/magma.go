package block

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/item/enchantment"
	"github.com/Origin-Net/FernMC/server/world"
)


type Magma struct {
	solid
	bassDrum
}


func (Magma) LightEmissionLevel() uint8 {
	return 3
}


func (Magma) EntityStepOn(_ cube.Pos, _ *world.Tx, e world.Entity) {
	if fireProof, ok := e.(interface{ FireProof() bool }); ok && fireProof.FireProof() {
		return
	}
	
	if sneaking, ok := e.(interface{ Sneaking() bool }); ok && sneaking.Sneaking() {
		return
	}
	if l, ok := e.(livingEntity); ok {
		l.Hurt(1, MagmaDamageSource{})
	}
}


func (m Magma) BreakInfo() BreakInfo {
	return newBreakInfo(0.5, pickaxeHarvestable, pickaxeEffective, oneOf(m)).withBlastResistance(30)
}


func (Magma) EncodeItem() (name string, meta int16) {
	return "minecraft:magma", 0
}


func (Magma) EncodeBlock() (string, map[string]any) {
	return "minecraft:magma", nil
}


type MagmaDamageSource struct{}

func (MagmaDamageSource) ReducedByResistance() bool { return true }
func (MagmaDamageSource) ReducedByArmour() bool     { return true }
func (MagmaDamageSource) Fire() bool                { return true }
func (MagmaDamageSource) AffectedByEnchantment(e item.EnchantmentType) bool {
	return e == enchantment.FireProtection
}
func (MagmaDamageSource) IgnoreTotem() bool { return false }
