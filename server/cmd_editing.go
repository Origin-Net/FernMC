package server

import (
	"strings"

	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type setblockCmd struct {
	Position         mgl64.Vec3
	Block            BlockName
	BlockData        cmd.Optional[string]
	OldBlockHandling cmd.Optional[OldBlockHandling]
}

func (s setblockCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	pos := cube.PosFromVec3(s.Position)
	name := strings.TrimPrefix(string(s.Block), "minecraft:")
	block, ok := world.BlockByName(name, nil)
	if !ok {
		o.Error("Unknown block: " + string(s.Block))
		return
	}
	tx.SetBlock(pos, block, nil)
	o.Printf("Set block at %v to %s", s.Position, string(s.Block))
}

type fillCmd struct {
	From             mgl64.Vec3
	To               mgl64.Vec3
	Block            BlockName
	BlockData        cmd.Optional[string]
	OldBlockHandling cmd.Optional[OldBlockHandling]
}

func (f fillCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	name := strings.TrimPrefix(string(f.Block), "minecraft:")
	block, ok := world.BlockByName(name, nil)
	if !ok {
		o.Error("Unknown block: " + string(f.Block))
		return
	}
	min := cube.PosFromVec3(mgl64.Vec3{
		min(f.From[0], f.To[0]),
		min(f.From[1], f.To[1]),
		min(f.From[2], f.To[2]),
	})
	max := cube.PosFromVec3(mgl64.Vec3{
		max(f.From[0], f.To[0]),
		max(f.From[1], f.To[1]),
		max(f.From[2], f.To[2]),
	})
	count := 0
	for x := min[0]; x <= max[0]; x++ {
		for y := min[1]; y <= max[1]; y++ {
			for z := min[2]; z <= max[2]; z++ {
				tx.SetBlock(cube.Pos{x, y, z}, block, nil)
				count++
			}
		}
	}
	o.Printf("Filled %d blocks", count)
}

func init() {
	cmd.Register(cmd.New("setblock", "Changes a block to another block", nil, setblockCmd{}))
	cmd.Register(cmd.New("fill", "Fills all or parts of a region with a specific block", nil, fillCmd{}))
}
