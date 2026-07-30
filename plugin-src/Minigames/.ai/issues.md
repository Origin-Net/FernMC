# API Gaps & Design Decisions

## Missing API: No Command Block

**Requirement**: Map builders need to define spawn points in-world (e.g., with a command block).

**Reality**: Dragonfly has no `block.CommandBlock` type. `grep -r CommandBlock server/block/` returns nothing.

**Possible solutions**:
1. **Use `block.Barrier` as spawn marker** (current implementation) — Standard community convention for minigame maps. Map builders place barrier blocks at spawn positions; the plugin scans for them. Barriers do not exist as survival-obtainable items so players can't interfere.
2. **Use `block.StructureBlock`** — Not confirmed to exist in Dragonfly either.
3. **Arena config file only** — No in-world marking; all spawns defined in `arena.yml` manually.
4. **Mixed approach** — Support both barrier markers AND config file overrides.

**Status**: Pending decision.

## Design Decision: Player Spawn Detection

**Requirement**: Plugin must detect which spawn a player was assigned to, to fill nearby chests.

**Approach**: Each spawn gets an index. Island chests are assigned to the nearest spawn by distance. Centre chests (ender chests) are shared.

**Alternatives**: Tag chests with named spawners/armor stands? Not yet investigated.

**Status**: Current implementation uses distance-based assignment. Acceptable?

## Design Decision: Storage Backend

**Requirement**: Persist player stats (wins, kills, etc.) across restarts.

**Options**:
1. **JSON files** (current) — Simple, human-readable, no dependencies. OK for small servers. Concurrent writes may corrupt under heavy load.
2. **SQLite** — Safe concurrent access, fast, requires `modernc.org/sqlite` dependency. Better for production.
3. **YAML** — Similar to JSON but more human-editable.

**Status**: Pending decision.

## Design Decision: Spawn Protection / Grace Period

**Requirement**: Players should not take damage immediately after spawning.

**Approach**: Use `player.Hurt` handler to cancel damage during a configurable grace period (default 5s).

**Status**: Implementation-ready once confirmed.

## Design Decision: Map Distribution Format

**Requirement**: Users should be able to add minigame maps easily.

**Approach**: `.mcworld` files (ZIP archives) placed in `plugins/Minigames/maps/` are automatically extracted on startup. Alternatively, pre-extracted world directories can be placed there directly.

**Status**: Implementation-ready once confirmed.

## Confirmed: Plugin Loading Mechanism

**Reality**: Dragonfly uses Go's `plugin` package (`-buildmode=plugin`). The existing `build_plugin.sh` confirms this. The server auto-discovers `plugin-src/Minigames/`, builds it, loads `plugins/Minigames.pl`, and calls `OnLoad`/`OnEnable`/etc.

**Requirements for the plugin**:
- `package main`
- Exported `var Plugin server.Plugin`
- `go.mod` pointing to the Dragonfly module (with `replace` directive for local development)

**Status**: Correct.

## Confirmed: Event Interception

**Reality**: There is no global event bus. Per-player event handling is done via `player.Handler` interface (implemented by each player's handler).

**Approach**: Replace each player's handler with a custom handler that delegates to the game's match handler.

**Status**: Implementation-ready.

## Confirmed: World API

- `world.Tx.BlocksWithin(pos, radius, ...blockTypes)` — Can search for barrier blocks and chests.
- `mcdb.Config{Blocks: reg}.Open(path)` — Can open an existing LevelDB world.
- `world.Config{Dim, Provider, Generator, Blocks}.New()` — Creates a runtime world instance.

**Status**: Implementation is correct.
