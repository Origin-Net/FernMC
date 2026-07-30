package enchantment

import (
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/world"
)


var Protection protection

type protection struct{}


func (protection) Name() string {
	return "Protection"
}


func (protection) MaxLevel() int {
	return 4
}


func (protection) Cost(level int) (int, int) {
	minCost := 1 + (level-1)*11
	return minCost, minCost + 11
}


func (protection) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityCommon
}


func (protection) Modifier() float64 {
	return 0.04
}


func (protection) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return t != BlastProtection && t != FireProtection && t != ProjectileProtection
}


func (protection) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Armour)
	return ok
}
