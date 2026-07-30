package framework

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)


type PlayerStats struct {
	Wins    int `json:"wins"`
	Kills   int `json:"kills"`
	Deaths  int `json:"deaths"`
	Games   int `json:"games"`
}


type StatsManager struct {
	mu       sync.RWMutex
	dir      string
	cache    map[string]*PlayerStats
	logger   *LogWrapper
}


func NewStatsManager(dataDir string, logger *LogWrapper) (*StatsManager, error) {
	dir := filepath.Join(dataDir, "stats")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create stats dir: %w", err)
	}
	return &StatsManager{
		dir:    dir,
		cache:  make(map[string]*PlayerStats),
		logger: logger,
	}, nil
}


func (sm *StatsManager) LoadPlayerStats(uuid uuid.UUID) *PlayerStats {
	sm.mu.RLock()
	if stats, ok := sm.cache[uuid.String()]; ok {
		sm.mu.RUnlock()
		return stats
	}
	sm.mu.RUnlock()

	stats, err := sm.readFromDisk(uuid)
	if err != nil {
		stats = &PlayerStats{}
	}

	sm.mu.Lock()
	sm.cache[uuid.String()] = stats
	sm.mu.Unlock()

	return stats
}


func (sm *StatsManager) SavePlayerStats(uuid uuid.UUID, stats *PlayerStats) error {
	sm.mu.Lock()
	sm.cache[uuid.String()] = stats
	sm.mu.Unlock()

	return sm.writeToDisk(uuid, stats)
}


func (sm *StatsManager) IncrementWins(uuid uuid.UUID) {
	stats := sm.LoadPlayerStats(uuid)
	stats.Wins++
	stats.Games++
	if err := sm.SavePlayerStats(uuid, stats); err != nil {
		sm.logger.Error("Failed to save stats", "uuid", uuid, "error", err)
	}
}


func (sm *StatsManager) IncrementKills(uuid uuid.UUID) {
	stats := sm.LoadPlayerStats(uuid)
	stats.Kills++
	if err := sm.SavePlayerStats(uuid, stats); err != nil {
		sm.logger.Error("Failed to save stats", "uuid", uuid, "error", err)
	}
}


func (sm *StatsManager) IncrementDeaths(uuid uuid.UUID) {
	stats := sm.LoadPlayerStats(uuid)
	stats.Deaths++
	if err := sm.SavePlayerStats(uuid, stats); err != nil {
		sm.logger.Error("Failed to save stats", "uuid", uuid, "error", err)
	}
}


func (sm *StatsManager) IncrementGames(uuid uuid.UUID) {
	stats := sm.LoadPlayerStats(uuid)
	stats.Games++
	if err := sm.SavePlayerStats(uuid, stats); err != nil {
		sm.logger.Error("Failed to save stats", "uuid", uuid, "error", err)
	}
}

func (sm *StatsManager) readFromDisk(uuid uuid.UUID) (*PlayerStats, error) {
	path := filepath.Join(sm.dir, uuid.String()+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var stats PlayerStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

func (sm *StatsManager) writeToDisk(uuid uuid.UUID, stats *PlayerStats) error {
	path := filepath.Join(sm.dir, uuid.String()+".json")
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
