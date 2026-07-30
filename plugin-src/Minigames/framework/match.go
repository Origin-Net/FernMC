package framework

import (
	"fmt"
	"sync"
	"time"

	"github.com/Origin-Net/FernMC/server/player"
	"github.com/google/uuid"
)


type Match interface {
	
	ID() string
	
	GameID() string
	
	State() MatchState
	
	Start()
	
	Done() <-chan struct{}
	
	Players() []*player.Player
	
	Alive() []*player.Player
	
	AddPlayer(p *player.Player) error
}


type MatchManager struct {
	mu       sync.RWMutex
	matches  map[string]Match
	worldMgr *WorldManager
	logger   *LogWrapper
}


func NewMatchManager(worldMgr *WorldManager, logger *LogWrapper) *MatchManager {
	return &MatchManager{
		matches:  make(map[string]Match),
		worldMgr: worldMgr,
		logger:   logger,
	}
}


func (m *MatchManager) Register(match Match) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.matches[match.ID()] = match
	m.logger.Info("Match registered", "id", match.ID(), "game", match.GameID())
}


func (m *MatchManager) Unregister(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.matches, id)
}


func (m *MatchManager) Get(id string) (Match, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	match, ok := m.matches[id]
	return match, ok
}


func (m *MatchManager) FindSuitableMatch(gameID string, maxPlayers int) Match {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var best Match
	var bestCount int

	for _, match := range m.matches {
		if match.GameID() != gameID {
			continue
		}
		if match.State() != MatchStateWaiting {
			continue
		}
		count := len(match.Players())
		if count >= maxPlayers {
			continue
		}
		if count > bestCount {
			best = match
			bestCount = count
		}
	}
	return best
}


func (m *MatchManager) MatchByPlayer(p *player.Player) (Match, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, match := range m.matches {
		for _, mp := range match.Players() {
			if mp == p {
				return match, true
			}
		}
	}
	return nil, false
}


func (m *MatchManager) Shutdown() {
	m.mu.RLock()
	matches := make([]Match, 0, len(m.matches))
	for _, match := range m.matches {
		matches = append(matches, match)
	}
	m.mu.RUnlock()

	for _, match := range matches {
		if s := match.State(); s == MatchStatePlaying || s == MatchStateWaiting || s == MatchStateCountdown {
			m.logger.Info("Shutting down match", "id", match.ID())
		}
	}
}


func NewMatchID(gameID string) string {
	return fmt.Sprintf("%s_%s_%d", gameID, uuid.New().String()[:8], time.Now().Unix())
}
