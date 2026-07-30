package block

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/world"
)





type CommandBlock struct {
	solid
	sourceWaterDisplacer

	
	Facing cube.Face

	
	
	Conditional bool
}


func (CommandBlock) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return true
}


func (CommandBlock) EncodeItem() (name string, meta int16) {
	return "minecraft:command_block", 0
}


func (c CommandBlock) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:command_block", map[string]any{
		"conditional_bit":  boolByte(c.Conditional),
		"facing_direction": int32(c.Facing),
	}
}


func allCommandBlocks() (b []world.Block) {
	for _, face := range cube.Faces() {
		b = append(b, CommandBlock{Facing: face, Conditional: false})
		b = append(b, CommandBlock{Facing: face, Conditional: true})
	}
	return
}
