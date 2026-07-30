# Changelog

## 2026-07-29 - Project Initialisation
- Created project structure: plugin-src/Minigames/
- Created .ai memory system
- Explored Dragonfly plugin API (plugin.go, handler.go, world/tx.go, etc.)
- Designed framework architecture (interfaces, data flow, lifecycle)
- Created framework/ package with core types (types.go)
- Created Game interface, GameRegistry, GameManager
- Created Match interface, MatchManager
- Created WorldManager with WorldPool (LevelDB copy + recycle)
- Created PlayerManager with tracking
- Created Countdown timer
- Documented Dragonfly API integration points
