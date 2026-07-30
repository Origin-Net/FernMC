package server

import (
	"github.com/Origin-Net/FernMC/server/entity"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/item/creative"
)

func init() {
	spawnEggGroup := creative.Group{
		Category: creative.ItemsCategory(),
		Name:     "spawn_eggs",
		Icon:     item.NewStack(item.SpawnEgg{ItemName: "minecraft:zombie_spawn_egg"}, 1).WithValue("mob_type", "zombie"),
	}
	creative.RegisterGroup(spawnEggGroup)
	for _, t := range entity.MobTypes() {
		mobName := t.EncodeEntity()[10:]
		if len(mobName) > 3 && mobName[len(mobName)-3:] == "_v2" {
			mobName = mobName[:len(mobName)-3]
		}
		itemName := "minecraft:" + mobName + "_spawn_egg"
		st := item.NewStack(item.SpawnEgg{ItemName: itemName}, 1).WithValue("mob_type", mobName)
		creative.RegisterItem(creative.Item{Stack: st, Group: spawnEggGroup.Name})
	}
}
