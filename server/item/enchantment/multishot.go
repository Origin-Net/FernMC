package enchantment

import (
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/world"
)


var Multishot multishot

type multishot struct{}


func (multishot) Name() string {
	return "Multishot"
}


func (multishot) MaxLevel() int {
	return 1
}


func (m multishot) Cost(level int) (int, int) {
	return 20, 50
}


func (multishot) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityRare
}


func (multishot) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return t != Piercing
}


func (multishot) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Crossbow)
	return ok
}


func (multishot) MultipleProjectiles() bool {
	return true
}
