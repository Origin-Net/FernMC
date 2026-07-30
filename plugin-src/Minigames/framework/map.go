package framework

import (
	"archive/zip"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Origin-Net/FernMC/server/block"
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/Origin-Net/FernMC/server/world/mcdb"
	"github.com/df-mc/goleveldb/leveldb"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"gopkg.in/yaml.v3"
)


type blockFound struct {
	pos  cube.Pos
	name string
}


type MapManager struct {
	dataDir   string
	mapsDir   string
	arenas    []Arena
	recent    []string
	maxRecent int
	mu        sync.RWMutex
	logger    *LogWrapper
	blockReg  world.BlockRegistry
}


func NewMapManager(dataDir string, logger *LogWrapper, blockReg world.BlockRegistry) *MapManager {
	return &MapManager{
		dataDir:   dataDir,
		mapsDir:   filepath.Join(dataDir, "maps"),
		maxRecent: 3,
		logger:    logger,
		blockReg:  blockReg,
	}
}



func (mm *MapManager) DiscoverMaps(gameID string) error {
	if err := os.MkdirAll(mm.mapsDir, 0755); err != nil {
		return fmt.Errorf("create maps dir: %w", err)
	}

	entries, err := os.ReadDir(mm.mapsDir)
	if err != nil {
		return fmt.Errorf("read maps dir: %w", err)
	}

	for _, entry := range entries {
		name := entry.Name()

		if strings.HasSuffix(name, ".mcworld") && !entry.IsDir() {
			if err := mm.extractMcworld(name); err != nil {
				mm.logger.Error("Failed to extract map", "file", name, "error", err)
			}
			continue
		}

		if entry.IsDir() {
			if err := mm.scanMapDirectory(name); err != nil {
				mm.logger.Error("Failed to scan map", "name", name, "error", err)
			}
		}
	}

	mm.logger.Info("Map discovery complete", "maps_found", len(mm.arenas))
	return nil
}

func (mm *MapManager) extractMcworld(filename string) error {
	zipPath := filepath.Join(mm.mapsDir, filename)
	mapName := strings.TrimSuffix(filename, ".mcworld")
	extractDir := filepath.Join(mm.mapsDir, mapName)

	if _, err := os.Stat(extractDir); err == nil {
		mm.logger.Info("Map already extracted", "name", mapName)
		return nil
	}

	mm.logger.Info("Extracting map", "file", filename)

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", filename, err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(extractDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}
		out, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			out.Close()
			return err
		}
		_, err = io.Copy(out, rc)
		out.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (mm *MapManager) scanMapDirectory(name string) error {
	mapDir := filepath.Join(mm.mapsDir, name)

	
	dbDir := filepath.Join(mapDir, "db")
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		
		entries, err := os.ReadDir(mapDir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				subDbDir := filepath.Join(mapDir, entry.Name(), "db")
				if _, err := os.Stat(subDbDir); err == nil {
					return mm.scanMapDirectory(filepath.Join(name, entry.Name()))
				}
			}
		}
		return fmt.Errorf("no LevelDB db directory found in %s", mapDir)
	}

	
	arenaPath := filepath.Join(mapDir, "arena.yml")
	arena, err := mm.loadArenaConfig(arenaPath)
	if err == nil {
		mm.arenas = append(mm.arenas, arena)
		mm.logger.Info("Loaded arena from config", "name", arena.Name, "spawns", len(arena.Spawns))
		return nil
	}

	var spawns []mgl64.Vec3
	var chests []ChestPosition

	
	foundBlocks, err := scanLevelDBForBlocks(dbDir, "minecraft:command_block", "minecraft:chest", "minecraft:ender_chest")
	if err != nil {
		mm.logger.Warn("Direct LevelDB scan failed, falling back to world scan", "error", err)
		
		provider, pErr := mcdb.Config{Blocks: mm.blockReg}.Open(mapDir)
		if pErr != nil {
			return fmt.Errorf("open world: %w", pErr)
		}
		defer provider.Close()

		w := world.Config{
			Dim:       world.Overworld,
			Provider:  provider,
			Generator: world.NopGenerator{},
			Blocks:    mm.blockReg,
		}.New()
		defer w.Close()

		w.Do(func(tx *world.Tx) {
			center := cube.Pos{0, 64, 0}
			var cmdBlocks []world.Block
			for _, face := range cube.Faces() {
				cmdBlocks = append(cmdBlocks, block.CommandBlock{Facing: face, Conditional: false})
				cmdBlocks = append(cmdBlocks, block.CommandBlock{Facing: face, Conditional: true})
			}
			for pos := range tx.BlocksWithin(center, 200, cmdBlocks...) {
				spawnPos := mgl64.Vec3{float64(pos.X()) + 0.5, float64(pos.Y()) + 1, float64(pos.Z()) + 0.5}
				spawns = append(spawns, spawnPos)
			}
			chestB := block.Chest{}
			enderChest := block.EnderChest{}
			for pos := range tx.BlocksWithin(center, 200, chestB) {
				chests = append(chests, ChestPosition{
					Pos:  mgl64.Vec3{float64(pos.X()) + 0.5, float64(pos.Y()), float64(pos.Z()) + 0.5},
					Type: ChestTypeNormal,
				})
			}
			for pos := range tx.BlocksWithin(center, 200, enderChest) {
				chests = append(chests, ChestPosition{
					Pos:  mgl64.Vec3{float64(pos.X()) + 0.5, float64(pos.Y()), float64(pos.Z()) + 0.5},
					Type: ChestTypeCenter,
				})
			}
		})
	} else {
		for _, b := range foundBlocks {
			switch b.name {
			case "minecraft:command_block":
				spawnPos := mgl64.Vec3{float64(b.pos.X()) + 0.5, float64(b.pos.Y()) + 1, float64(b.pos.Z()) + 0.5}
				spawns = append(spawns, spawnPos)
			case "minecraft:chest":
				chests = append(chests, ChestPosition{
					Pos:  mgl64.Vec3{float64(b.pos.X()) + 0.5, float64(b.pos.Y()), float64(b.pos.Z()) + 0.5},
					Type: ChestTypeNormal,
				})
			case "minecraft:ender_chest":
				chests = append(chests, ChestPosition{
					Pos:  mgl64.Vec3{float64(b.pos.X()) + 0.5, float64(b.pos.Y()), float64(b.pos.Z()) + 0.5},
					Type: ChestTypeCenter,
				})
			}
		}
	}

	if len(spawns) == 0 {
		return fmt.Errorf("no spawn points (command blocks) found in %s", name)
	}

	arena = Arena{
		Name:       name,
		WorldName:  name,
		Spawns:     spawns,
		Chests:     chests,
		MaxPlayers: len(spawns),
	}

	
	mm.saveArenaConfig(mapDir, arena)

	mm.mu.Lock()
	mm.arenas = append(mm.arenas, arena)
	mm.mu.Unlock()

	mm.logger.Info("Scanned map",
		"name", name,
		"spawns", len(spawns),
		"chests", len(chests),
	)

	return nil
}

func (mm *MapManager) loadArenaConfig(path string) (Arena, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Arena{}, err
	}
	var arena Arena
	if err := yaml.Unmarshal(data, &arena); err != nil {
		return Arena{}, err
	}
	return arena, nil
}

func (mm *MapManager) saveArenaConfig(dir string, arena Arena) {
	data, err := yaml.Marshal(arena)
	if err != nil {
		mm.logger.Error("Failed to marshal arena config", "name", arena.Name, "error", err)
		return
	}
	if err := os.WriteFile(filepath.Join(dir, "arena.yml"), data, 0644); err != nil {
		mm.logger.Error("Failed to save arena config", "name", arena.Name, "error", err)
	}
}


func (mm *MapManager) SelectArena(gameID string) (Arena, error) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if len(mm.arenas) == 0 {
		return Arena{}, fmt.Errorf("no arenas available")
	}

	if len(mm.arenas) == 1 {
		return mm.arenas[0], nil
	}

	
	type candidate struct {
		arena  Arena
		weight float64
	}
	var candidates []candidate

	for _, a := range mm.arenas {
		w := 1.0
		for i, recent := range mm.recent {
			if recent == a.Name {
				w -= 0.3 * float64(mm.maxRecent-i) / float64(mm.maxRecent)
			}
		}
		if w < 0.1 {
			w = 0.1
		}
		candidates = append(candidates, candidate{arena: a, weight: w})
	}

	
	var totalWeight float64
	for _, c := range candidates {
		totalWeight += c.weight
	}

	r := rand.Float64() * totalWeight
	for _, c := range candidates {
		r -= c.weight
		if r <= 0 {
			mm.recordSelection(c.arena.Name)
			return c.arena, nil
		}
	}

	fallback := candidates[len(candidates)-1].arena
	mm.recordSelection(fallback.Name)
	return fallback, nil
}

func (mm *MapManager) recordSelection(name string) {
	mm.recent = append(mm.recent, name)
	if len(mm.recent) > mm.maxRecent {
		mm.recent = mm.recent[1:]
	}
}




func scanLevelDBForBlocks(dbDir string, targetNames ...string) ([]blockFound, error) {
	ldb, err := leveldb.OpenFile(dbDir, nil)
	if err != nil {
		return nil, fmt.Errorf("open leveldb: %w", err)
	}
	defer ldb.Close()

	targetSet := make(map[string]bool, len(targetNames))
	for _, n := range targetNames {
		targetSet[n] = true
	}

	idToBlock := map[string]string{
		"CommandBlock": "minecraft:command_block",
		"Chest":        "minecraft:chest",
		"EnderChest":   "minecraft:ender_chest",
	}

	var results []blockFound
	iter := ldb.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		key := iter.Key()
		val := iter.Value()

		
		
		if len(val) == 0 {
			continue
		}
		
		if val[0] != 0x0A {
			continue
		}

		
		var cx, cz int32
		switch len(key) {
		case 9: 
			cx = int32(binary.LittleEndian.Uint32(key[0:4]))
			cz = int32(binary.LittleEndian.Uint32(key[4:8]))
		case 10: 
			cx = int32(binary.LittleEndian.Uint32(key[0:4]))
			cz = int32(binary.LittleEndian.Uint32(key[4:8]))
		default:
			continue
		}
		
		if (int(cx)*16+16 < -200) || (int(cx)*16 > 200) ||
			(int(cz)*16+16 < -200) || (int(cz)*16 > 200) {
			continue
		}

		var nbtData map[string]any
		if err := nbt.UnmarshalEncoding(val, &nbtData, nbt.LittleEndian); err != nil {
			continue
		}
		id, _ := nbtData["id"].(string)
		if id == "" {
			continue
		}
		blockName, ok := idToBlock[id]
		if !ok || !targetSet[blockName] {
			continue
		}

		x, _ := nbtData["x"].(int32)
		y, _ := nbtData["y"].(int32)
		z, _ := nbtData["z"].(int32)

		if x < -200 || x > 200 || z < -200 || z > 200 {
			continue
		}

		results = append(results, blockFound{
			pos:  cube.Pos{int(x), int(y), int(z)},
			name: blockName,
		})
	}

	return results, nil
}


func (mm *MapManager) ArenaByName(name string) (Arena, bool) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	for _, a := range mm.arenas {
		if a.Name == name {
			return a, true
		}
	}
	return Arena{}, false
}


func (mm *MapManager) Arenas() []Arena {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	out := make([]Arena, len(mm.arenas))
	copy(out, mm.arenas)
	return out
}
