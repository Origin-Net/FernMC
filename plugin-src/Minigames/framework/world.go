package framework

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Origin-Net/FernMC/server/block"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/Origin-Net/FernMC/server/world/biome"
	"github.com/Origin-Net/FernMC/server/world/generator"
	"github.com/Origin-Net/FernMC/server/world/mcdb"
	"github.com/google/uuid"
)


type WorldPool struct {
	mu      sync.Mutex
	idle    []*world.World
	inUse   map[string]*world.World
	maxIdle int
	dir     string
	reg     world.BlockRegistry
	logger  *LogWrapper
	closed  bool
}


func NewWorldPool(dir string, maxIdle int, reg world.BlockRegistry, logger *LogWrapper) *WorldPool {
	return &WorldPool{
		idle:    make([]*world.World, 0, maxIdle),
		inUse:   make(map[string]*world.World),
		maxIdle: maxIdle,
		dir:     dir,
		reg:     reg,
		logger:  logger,
	}
}


func (p *WorldPool) Acquire(templateDir string, arena Arena) (*world.World, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil, fmt.Errorf("world pool is closed")
	}

	
	if len(p.idle) > 0 {
		w := p.idle[len(p.idle)-1]
		p.idle = p.idle[:len(p.idle)-1]
		p.inUse[arena.WorldName] = w
		return w, nil
	}

	
	matchID := uuid.New().String()[:8]
	matchDir := filepath.Join(p.dir, "matches", matchID)
	if err := copyDir(templateDir, matchDir); err != nil {
		return nil, fmt.Errorf("copy world: %w", err)
	}

	w, err := p.openWorld(matchDir, arena.Name)
	if err != nil {
		return nil, err
	}
	p.inUse[arena.WorldName] = w
	return w, nil
}


func (p *WorldPool) Preload(templateDir string, arena Arena, count int) error {
	for range count {
		matchID := uuid.New().String()[:8]
		matchDir := filepath.Join(p.dir, "matches", matchID)
		if err := copyDir(templateDir, matchDir); err != nil {
			return fmt.Errorf("preload copy: %w", err)
		}
		w, err := p.openWorld(matchDir, arena.Name)
		if err != nil {
			return err
		}
		p.idle = append(p.idle, w)
	}
	return nil
}

func (p *WorldPool) openWorld(matchDir string, name string) (*world.World, error) {
	provider, err := mcdb.Config{
		Blocks: p.reg,
	}.Open(matchDir)
	if err != nil {
		return nil, fmt.Errorf("open world %s: %w", name, err)
	}

	w := world.Config{
		Dim:       world.Overworld,
		Provider:  provider,
		Generator: world.NopGenerator{},
		Blocks:    p.reg,
	}.New()

	return w, nil
}


func (p *WorldPool) Release(w *world.World, reset bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	name := findWorldName(p.inUse, w)
	if name != "" {
		delete(p.inUse, name)
	}

	if p.closed || !reset {
		w.Close()
		return
	}

	if len(p.idle) < p.maxIdle {
		p.idle = append(p.idle, w)
	} else {
		w.Close()
	}
}

func findWorldName(m map[string]*world.World, w *world.World) string {
	for name, w2 := range m {
		if w2 == w {
			return name
		}
	}
	return ""
}


func (p *WorldPool) Shutdown() {
	p.mu.Lock()
	p.closed = true
	all := append(p.idle, values(p.inUse)...)
	p.idle = nil
	p.inUse = make(map[string]*world.World)
	p.mu.Unlock()

	for _, w := range all {
		w.Close()
	}
}

func values(m map[string]*world.World) []*world.World {
	out := make([]*world.World, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}


type WorldManager struct {
	pool      *WorldPool
	dataDir   string
	logger    *LogWrapper
	matchDirs map[string]string
	mu        sync.Mutex
}


func NewWorldManager(dataDir string, logger *LogWrapper, reg world.BlockRegistry) *WorldManager {
	return &WorldManager{
		pool:      NewWorldPool(dataDir, 3, reg, logger),
		dataDir:   dataDir,
		logger:    logger,
		matchDirs: make(map[string]string),
	}
}


func (wm *WorldManager) OpenWorld(dir string) (*world.World, error) {
	provider, err := mcdb.Config{
		Blocks: wm.pool.reg,
	}.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("open world %s: %w", dir, err)
	}
	w := world.Config{
		Dim:       world.Overworld,
		Provider:  provider,
		Generator: world.NopGenerator{},
		Blocks:    wm.pool.reg,
	}.New()
	return w, nil
}


func (wm *WorldManager) CreateWorld(dir string) (*world.World, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create world dir: %w", err)
	}
	provider, err := mcdb.Config{
		Blocks: wm.pool.reg,
	}.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("open new world %s: %w", dir, err)
	}
	w := world.Config{
		Dim:       world.Overworld,
		Provider:  provider,
		Generator: generator.NewFlat(biome.Plains{}, []world.Block{block.Grass{}, block.Dirt{}, block.Dirt{}, block.Bedrock{}}),
		Blocks:    wm.pool.reg,
	}.New()
	return w, nil
}


func (wm *WorldManager) AcquireWorld(arena Arena) (*world.World, error) {
	templateDir := filepath.Join(wm.dataDir, "maps", arena.WorldName)
	w, err := wm.pool.Acquire(templateDir, arena)
	if err != nil {
		return nil, fmt.Errorf("acquire world: %w", err)
	}
	return w, nil
}


func (wm *WorldManager) ReturnWorld(worldName string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	
	wm.logger.Info("World returned to pool", "name", worldName)
}


func (wm *WorldManager) Shutdown() {
	wm.pool.Shutdown()
}


func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if entry.Name() == "db" {
				if err := copyLevelDBDir(srcPath, dstPath); err != nil {
					return err
				}
				continue
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
			continue
		}

		data, err := os.ReadFile(srcPath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dstPath, data, 0644); err != nil {
			return err
		}
	}
	return nil
}

func copyLevelDBDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(src, entry.Name()))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dst, entry.Name()), data, 0644); err != nil {
			return err
		}
	}
	return nil
}
