package cmd

import (
	"github.com/Origin-Net/FernMC/server/internal/sliceutil"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"slices"
)



type Target interface {
	
	Position() mgl64.Vec3
}


type NamedTarget interface {
	Target
	
	Name() string
}


func targets(tx *world.Tx) (entities []Target, players []NamedTarget) {
	if tx == nil {
		return nil, nil
	}
	ent := sliceutil.Convert[Target](slices.Collect(tx.Entities()))
	pl := sliceutil.Convert[NamedTarget](slices.Collect(tx.Players()))
	return ent, pl
}
