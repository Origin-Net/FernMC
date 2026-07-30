package world

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"time"
)

var (
	
	
	
	Overworld overworld
	
	
	Nether nether
	
	
	End end
)

var dimensionReg = newDimensionRegistry(map[int]Dimension{
	0: Overworld,
	1: Nether,
	2: End,
})




func DimensionByID(id int) (Dimension, bool) {
	return dimensionReg.Lookup(id)
}



func DimensionID(dim Dimension) (int, bool) {
	return dimensionReg.LookupID(dim)
}

type dimensionRegistry struct {
	dimensions map[int]Dimension
	ids        map[Dimension]int
}


func newDimensionRegistry(dim map[int]Dimension) *dimensionRegistry {
	ids := make(map[Dimension]int, len(dim))
	for k, v := range dim {
		ids[v] = k
	}
	return &dimensionRegistry{dimensions: dim, ids: ids}
}




func (reg *dimensionRegistry) Lookup(id int) (Dimension, bool) {
	dim, ok := reg.dimensions[id]
	if !ok {
		dim = Overworld
	}
	return dim, ok
}



func (reg *dimensionRegistry) LookupID(dim Dimension) (int, bool) {
	id, ok := reg.ids[dim]
	return id, ok
}

type (
	
	
	
	Dimension interface {
		
		
		Range() cube.Range
		WaterEvaporates() bool
		LavaSpreadDuration() time.Duration
		WeatherCycle() bool
		TimeCycle() bool
	}
	overworld struct{}
	nether    struct{}
	end       struct{}
)

func (overworld) Range() cube.Range                 { return cube.Range{-64, 319} }
func (overworld) WaterEvaporates() bool             { return false }
func (overworld) LavaSpreadDuration() time.Duration { return time.Second * 3 / 2 }
func (overworld) WeatherCycle() bool                { return true }
func (overworld) TimeCycle() bool                   { return true }
func (overworld) String() string                    { return "Overworld" }

func (nether) Range() cube.Range                 { return cube.Range{0, 127} }
func (nether) WaterEvaporates() bool             { return true }
func (nether) LavaSpreadDuration() time.Duration { return time.Second / 4 }
func (nether) WeatherCycle() bool                { return false }
func (nether) TimeCycle() bool                   { return false }
func (nether) String() string                    { return "Nether" }

func (end) Range() cube.Range                 { return cube.Range{0, 255} }
func (end) WaterEvaporates() bool             { return false }
func (end) LavaSpreadDuration() time.Duration { return time.Second * 3 / 2 }
func (end) WeatherCycle() bool                { return false }
func (end) TimeCycle() bool                   { return false }
func (end) String() string                    { return "End" }
