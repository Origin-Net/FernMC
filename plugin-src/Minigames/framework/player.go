package framework

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/Origin-Net/FernMC/server"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)


type PlayerState struct {
	Player       *player.Player
	MatchID      string
	IsAlive      bool
	IsSpectating bool
	Kills        int
	Deaths       int
}


type PlayerManager struct {
	mu         sync.RWMutex
	players    map[uuid.UUID]*PlayerState
	lobbyWorld *world.World
	lobbySpawn LobbySpawn
	stats      *StatsManager
	srv        *server.Server
	dataDir    string
	logger     *LogWrapper
}


type LobbySpawn struct {
	Position mgl64.Vec3  `yaml:"position"`
	Yaw      float64 `yaml:"yaw"`
	Pitch    float64 `yaml:"pitch"`
}


func NewPlayerManager(dataDir string, logger *LogWrapper) *PlayerManager {
	pm := &PlayerManager{
		players: make(map[uuid.UUID]*PlayerState),
		dataDir: dataDir,
		logger:  logger,
	}
	pm.loadLobbySpawn()
	return pm
}


func (pm *PlayerManager) SetStatsManager(sm *StatsManager) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.stats = sm
}


func (pm *PlayerManager) SetServer(srv *server.Server) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.srv = srv
}


func (pm *PlayerManager) SetLobbyWorld(w *world.World) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.lobbyWorld = w
}


func (pm *PlayerManager) SetSpawn(pos mgl64.Vec3, yaw, pitch float64) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.lobbySpawn = LobbySpawn{Position: pos, Yaw: yaw, Pitch: pitch}
	pm.saveLobbySpawn()
	pm.logger.Info("Lobby spawn set", "pos", pos, "yaw", yaw, "pitch", pitch)
}


func (pm *PlayerManager) LobbySpawn() LobbySpawn {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.lobbySpawn
}

func (pm *PlayerManager) spawnConfigPath() string {
	return filepath.Join(pm.dataDir, "config", "lobby_spawn.yml")
}

func (pm *PlayerManager) saveLobbySpawn() {
	data, err := yaml.Marshal(pm.lobbySpawn)
	if err != nil {
		pm.logger.Error("Failed to marshal lobby spawn", "error", err)
		return
	}
	path := pm.spawnConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		pm.logger.Error("Failed to create config dir", "error", err)
		return
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		pm.logger.Error("Failed to save lobby spawn", "error", err)
	}
}

func (pm *PlayerManager) loadLobbySpawn() {
	data, err := os.ReadFile(pm.spawnConfigPath())
	if err != nil {
		if !os.IsNotExist(err) {
			pm.logger.Error("Failed to read lobby spawn config", "error", err)
		}
		return
	}
	if err := yaml.Unmarshal(data, &pm.lobbySpawn); err != nil {
		pm.logger.Error("Failed to parse lobby spawn config", "error", err)
	}
}


func (pm *PlayerManager) AddPlayer(p *player.Player, matchID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.players[p.UUID()] = &PlayerState{
		Player:  p,
		MatchID: matchID,
		IsAlive: true,
	}
}


func (pm *PlayerManager) RemovePlayer(p *player.Player) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.players, p.UUID())
}


func (pm *PlayerManager) SetSpectating(p *player.Player) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if state, ok := pm.players[p.UUID()]; ok {
		state.IsAlive = false
		state.IsSpectating = true
	}
}


func (pm *PlayerManager) SetDead(p *player.Player) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if state, ok := pm.players[p.UUID()]; ok {
		state.IsAlive = false
		state.Deaths++
	}
}


func (pm *PlayerManager) IncrementKills(p *player.Player) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if state, ok := pm.players[p.UUID()]; ok {
		state.Kills++
	}
}


func (pm *PlayerManager) IsInMatch(p *player.Player) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	_, ok := pm.players[p.UUID()]
	return ok
}


func (pm *PlayerManager) PlayerState(p *player.Player) (*PlayerState, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	state, ok := pm.players[p.UUID()]
	return state, ok
}


func clearInventory(p *player.Player) {
	p.Inventory().Clear()
	p.Armour().Inventory().Clear()
}


func (pm *PlayerManager) ReturnToLobby(p *player.Player) {
	pm.RemovePlayer(p)
	if pm.lobbyWorld == nil {
		return
	}
	lobbyWorld := pm.lobbyWorld
	spawn := pm.LobbySpawn()
	spawnPos := spawn.Position
	if spawnPos.ApproxEqual(mgl64.Vec3{}) {
		spawnPos = lobbyWorld.Spawn().Vec3()
	}
	p.Do(func(tx *world.Tx, pl *player.Player) {
		clearInventory(pl)
		h := tx.RemoveEntity(pl)
		lobbyWorld.Do(func(tx2 *world.Tx) {
			np := tx2.AddEntity(h).(*player.Player)
			np.Handle(nil)
			np.SetGameMode(world.GameModeAdventure)
			np.Teleport(spawnPos)

			
			current := np.Rotation()
			dy := spawn.Yaw - current.Yaw()
			dp := spawn.Pitch - current.Pitch()
			if !mgl64.FloatEqual(dy, 0) || !mgl64.FloatEqual(dp, 0) {
				np.Move(mgl64.Vec3{}, dy, dp)
			}
			pm.SendLobbyScoreboard(np)
		})
	})
	p.Message("§aReturned to lobby.")
}


func (pm *PlayerManager) RemoveAllFromMatch(matchID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	for id, state := range pm.players {
		if state.MatchID == matchID {
			delete(pm.players, id)
		}
	}
}
