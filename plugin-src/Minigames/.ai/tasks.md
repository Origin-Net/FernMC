# Implementation Plan (Pending Approval)

## Files to Create

| # | File | Purpose | Dependencies |
|---|------|---------|-------------|
| 1 | `plugin-src/Minigames/plugin.go` | Plugin entry point: `var Plugin server.Plugin`, implements `Meta()`, `OnLoad`, `OnEnable`, `OnDisable`, `OnUnload`. Registers commands, inits framework. | All framework packages, `server.Plugin` |
| 2 | `plugin-src/Minigames/go.mod` | Module file with `replace` directive to local Dragonfly | — |
| 3 | `plugin-src/Minigames/commands.go` | Commands: `/join <game>`, `/leave`, `/games` (list games). Uses `cmd.Command` / `cmd.New`. | framework, `server/cmd` |
| 4 | `plugin-src/Minigames/framework/types.go` | ✅ Done — Core types |
| 5 | `plugin-src/Minigames/framework/game.go` | ✅ Done — Game interface, GameRegistry, GameManager |
| 6 | `plugin-src/Minigames/framework/match.go` | ✅ Done — Match interface, MatchManager, lifecycle states |
| 7 | `plugin-src/Minigames/framework/world.go` | ✅ Done — WorldManager, WorldPool, world lifecycle |
| 8 | `plugin-src/Minigames/framework/player.go` | ✅ Done — PlayerManager, player tracking |
| 9 | `plugin-src/Minigames/framework/countdown.go` | ✅ Done — CountdownTimer |
| 10 | `plugin-src/Minigames/framework/spectator.go` | ✅ Done — SpectatorManager |
| 11 | `plugin-src/Minigames/framework/loot.go` | ✅ Done — LootManager, loot table loading |
| 12 | `plugin-src/Minigames/framework/stats.go` | ✅ Done — StatsManager (JSON) |
| 13 | `plugin-src/Minigames/framework/map.go` | ✅ Done — MapManager, map discovery, arena scanning |
| 14 | `plugin-src/Minigames/framework/config.go` | ✅ Done — ConfigManager, default config |
| 15 | `plugin-src/Minigames/framework/log.go` | ✅ Done — LogWrapper |
| 16 | `plugin-src/Minigames/games/skywars/game.go` | ✅ Done — SkyWars game struct + chest filler |
| 17 | `plugin-src/Minigames/games/skywars/match.go` | ⏳ Not written — Match lifecycle (wait → in-game → end), ref logic, kill feed |
| 18 | `plugin-src/Minigames/games/skywars/handler.go` | ⏳ Not written — Player handler for death, chat, block break |
| 19 | `plugin-src/Minigames/games/skywars/config.go` | ⏳ Not written — SkyWars-specific config |
| 20 | `plugin-src/Minigames/config/config.yml` | ⏳ Not written — Default config file |
| 21 | `plugin-src/Minigames/config/loot/normal.yml` | ⏳ Not written — Island loot table |
| 22 | `plugin-src/Minigames/config/loot/center.yml` | ⏳ Not written — Centre loot table |

## Build & Test

| Step | Command | Expected |
|------|---------|----------|
| Build | `go build -buildmode=plugin -o plugins/Minigames.pl ./plugin-src/Minigames/` | Produces `plugins/Minigames.pl` |
| Auto-load | Start Dragonfly server | Server builds from `plugin-src/Minigames/` automatically, loads `plugins/Minigames.pl` |
| Test join | Player runs `/join skywars` | Player teleported to waiting lobby |
| Test start | Enough players join | Countdown begins, game starts |
| Test death | Player dies | Player enters spectator mode |
| Test end | One player remains | Winner announced, players returned to lobby/stats saved |
| Test reload | `/plugin reload Minigames` | Plugin reloaded, commands re-registered |

## Open Questions for Approval

1. Spawn marking: Use barrier blocks in-world, arena config file, or both?
2. Storage: JSON files or SQLite?
3. Island chest assignment: Nearest-chest distance or explicit chest-spawn linkage?
4. Map format: `.mcworld` auto-extraction or manual LevelDB directories?
5. Plugin name: "Minigames" or "SkyWars" (if SkyWars-only initially)?
