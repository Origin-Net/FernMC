package cmd

import (
	"strings"

	"github.com/Origin-Net/FernMC/server/entity"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type MobNameEnum string

func (MobNameEnum) Type() string { return "EntityType" }
func (MobNameEnum) Options(Source) []string {
	mobs := entity.MobTypes()
	n := make([]string, len(mobs))
	for i, m := range mobs {
		n[i] = strings.TrimPrefix(m.EncodeEntity(), "minecraft:")
	}
	return n
}

type Summon struct {
	Mob MobNameEnum
}

func (s Summon) Run(src Source, o *Output, tx *world.Tx) {
	if tx == nil {
		o.Error("Cannot summon from console without a world target.")
		return
	}
	name := string(s.Mob)
	if strings.HasPrefix(name, "minecraft:") {
		name = name[10:]
	}
	mobType, ok := entity.MobByName(name)
	if !ok {
		o.Error("Unknown mob: " + string(s.Mob))
		return
	}
	pos := src.Position()
	handle := world.EntitySpawnOpts{Position: pos.Add(mgl64.Vec3{0, 1, 0})}.New(mobType, entity.MobBehaviourConfig)
	tx.AddEntity(handle)
	o.Print("Summoned " + string(s.Mob))
}

func init() {
	Register(New("summon", "Summons a mob", nil, Summon{}))
}
