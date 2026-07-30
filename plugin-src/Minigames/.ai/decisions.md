# Design Decisions

## ✅ Command Block as Spawn Marker (2026-07-29)

**Decision**: Use `minecraft:command_block` as the SkyWars spawn marker block.

**Rationale**: Command blocks don't exist in Dragonfly's block registry. Rather than substituting a different block (e.g., barrier), we added proper command block support to the Dragonfly core:
- New `block.CommandBlock` type with `Facing` (6 directions) and `Conditional` bool
- Full registration in block hash + accessor system
- Encoded as `minecraft:command_block` with `facing_direction` (int) and `conditional_bit` (bool) block states
- Compatible with vanilla Bedrock world saves containing command blocks

**Impact**:
- `server/block/command_block.go` — new block type
- `server/block/hash.go` - add `hashCommandBlock` + `Hash()` method  
- `server/block/register.go` — register all 12 variants (6 faces × 2 conditional)

## ✅ Plugin Architecture (2026-07-29)

**Decision**: Implement as a Go `plugin` (`-buildmode=plugin`) loaded by Dragonfly's built-in `server/plugin.go` system.

**Structure**:
- `plugin-src/Minigames/` — separate Go module with its own `go.mod`
- Export `var Plugin server.Plugin` symbol
- Server auto-detects and builds from source on startup

## ✅ Spawn Lifecycle

1. **Map discovery** (`MapManager.DiscoverMaps`): Opens the LevelDB world, scans for `block.CommandBlock` via `tx.BlocksWithin`, records positions in `Arena.Spawns`
2. **Match start** (`Match.beginGame`): Teleports each player to `Spawn[i]` position + 1 block above, replaces the command block with `block.Glass{}`

## Summary of Framework Files

| File | Status | Purpose |
|------|--------|---------|
| `plugin.go` | ✅ | Plugin entry point, command registration |
| `commands.go` | ✅ | /join, /leave, /games commands |
| `go.mod` | ✅ | Module definition with replace directive |
| `framework/types.go` | ✅ | Core types (MatchState, Arena, ChestPosition, etc.) |
| `framework/framework.go` | ✅ | Central orchestrator |
| `framework/game.go` | ✅ | Game interface, GameRegistry, GameManager |
| `framework/match.go` | ✅ | Match interface, MatchManager |
| `framework/world.go` | ✅ | WorldPool, WorldManager |
| `framework/player.go` | ✅ | PlayerManager with lobby return |
| `framework/countdown.go` | ✅ | CountdownTimer |
| `framework/spectator.go` | ✅ | SpectatorManager |
| `framework/loot.go` | ✅ | LootManager (YAML tables) |
| `framework/stats.go` | ✅ | StatsManager (JSON) |
| `framework/map.go` | ✅ | MapManager (scans for command blocks + chests) |
| `framework/config.go` | ✅ | ConfigManager (YAML config) |
| `games/skywars/game.go` | ✅ | SkyWars game type |
| `games/skywars/match.go` | ✅ | SkyWars match lifecycle |
| `games/skywars/handler.go` | ✅ | Per-player event handler |
| `games/skywars/config.go` | ✅ | SkyWars-specific config |
