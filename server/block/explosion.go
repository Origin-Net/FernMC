package block

import (
	"math"
	"math/rand/v2"
	"time"

	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/block/cube/trace"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/Origin-Net/FernMC/server/world/particle"
	"github.com/Origin-Net/FernMC/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)



type ExplosionConfig struct {
	
	Size float64
	
	
	
	RandSource rand.Source
	
	
	SpawnFire bool
	
	
	
	ItemDropChance float64

	
	
	Sound world.Sound
	
	
	Particle world.Particle
}


type ExplodableEntity interface {
	
	
	Explode(explosionPos mgl64.Vec3, impact float64, c ExplosionConfig)
}


type Explodable interface {
	
	Explode(explosionPos mgl64.Vec3, pos cube.Pos, tx *world.Tx, c ExplosionConfig)
}


var rays = make([]mgl64.Vec3, 0, 1352)


func init() {
	for x := 0.0; x < 16; x++ {
		for y := 0.0; y < 16; y++ {
			for z := 0.0; z < 16; z++ {
				if x != 0 && x != 15 && y != 0 && y != 15 && z != 0 && z != 15 {
					continue
				}
				rays = append(rays, mgl64.Vec3{x/15*2 - 1, y/15*2 - 1, z/15*2 - 1}.Normalize().Mul(0.3))
			}
		}
	}
}


func (c ExplosionConfig) Explode(tx *world.Tx, explosionPos mgl64.Vec3) {
	if c.Sound == nil {
		c.Sound = sound.Explosion{}
	}
	if c.Particle == nil {
		c.Particle = particle.HugeExplosion{}
	}
	if c.RandSource == nil {
		t := uint64(time.Now().UnixNano())
		c.RandSource = rand.NewPCG(t, t)
	}
	if c.Size == 0 {
		c.Size = 4
	}
	if c.ItemDropChance == 0 {
		c.ItemDropChance = 1.0 / c.Size
	}

	r, d := rand.New(c.RandSource), c.Size*2
	box := cube.Box(
		math.Floor(explosionPos[0]-d-1),
		math.Floor(explosionPos[1]-d-1),
		math.Floor(explosionPos[2]-d-1),
		math.Ceil(explosionPos[0]+d+1),
		math.Ceil(explosionPos[1]+d+1),
		math.Ceil(explosionPos[2]+d+1),
	)

	affectedEntities := make([]world.Entity, 0, 32)
	for e := range tx.EntitiesWithin(box.Grow(2)) {
		pos := e.Position()
		dist := pos.Sub(explosionPos).Len()
		if dist > d || dist == 0 {
			continue
		}

		affectedEntities = append(affectedEntities, e)
	}

	affectedBlocks := make([]cube.Pos, 0, 32)
	for _, ray := range rays {
		pos := explosionPos
		for blastForce := c.Size * (0.7 + r.Float64()*0.6); blastForce > 0.0; blastForce -= 0.225 {
			current := cube.PosFromVec3(pos)
			currentBlock := tx.Block(current)

			resistance := 0.0
			if l, ok := tx.Liquid(current); ok {
				resistance = l.BlastResistance()
			} else if i, ok := currentBlock.(Breakable); ok {
				resistance = i.BreakInfo().BlastResistance
			} else if _, ok = currentBlock.(Air); !ok {
				
				break
			}

			pos = pos.Add(ray)
			if blastForce -= (resistance/5 + 0.3) * 0.3; blastForce > 0 {
				affectedBlocks = append(affectedBlocks, current)
			}
		}
	}

	ctx := tx.Event()
	spawnFire := c.SpawnFire
	itemDropChance := c.ItemDropChance
	if tx.World().Handler().HandleExplosion(ctx, explosionPos, &affectedEntities, &affectedBlocks, &itemDropChance, &spawnFire); ctx.Cancelled() {
		return
	}

	for _, e := range affectedEntities {
		if explodable, ok := e.(ExplodableEntity); ok {
			impact := (1 - e.Position().Sub(explosionPos).Len()/d) * exposure(tx, explosionPos, e)
			explodable.Explode(explosionPos, impact, c)
		}
	}

	for _, pos := range affectedBlocks {
		bl := tx.Block(pos)
		if explodable, ok := bl.(Explodable); ok {
			explodable.Explode(explosionPos, pos, tx, c)
		} else if breakable, ok := bl.(Breakable); ok {
			
			tx.SetBlock(pos, nil, nil)
			breakHandler := breakable.BreakInfo().BreakHandler
			if breakHandler != nil {
				breakHandler(pos, tx, nil)
			}
			if itemDropChance > r.Float64() {
				for _, drop := range breakable.BreakInfo().Drops(item.ToolNone{}, nil) {
					dropItem(tx, drop, pos.Vec3Centre())
				}
			}
		}
	}

	if spawnFire {
		for _, pos := range affectedBlocks {
			if r.IntN(3) == 0 {
				if _, ok := tx.Block(pos).(Air); ok && tx.Block(pos.Side(cube.FaceDown)).Model().FaceSolid(pos, cube.FaceUp, tx) {
					Fire{}.Start(tx, pos)
				}
			}
		}
	}

	tx.AddParticle(explosionPos, c.Particle)
	tx.PlaySound(explosionPos, c.Sound)
}


func exposure(tx *world.Tx, origin mgl64.Vec3, e world.Entity) float64 {
	pos := e.Position()
	box := e.H().Type().BBox(e).Translate(pos)

	boxMin, boxMax := box.Min(), box.Max()
	diff := boxMax.Sub(boxMin).Mul(2.0).Add(mgl64.Vec3{1, 1, 1})

	step := mgl64.Vec3{1.0 / diff[0], 1.0 / diff[1], 1.0 / diff[2]}
	if step[0] < 0.0 || step[1] < 0.0 || step[2] < 0.0 {
		return 0.0
	}

	xOffset := (1.0 - math.Floor(diff[0])/diff[0]) / 2.0
	zOffset := (1.0 - math.Floor(diff[2])/diff[2]) / 2.0

	var checks, misses float64
	for x := 0.0; x <= 1.0; x += step[0] {
		for y := 0.0; y <= 1.0; y += step[1] {
			for z := 0.0; z <= 1.0; z += step[2] {
				point := mgl64.Vec3{
					lerp(x, boxMin[0], boxMax[0]) + xOffset,
					lerp(y, boxMin[1], boxMax[1]),
					lerp(z, boxMin[2], boxMax[2]) + zOffset,
				}
				var collided bool
				trace.TraverseBlocks(origin, point, func(pos cube.Pos) (cont bool) {
					_, collided = trace.BlockIntercept(pos, tx, tx.Block(pos), origin, point)
					return !collided
				})

				if !collided {
					misses++
				}
				checks++
			}
		}
	}
	return misses / checks
}


func lerp(a, b, t float64) float64 {
	return b + a*(t-b)
}
