package block

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/block/model"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/world"
	"time"
)



type WoodFence struct {
	transparent
	bass
	sourceWaterDisplacer

	
	
	Wood WoodType
}


func (w WoodFence) BreakInfo() BreakInfo {
	return newBreakInfo(2, alwaysHarvestable, axeEffective, oneOf(w)).withBlastResistance(15)
}


func (WoodFence) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}


func (w WoodFence) FlammabilityInfo() FlammabilityInfo {
	if !w.Wood.Flammable() {
		return newFlammabilityInfo(0, 0, false)
	}
	return newFlammabilityInfo(5, 20, true)
}


func (w WoodFence) FuelInfo() item.FuelInfo {
	if !w.Wood.Flammable() {
		return item.FuelInfo{}
	}
	return newFuelInfo(time.Second * 15)
}


func (w WoodFence) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:" + w.Wood.String() + "_fence", nil
}


func (w WoodFence) Model() world.BlockModel {
	return model.Fence{Wood: true}
}


func (w WoodFence) EncodeItem() (name string, meta int16) {
	return "minecraft:" + w.Wood.String() + "_fence", 0
}


func allFence() (fence []world.Block) {
	for _, w := range WoodTypes() {
		fence = append(fence, WoodFence{Wood: w})
	}
	return
}
