package world

import "github.com/Origin-Net/FernMC/server/block/cube"

type SpawnEggHandler func(tx *Tx, pos cube.Pos) bool

var spawnEggHandlers = map[string]SpawnEggHandler{}

func RegisterSpawnEggHandler(name string, handler SpawnEggHandler) {
	spawnEggHandlers[name] = handler
}

func SpawnEggHandlerByName(name string) (SpawnEggHandler, bool) {
	h, ok := spawnEggHandlers[name]
	return h, ok
}
