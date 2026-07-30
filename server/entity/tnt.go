package entity

import (
	"github.com/Origin-Net/FernMC/server/block"
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/internal/nbtconv"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math"
	"math/rand/v2"
	"time"
)


func NewTNT(opts world.EntitySpawnOpts, fuse time.Duration) *world.EntityHandle {
	conf := tntConf
	conf.ExistenceDuration = fuse
	if opts.Velocity.Len() == 0 {
		angle := rand.Float64() * math.Pi * 2
		opts.Velocity = mgl64.Vec3{-math.Sin(angle) * 0.02, 0.1, -math.Cos(angle) * 0.02}
	}
	return opts.New(TNTType, conf)
}

var tntConf = PassiveBehaviourConfig{
	Gravity: 0.04,
	Drag:    0.02,
	Expire:  explodeTNT,
}


func explodeTNT(e *Ent, tx *world.Tx) {
	block.ExplosionConfig{ItemDropChance: 1}.Explode(tx, e.Position())
}


var TNTType tntType

type tntType struct{}

func (t tntType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Ent{tx: tx, handle: handle, data: data}
}

func (tntType) EncodeEntity() string   { return "minecraft:tnt" }
func (tntType) NetworkOffset() float64 { return 0.49 }
func (tntType) BBox(world.Entity) cube.BBox {
	return cube.Box(-0.49, 0, -0.49, 0.49, 0.98, 0.49)
}

func (t tntType) DecodeNBT(m map[string]any, data *world.EntityData) {
	conf := tntConf
	conf.ExistenceDuration = nbtconv.TickDuration[uint8](m, "Fuse")
	data.Data = conf.New()
}

func (tntType) EncodeNBT(data *world.EntityData) map[string]any {
	return map[string]any{"Fuse": uint8(data.Data.(*PassiveBehaviour).Fuse().Milliseconds() / 50)}
}
