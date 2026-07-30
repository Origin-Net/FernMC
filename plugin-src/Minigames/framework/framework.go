package framework

import (
	"fmt"
	"log/slog"

	"github.com/Origin-Net/FernMC/server"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/world"
)


type Framework struct {
	dataDir     string
	logger      *LogWrapper
	config      *ConfigManager
	games       *GameManager
	matches     *MatchManager
	worlds      *WorldManager
	players     *PlayerManager
	stats       *StatsManager
	loot        *LootManager
	maps        *MapManager
	holograms   *HologramManager
	defaultWorld *world.World
	srv         *server.Server
}



func (f *Framework) SetupLobby() error {
	if f.defaultWorld == nil {
		return fmt.Errorf("default world not set")
	}
	f.players.SetLobbyWorld(f.defaultWorld)
	f.holograms.SetLobbyWorld(f.defaultWorld)
	f.logger.Info("Lobby world set to default world")
	return nil
}


func (f *Framework) SetDefaultWorld(w *world.World) {
	f.defaultWorld = w
}


func (f *Framework) SetServer(srv *server.Server) {
	f.srv = srv
	f.players.SetServer(srv)
}


type lobbyHandler struct {
	world.Handler
	pm *PlayerManager
}

func (h lobbyHandler) HandleEntitySpawn(tx *world.Tx, e world.Entity) {
	if p, ok := e.(*player.Player); ok {
		h.pm.SendLobbyScoreboard(p)
	}
	h.Handler.HandleEntitySpawn(tx, e)
}


func (f *Framework) AttachLobbyHandler() {
	if f.defaultWorld == nil {
		f.logger.Warn("Cannot attach lobby handler: default world not set")
		return
	}
	existing := f.defaultWorld.Handler()
	slog.Info("Attaching lobby handler", "existing_type", fmt.Sprintf("%T", existing))
	h := lobbyHandler{
		Handler: existing,
		pm:      f.players,
	}
	f.defaultWorld.Handle(h)
	f.logger.Info("Lobby handler attached for scoreboard on player spawn")
}


func New(dataDir string, logger *LogWrapper, blockReg world.BlockRegistry) (*Framework, error) {
	config, err := NewConfigManager(dataDir, logger)
	if err != nil {
		return nil, err
	}

	worlds := NewWorldManager(dataDir, logger, blockReg)
	players := NewPlayerManager(dataDir, logger)
	stats, err := NewStatsManager(dataDir, logger)
	if err != nil {
		return nil, err
	}
	loot, err := NewLootManager(dataDir, logger)
	if err != nil {
		return nil, err
	}
	maps := NewMapManager(dataDir, logger, blockReg)
	holograms := NewHologramManager(dataDir, logger)
	matchMgr := NewMatchManager(worlds, logger)

	registry := NewGameRegistry()
	gameMgr := NewGameManager(registry, matchMgr, worlds, players, maps, logger)

	players.SetStatsManager(stats)

	return &Framework{
		dataDir: dataDir,
		logger:  logger,
		config:  config,
		games:   gameMgr,
		matches: matchMgr,
		worlds:  worlds,
		players: players,
		stats:   stats,
		loot:      loot,
		maps:      maps,
		holograms: holograms,
	}, nil
}


func (f *Framework) DiscoverMaps() error {
	return f.maps.DiscoverMaps("")
}


func (f *Framework) Shutdown() {
	f.matches.Shutdown()
	f.worlds.Shutdown()
}


func (f *Framework) DataFolder() string { return f.dataDir }


func (f *Framework) Logger() *LogWrapper { return f.logger }


func (f *Framework) Config() *ConfigManager { return f.config }


func (f *Framework) Games() *GameManager { return f.games }


func (f *Framework) Matches() *MatchManager { return f.matches }


func (f *Framework) Worlds() *WorldManager { return f.worlds }


func (f *Framework) Players() *PlayerManager { return f.players }


func (f *Framework) Stats() *StatsManager { return f.stats }


func (f *Framework) Loot() *LootManager { return f.loot }


func (f *Framework) Holograms() *HologramManager { return f.holograms }


func (f *Framework) Maps() *MapManager { return f.maps }
