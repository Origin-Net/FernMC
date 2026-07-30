# Architecture

## Package Structure
```
plugin-src/Minigames/
├── plugin.go              # Plugin entry point (main package, exports Plugin)
├── commands.go            # /join, /leave, /games commands
├── framework/
│   ├── types.go           # MatchState, common types
│   ├── game.go            # Game interface, GameManager, GameRegistry
│   ├── match.go           # Match interface, MatchManager
│   ├── world.go           # WorldPool, world creation/copy/reset
│   ├── player.go          # PlayerTracker, player state
│   ├── countdown.go       # Reusable countdown timer
│   ├── spectator.go       # Spectator mode management
│   ├── loot.go            # Loot configuration and generation
│   ├── stats.go           # Player statistics (JSON)
│   ├── map.go             # Map scanning, rotation
│   └── config.go          # Configuration loading
├── games/
│   └── skywars/
│       ├── game.go        # SkywarsGame - implements framework.Game
│       ├── match.go       # Skywars match logic
│       ├── handler.go     # Player event handler (death, quit, etc.)
│       └── config.go      # Skywars-specific config
```

## Core Interfaces

### Game (framework/game.go)
```go
type Game interface {
    Name() string
    ID() string
    MinPlayers() int
    MaxPlayers() int
    CreateMatch(world *world.World, arena Arena, cfg GameConfig) (Match, error)
}
```

### Match (framework/match.go)
```go
type Match interface {
    ID() string
    Game() Game
    State() MatchState
    World() *world.World
    Players() []uuid.UUID
    AddPlayer(p *player.Player) error
    RemovePlayer(p *player.Player, reason RemoveReason)
    Start() error
    End(winner uuid.UUID, reason EndReason)
    HandlePlayerDeath(p *player.Player, killer *player.Player)
    HandlePlayerQuit(p *player.Player)
    Tick(currentTick int64)
}
```

## Data Flow
1. Player joins server → auto-added to player tracking
2. Player runs `/join` → GameManager finds best match
3. If no suitable match → new match created with idle world from pool
4. Match WAITING → countdown starts when MinPlayers reached
5. Countdown ends → STARTING state (teleport to spawns, give kits)
6. Match PLAYING → PvP enabled, chests lootable
7. Player death → SPECTATING state (fly, track remaining)
8. Last player standing → match ends (FINISHED)
9. Statistics saved → players transferred to lobby
10. World returned to pool → state resets

## Match State Machine
```
WAITING → COUNTDOWN → STARTING → PLAYING → FINISHED → CLOSED
                ↑                      ↓
          (players leave)     (last player dies)
```

## World Lifecycle
```
IDLE → IN_USE (match running) → RESETTING → IDLE
  ↑                                    │
  └────────────────────────────────────┘
```

## Player Lifecycle
```
JOINED_MATCH → PLAYING → SPECTATING → TRANSFERRED
                           ↓
                        QUIT_MATCH → REMOVED
```

## Dependencies
- framework/ has no dependencies on games/
- games/ depends on framework/
- plugin.go depends on both
