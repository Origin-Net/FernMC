package world

import (
	"sync"
	"sync/atomic"

	"github.com/Origin-Net/FernMC/server/block/cube"
)



type Settings struct {
	sync.Mutex
	ref atomic.Int32

	
	Name string
	
	Spawn cube.Pos
	
	Time int64
	
	TimeCycle bool
	
	RainTime int64
	
	Raining bool
	
	ThunderTime int64
	
	Thundering bool
	
	WeatherCycle bool
	
	RequiredSleepTicks int64
	
	
	CurrentTick int64
	
	DefaultGameMode GameMode
	
	
	Difficulty Difficulty
	
	
	TickRange int32
}


func defaultSettings() *Settings {
	return &Settings{
		Name:            "World",
		DefaultGameMode: GameModeSurvival,
		Difficulty:      DifficultyNormal,
		TimeCycle:       true,
		WeatherCycle:    true,
		TickRange:       6,
	}
}
