package block

import (
	"math/rand/v2"

	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/world"
)


type Grass struct {
	solid
}



var plantSelection = []world.Block{
	Flower{Type: OxeyeDaisy()},
	Flower{Type: PinkTulip()},
	Flower{Type: Cornflower()},
	Flower{Type: WhiteTulip()},
	Flower{Type: RedTulip()},
	Flower{Type: OrangeTulip()},
	Flower{Type: Dandelion()},
	Flower{Type: Poppy()},
}


func init() {
	for i := 0; i < 8; i++ {
		plantSelection = append(plantSelection, Fern{})
	}
	for i := 0; i < 12; i++ {
		plantSelection = append(plantSelection, ShortGrass{})
	}
}


func (g Grass) SoilFor(block world.Block) bool {
	switch block.(type) {
	case ShortGrass, Fern, DoubleTallGrass, Flower, DoubleFlower, NetherSprouts, PinkPetals, SugarCane, DeadBush, BambooSapling, Bamboo:
		return true
	}
	return false
}


func (g Grass) RandomTick(pos cube.Pos, tx *world.Tx, r *rand.Rand) {
	aboveLight := tx.Light(pos.Side(cube.FaceUp))
	if aboveLight < 4 {
		
		tx.SetBlock(pos, Dirt{}, nil)
		return
	}
	if aboveLight < 9 {
		
		return
	}

	
	n := r.Uint32()

	
	for i := 0; i < 4; i++ {
		x, y, z := int(n)%3, int(n>>2)%5, int(n>>5)%3
		n >>= 7

		spreadPos := pos.Add(cube.Pos{x - 1, y - 3, z - 1})
		
		if tx.Light(spreadPos.Side(cube.FaceUp)) < 4 {
			continue
		}
		b := tx.Block(spreadPos)
		if dirt, ok := b.(Dirt); !ok || dirt.Coarse {
			continue
		}
		tx.SetBlock(spreadPos, g, nil)
	}
}


func (g Grass) BoneMeal(pos cube.Pos, tx *world.Tx) (result item.BoneMealResult) {
	result = item.BoneMealResultNone
	for range 14 {
		c := pos.Add(cube.Pos{rand.IntN(6) - 3, 0, rand.IntN(6) - 3})
		above := c.Side(cube.FaceUp)
		_, air := tx.Block(above).(Air)
		_, grass := tx.Block(c).(Grass)
		if air && grass {
			tx.SetBlock(above, plantSelection[rand.IntN(len(plantSelection))], nil)
			result = item.BoneMealResultArea
		}
	}
	return
}


func (g Grass) BreakInfo() BreakInfo {
	return newBreakInfo(0.6, alwaysHarvestable, shovelEffective, silkTouchOneOf(Dirt{}, g))
}


func (Grass) CompostChance() float64 {
	return 0.3
}


func (Grass) EncodeItem() (name string, meta int16) {
	return "minecraft:grass_block", 0
}


func (Grass) EncodeBlock() (string, map[string]any) {
	return "minecraft:grass_block", nil
}


func (g Grass) Till() (world.Block, bool) {
	return Farmland{}, true
}


func (g Grass) Shovel() (world.Block, bool) {
	return DirtPath{}, true
}
