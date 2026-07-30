package enchantment

import (
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/world"
)



var SwiftSneak swiftSneak

type swiftSneak struct{}


func (swiftSneak) Name() string {
	return "Swift Sneak"
}


func (swiftSneak) MaxLevel() int {
	return 3
}


func (swiftSneak) Cost(level int) (int, int) {
	minCost := level * 25
	return minCost, minCost + 50
}


func (swiftSneak) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityVeryRare
}


func (swiftSneak) CompatibleWithEnchantment(item.EnchantmentType) bool {
	return true
}


func (swiftSneak) Treasure() bool {
	return true
}


func (swiftSneak) CompatibleWithItem(i world.Item) bool {
	b, ok := i.(item.LeggingsType)
	return ok && b.Leggings()
}
