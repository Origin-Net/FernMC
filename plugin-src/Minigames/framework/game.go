package framework

import (
	"fmt"
	"sync"

	"github.com/Origin-Net/FernMC/server/player"
)


type Game interface {
	
	ID() string
	
	Name() string
	
	Description() string
	
	MinPlayers() int
	
	MaxPlayers() int
	
	CreateMatch(id string, arena Arena, players []*player.Player) Match
}


type GameRegistry struct {
	mu    sync.RWMutex
	games map[string]Game
}


func NewGameRegistry() *GameRegistry {
	return &GameRegistry{
		games: make(map[string]Game),
	}
}


func (r *GameRegistry) Register(g Game) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := g.ID()
	if _, exists := r.games[id]; exists {
		return fmt.Errorf("game %s is already registered", id)
	}
	r.games[id] = g
	return nil
}


func (r *GameRegistry) Get(id string) (Game, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.games[id]
	return g, ok
}


func (r *GameRegistry) All() []Game {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Game, 0, len(r.games))
	for _, g := range r.games {
		out = append(out, g)
	}
	return out
}


type GameManager struct {
	registry  *GameRegistry
	matchMgr  *MatchManager
	worldMgr  *WorldManager
	playerMgr *PlayerManager
	mapMgr    *MapManager
	logger    *LogWrapper
}


func (m *GameManager) Registry() *GameRegistry { return m.registry }


func NewGameManager(registry *GameRegistry, matchMgr *MatchManager, worldMgr *WorldManager, playerMgr *PlayerManager, mapMgr *MapManager, logger *LogWrapper) *GameManager {
	return &GameManager{
		registry:  registry,
		matchMgr:  matchMgr,
		worldMgr:  worldMgr,
		playerMgr: playerMgr,
		mapMgr:    mapMgr,
		logger:    logger,
	}
}


func (m *GameManager) JoinGame(p *player.Player, gameID string) error {
	game, ok := m.registry.Get(gameID)
	if !ok {
		return fmt.Errorf("unknown game: %s", gameID)
	}

	if m.playerMgr.IsInMatch(p) {
		return fmt.Errorf("you are already in a match")
	}

	match := m.matchMgr.FindSuitableMatch(gameID, game.MaxPlayers())
	if match == nil {
		arena, err := m.mapMgr.SelectArena(gameID)
		if err != nil {
			return fmt.Errorf("no arena available: %w", err)
		}
		id := NewMatchID(gameID)
		match = game.CreateMatch(id, arena, nil)
		m.matchMgr.Register(match)
	}

	if err := match.AddPlayer(p); err != nil {
		return err
	}
	m.playerMgr.AddPlayer(p, match.ID())
	m.logger.Info("Player joined match", "player", p.Name(), "match", match.ID())

	if len(match.Players()) >= game.MinPlayers() {
		
		match.Start()
	}

	return nil
}
