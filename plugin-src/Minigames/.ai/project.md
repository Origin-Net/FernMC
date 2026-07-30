# Minigames Framework for Dragonfly

## Project Goal
Production-ready reusable minigames framework for Dragonfly Minecraft Bedrock server.

## Technology Stack
- Go 1.26
- Dragonfly (github.com/df-mc/dragonfly)
- Go plugin system (`-buildmode=plugin`)
- LevelDB (world storage)
- JSON file storage (player statistics)
- YAML configuration (loot tables, game config)
- WaterdogPE proxy integration

## Network Architecture
- Bedrock Player → WaterdogPE Proxy → Lobby Server → Skywars Server
- One Skywars server manages multiple match worlds
- Matches are worlds, not separate servers

## Permanent Rules
- Use existing Dragonfly plugin API - do NOT modify Dragonfly core
- Matches are independent worlds within one server process
- All game-specific code lives in games/ subdirectory
- Framework package must be reusable across game types
- No hardcoded values in framework code
