package skywars

import (
	"fmt"
	"sync"

	"github.com/Origin-Net/FernMC/server/player"
	"gopkg.in/yaml.v3"
	"github.com/Origin-Net/FernMC/plugin-src/Minigames/framework"
)


type SkyWars struct {
	mu            sync.RWMutex
	config        *Config
	framework     *framework.Framework
}


func New(fw *framework.Framework) *SkyWars {
	return &SkyWars{
		config:    DefaultConfig(),
		framework: fw,
	}
}


func (s *SkyWars) ID() string {
	return "skywars"
}


func (s *SkyWars) Name() string {
	return "SkyWars"
}


func (s *SkyWars) Description() string {
	return "Battle on floating islands – the last one standing wins!"
}


func (s *SkyWars) MinPlayers() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config.MinPlayers
}


func (s *SkyWars) MaxPlayers() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config.MaxPlayers
}


func (s *SkyWars) CreateMatch(id string, arena framework.Arena, players []*player.Player) framework.Match {
	return NewMatch(id, arena, players, s.config, s.framework)
}


func (s *SkyWars) LoadConfig(data []byte) error {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parse skywars config: %w", err)
	}
	s.mu.Lock()
	s.config = &cfg
	s.mu.Unlock()
	return nil
}


func (s *SkyWars) Config() *Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}


func (s *SkyWars) ChestFiller(chestType framework.ChestType) []map[string]int {
	switch chestType {
	case framework.ChestTypeCenter:
		return s.framework.Loot().GenerateLoot("center")
	default:
		return s.framework.Loot().GenerateLoot("normal")
	}
}


func getRandomSpawnChests(arena framework.Arena, spawnIdx int, count int) []framework.ChestPosition {
	
	type distChest struct {
		chest framework.ChestPosition
		dist  float64
	}

	spawn := arena.Spawns[spawnIdx]
	var nearby []distChest

	for _, c := range arena.Chests {
		if c.Type == framework.ChestTypeCenter {
			
			continue
		}
		dx := c.Pos.X() - spawn.X()
		dz := c.Pos.Z() - spawn.Z()
		dist := dx*dx + dz*dz
		nearby = append(nearby, distChest{chest: c, dist: dist})
	}

	
	for i := 1; i < len(nearby); i++ {
		for j := i; j > 0 && nearby[j].dist < nearby[j-1].dist; j-- {
			nearby[j], nearby[j-1] = nearby[j-1], nearby[j]
		}
	}

	var result []framework.ChestPosition
	for _, nc := range nearby {
		if len(result) >= count {
			break
		}
		result = append(result, nc.chest)
	}

	return result
}


func generateIslandLoot(chests []framework.ChestPosition, generator func(framework.ChestType) []map[string]int) map[framework.ChestPosition][]map[string]int {
	chestLoot := make(map[framework.ChestPosition][]map[string]int)
	for _, chest := range chests {
		chestLoot[chest] = generator(chest.Type)
	}
	return chestLoot
}
