package particle

import (
	"image/color"

	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/Origin-Net/FernMC/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)



type Flame struct {
	particle
	
	Colour color.RGBA
}



type Dust struct {
	particle

	
	Colour color.RGBA
}



type BlockBreak struct {
	particle
	
	
	Block world.Block
}



type PunchBlock struct {
	particle
	
	
	Block world.Block
	
	Face cube.Face
}


type BlockForceField struct{ particle }


type BoneMeal struct {
	particle

	
	
	
	Area bool
}


type Note struct {
	particle

	
	Instrument sound.Instrument
	
	Pitch int
}


type DragonEggTeleport struct {
	particle

	
	Diff cube.Pos
}


type Evaporate struct{ particle }


type WaterDrip struct{ particle }


type LavaDrip struct{ particle }


type Lava struct{ particle }


type DustPlume struct{ particle }


type particle struct{}


func (particle) Spawn(*world.World, mgl64.Vec3) {}
