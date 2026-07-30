package framework

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/entity"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
)


type Hologram struct {
	UUID     string          `json:"uuid"`
	World    string          `json:"world"`
	X, Y, Z  float64         `json:"position"`
	Lines    []string        `json:"lines"`
	handles []*world.EntityHandle
	mu       sync.Mutex
}


type HologramManager struct {
	mu        sync.RWMutex
	holograms  map[string]*Hologram
	dir       string
	logger    *LogWrapper
	lobbyWorld *world.World
}


func NewHologramManager(dataDir string, logger *LogWrapper) *HologramManager {
	hm := &HologramManager{
		holograms: make(map[string]*Hologram),
		dir:       filepath.Join(dataDir, "holograms"),
		logger:    logger,
	}
	_ = os.MkdirAll(hm.dir, 0755)
	hm.loadAll()
	return hm
}


func (hm *HologramManager) SetLobbyWorld(w *world.World) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.lobbyWorld = w
	
	for name, h := range hm.holograms {
		hm.removeEntities(h)
		if w != nil {
			hm.addEntities(name, h, w)
		}
	}
}


func (hm *HologramManager) Create(name string, pos mgl64.Vec3, lines []string) (*Hologram, error) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	if _, ok := hm.holograms[name]; ok {
		return nil, fmt.Errorf("hologram %q already exists", name)
	}

	worldName := "lobby"
	if hm.lobbyWorld != nil {
		worldName = hm.lobbyWorld.Name()
	}

	h := &Hologram{
		UUID:  uuid.New().String(),
		World: worldName,
		X:     pos.X(), Y: pos.Y(), Z: pos.Z(),
		Lines: lines,
	}
	if len(h.Lines) == 0 {
		h.Lines = []string{"§eEmpty Hologram"}
	}

	hm.holograms[name] = h
	if hm.lobbyWorld != nil {
		hm.addEntities(name, h, hm.lobbyWorld)
	}
	hm.save(name, h)
	return h, nil
}


func (hm *HologramManager) Delete(name string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	h, ok := hm.holograms[name]
	if !ok {
		return
	}
	hm.removeEntities(h)
	delete(hm.holograms, name)
	path := filepath.Join(hm.dir, name+".json")
	_ = os.Remove(path)
}


func (hm *HologramManager) Get(name string) (*Hologram, bool) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	h, ok := hm.holograms[name]
	return h, ok
}


func (hm *HologramManager) All() map[string]*Hologram {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	out := make(map[string]*Hologram, len(hm.holograms))
	for k, v := range hm.holograms {
		out[k] = v
	}
	return out
}


func (hm *HologramManager) SetLines(name string, lines []string) error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	h, ok := hm.holograms[name]
	if !ok {
		return fmt.Errorf("hologram %q not found", name)
	}

	hm.removeEntities(h)
	h.mu.Lock()
	h.Lines = lines
	h.mu.Unlock()

	if hm.lobbyWorld != nil {
		hm.addEntities(name, h, hm.lobbyWorld)
	}
	hm.save(name, h)
	return nil
}


func (hm *HologramManager) MoveTo(name string, pos mgl64.Vec3) error {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	h, ok := hm.holograms[name]
	if !ok {
		return fmt.Errorf("hologram %q not found", name)
	}

	hm.removeEntities(h)
	h.mu.Lock()
	h.X, h.Y, h.Z = pos.Elem()
	h.mu.Unlock()

	if hm.lobbyWorld != nil {
		hm.addEntities(name, h, hm.lobbyWorld)
	}
	hm.save(name, h)
	return nil
}

func (hm *HologramManager) addEntities(name string, h *Hologram, w *world.World) {
	h.mu.Lock()
	defer h.mu.Unlock()

	lineH := 0.3
	startY := h.Y + float64(len(h.Lines)-1)*lineH

	for i, line := range h.Lines {
		pos := mgl64.Vec3{h.X, startY - float64(i)*lineH, h.Z}
		handle := entity.NewText(line, pos)
		w.Do(func(tx *world.Tx) {
			tx.AddEntity(handle)
		})
		h.handles = append(h.handles, handle)
	}
}

func (hm *HologramManager) removeEntities(h *Hologram) {
	h.mu.Lock()
	handles := h.handles
	h.handles = nil
	h.mu.Unlock()

	for _, handle := range handles {
		handle.Do(func(tx *world.Tx, e world.Entity) {
			tx.RemoveEntity(e)
		})
	}
}

func (hm *HologramManager) save(name string, h *Hologram) {
	path := filepath.Join(hm.dir, name+".json")
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		hm.logger.Error("Failed to marshal hologram", "name", name, "error", err)
		return
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		hm.logger.Error("Failed to save hologram", "name", name, "error", err)
	}
}

func (hm *HologramManager) loadAll() {
	entries, err := os.ReadDir(hm.dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		name := entry.Name()[:len(entry.Name())-len(".json")]
		path := filepath.Join(hm.dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			hm.logger.Error("Failed to read hologram file", "path", path, "error", err)
			continue
		}
		var h Hologram
		if err := json.Unmarshal(data, &h); err != nil {
			hm.logger.Error("Failed to parse hologram", "name", name, "error", err)
			continue
		}
		hm.holograms[name] = &h
		hm.logger.Info("Loaded hologram", "name", name, "lines", len(h.Lines))
	}
}



func HologramCommands(hm *HologramManager) []cmd.Command {
	return []cmd.Command{
		cmd.New("hd", "Hologram management", nil,
			hdCmd{hm: hm},
		),
	}
}

type hdCmd struct {
	hm    *HologramManager
	Verb  cmd.Optional[string] `cmd:"verb"`   
	Name  cmd.Optional[string] `cmd:"name"`
	Text  cmd.Optional[string] `cmd:"text"`
	Index cmd.Optional[int]    `cmd:"index"`
}

func (c hdCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, isPlayer := src.(*player.Player)
	verb, _ := c.Verb.Load()
	name, _ := c.Name.Load()
	text, _ := c.Text.Load()
	idx, _ := c.Index.Load()

	if verb == "" {
		verb = "list"
	}

	switch verb {
	case "list":
		all := c.hm.All()
		if len(all) == 0 {
			o.Print("§7No holograms.")
			return
		}
		o.Print("§6Holograms:")
		for hn, h := range all {
			o.Print(fmt.Sprintf("  §e%s §7(%s) §8[%d lines]", hn, h.World, len(h.Lines)))
		}

	case "create":
		if !isPlayer {
			o.Error("Only players can use this command.")
			return
		}
		h, err := c.hm.Create(name, p.Position(), []string{text})
		if err != nil {
			o.Error(err.Error())
			return
		}
		_ = h
		o.Print(fmt.Sprintf("§aHologram §e%s §acreated.", name))

	case "delete":
		c.hm.Delete(name)
		o.Print(fmt.Sprintf("§cHologram §e%s §cdeleted.", name))

	case "addline":
		h, ok := c.hm.Get(name)
		if !ok {
			o.Error(fmt.Sprintf("Hologram %q not found.", name))
			return
		}
		lines := append(h.Lines, text)
		if err := c.hm.SetLines(name, lines); err != nil {
			o.Error(err.Error())
			return
		}
		o.Print(fmt.Sprintf("§aAdded line to §e%s§a.", name))

	case "removeline":
		h, ok := c.hm.Get(name)
		if !ok {
			o.Error(fmt.Sprintf("Hologram %q not found.", name))
			return
		}
		if idx < 0 || idx >= len(h.Lines) {
			o.Error(fmt.Sprintf("Index out of range (0-%d).", len(h.Lines)-1))
			return
		}
		lines := make([]string, 0, len(h.Lines)-1)
		lines = append(lines, h.Lines[:idx]...)
		lines = append(lines, h.Lines[idx+1:]...)
		if err := c.hm.SetLines(name, lines); err != nil {
			o.Error(err.Error())
			return
		}
		o.Print(fmt.Sprintf("§aRemoved line %d from §e%s§a.", idx, name))

	case "setline":
		h, ok := c.hm.Get(name)
		if !ok {
			o.Error(fmt.Sprintf("Hologram %q not found.", name))
			return
		}
		if idx < 0 || idx >= len(h.Lines) {
			o.Error(fmt.Sprintf("Index out of range (0-%d).", len(h.Lines)-1))
			return
		}
		lines := make([]string, len(h.Lines))
		copy(lines, h.Lines)
		lines[idx] = text
		if err := c.hm.SetLines(name, lines); err != nil {
			o.Error(err.Error())
			return
		}
		o.Print(fmt.Sprintf("§aUpdated line %d of §e%s§a.", idx, name))

	case "movehere":
		if !isPlayer {
			o.Error("Only players can use this command.")
			return
		}
		if err := c.hm.MoveTo(name, p.Position()); err != nil {
			o.Error(err.Error())
			return
		}
		o.Print(fmt.Sprintf("§aMoved §e%s §ato your position.", name))

	default:
		o.Error(fmt.Sprintf("Unknown verb %q. Use: create, delete, addline, removeline, setline, list, movehere", c.Verb))
	}
}
