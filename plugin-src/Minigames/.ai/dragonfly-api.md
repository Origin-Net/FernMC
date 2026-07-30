# Dragonfly Plugin API Reference

## Plugin System
- File: `server/plugin.go`
- Interface: `server.Plugin` (Meta, OnLoad, OnEnable, OnDisable, OnUnload)
- Entry: export `var Plugin server.Plugin` in main package
- Build: `go build -buildmode=plugin -o plugins/<name>.pl plugin-src/<name>/`
- Context: `PluginContext{DataFolder, Logger, Server}`
- Server access: `GetServer()` returns `*Server`

## Server API
- `Server.World()` → overworld `*world.World`
- `Server.Players(tx)` → `iter.Seq[*player.Player]`
- `Server.Player(uuid)` → `(*world.EntityHandle, bool)`
- `Server.PlayerByName(name)` → `(*world.EntityHandle, bool)`
- `Server.Logger()` → `*slog.Logger`
- `Server.World().BlockRegistry()` → `world.BlockRegistry`
- `Server.World().Config()` → `world.Config` (has .Entities, .Blocks)

## World API
- Creation: `world.Config{...}.New()` 
- Config fields: Log, Dim, Provider, Generator, Entities, Blocks, SaveInterval, RandomTickSpeed, ReadOnly
- `world.New()` → in-memory world with NopProvider
- `World.Save()` → saves chunks to provider
- `World.Close()` → closes world, stops ticking
- `World.Name()` → string
- `World.Dimension()` → Dimension
- `World.Handle(h Handler)` → sets world event handler
- `World.Handler()` → returns current handler
- Synchronous mode: `Config{Synchronous: true}.New()` for testing

## Transaction System
- `World.Do(func(tx *Tx))` → queues transaction
- `World.Call(ctx, func(tx *Tx) (T, error))` → sync call with result
- `Tx.SetBlock(pos, block, opts)` → set block
- `Tx.Block(pos)` → get block (loads chunk)
- `Tx.BlockLoaded(pos)` → get block (only if loaded)
- `Tx.BlocksWithin(pos, radius, blocks...)` → iter.Seq[cube.Pos] - scans for matching blocks
- `Tx.Players()` → iter.Seq[Entity]
- `Tx.Entities()` → iter.Seq[Entity]
- `Tx.AddEntity(handle)` → Entity
- `Tx.RemoveEntity(e)` → *EntityHandle
- `Tx.AddParticle(pos, particle)` → spawn particle
- `Tx.PlaySound(pos, sound)` → play sound

## Provider API
- `mcdb.Config{}.Open(dir)` → opens LevelDB world, returns `*DB` (implements `world.Provider`)
- `NopProvider{}` → in-memory storage (no persistence)
- Provider from server config: `srv.World().Config().Provider`

## Player API
- `Player.Handle(h Handler)` → set player event handler
- `Player.Name()`, `Player.UUID()`, `Player.XUID()`
- `Player.Position()`, `Player.World()`
- `Player.GameMode()`, `Player.SetGameMode(mode)`
- `Player.Health()`, `Player.MaxHealth()`, `Player.SetMaxHealth(h)`
- `Player.Hurt(dmg, src)`, `Player.Heal(health, src)`
- `Player.Kill()`, `Player.Dead()`, `Player.Respawn()`
- `Player.SetFood(level)`, `Player.Food()`
- `Player.AddEffect(e)`, `Player.RemoveEffect(t)`, `Player.Effects()`
- `Player.Message(msg)`, `Player.SendTip(msg)`, `Player.SendPopup(msg)`
- `Player.SendTitle(title)`, `Player.RemoveTitle()`
- `Player.SendScoreboard(sb)`, `Player.RemoveScoreboard()`
- `Player.SendBossBar(bar)`, `Player.RemoveBossBar()`
- `Player.Teleport(pos)`
- `Player.Transfer(address)` → WaterdogPE transfer
- `Player.Disconnect(msg)`
- `Player.Inventory()`, `Player.EnderChestInventory()`
- `Player.HeldItems()`, `Player.SetHeldItems(main, off)`
- `Player.SetAllowFlight(bool)`, `Player.SetFlying(bool)`
- `Player.SetInvisible()`, `Player.SetImmobile()`
- `Player.SetSpeed(s)`, `Player.SetFlightSpeed(s)`
- `Player.PlaySound(sound)`, `Player.PlaySoundByName(name, pos)`
- `Player.SetNameTag(name)`, `Player.SetScoreTag(a...)`
- `Player.H()` → `*world.EntityHandle` (for cross-world operations)
- `Player.Tx()` → `*world.Tx` (only valid during transaction)

## Player Handler
- Interface: `player.Handler`
- Key methods: HandleHurt, HandleDeath, HandleRespawn, HandleQuit, HandleFoodLoss, HandleAttackEntity, HandleItemDrop, HandleChat, HandleChangeWorld, HandleTransfer
- Embed `player.NopHandler` for default no-op behaviour
- Set via `p.Handle(myHandler)` - replaces the current handler
- Death control: `HandleDeath` receives `keepInv *bool`
- Respawn control: `HandleRespawn` receives `pos *mgl64.Vec3`, `w **world.World`
- Context: `*player.Context` - `Cancel()` to cancel event, `.Cancelled()` to check

## Command System
- Registration: `cmd.New(name, desc, aliases, runnables...)` + `cmd.Register(cmd)`
- Struct-based with reflection: exported fields = command args
- Types: `[]cmd.Target` (player selectors), `string`, `int`, `float64`, `bool`, `mgl64.Vec3`, `cmd.Varargs` (rest), `cmd.SubCommand`, `cmd.Optional[T]`
- Enums: any `string` type with `Type()` and `Options(Source)` methods
- `cmd.Source` interface - `Position()`, `SendCommandOutput(output)`
- `cmd.Output` - `Print(msg)`, `Error(msg)`, `Printf(fmt, args...)`, `SetResult(v)`

## World Handler
- Interface: `world.Handler`
- Methods: HandleClose(*Tx), HandleLoadChunk(*Tx, pos), HandleUnloadChunk(*Tx, pos)
- Set via `world.Handle(h)`

## Entity System
- `entity.DefaultRegistry` - all default entity types
- `entity.SpawnEntityType` - various entity types
- `EntityType.CreateEntity(initInfo)` → `*EntityHandle`
- `entity.Effect` - lasting and instant effects
- `world.Entity` interface: Position(), Close(), etc.

## Block System
- Block identification: `EncodeBlock() (string, map[string]any)`
- Block lookup: `world.BlockByName(name, props)` → `(Block, bool)`
- `world.Blocks()` → all registered blocks
- Block state: immutable, `blockState.WithProperty(prop, val)` returns new state
- Block entities: `SetBlockEntity(pos, block)` - chests, furnaces, etc.
- Key block: `barrier.Barrier{}` - used for spawn detection

## Missing/Unavailable APIs
- No command block type (`minecraft:command_block`) - use BARRIER blocks for spawns
- No chunk iteration on World (unexported `chunks` map) - use `BlocksWithin` with large radius
- No event bus for plugins - use player handlers for event interception
- No scheduler API - use `World.DoAfter()` or spawn goroutines
- No world rename method - world name is in Settings/level.dat
