package skywars

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Origin-Net/FernMC/server/block"
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/entity"
	"github.com/Origin-Net/FernMC/server/entity/effect"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/player/form"
	"github.com/Origin-Net/FernMC/server/player/title"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"github.com/Origin-Net/FernMC/plugin-src/Minigames/framework"
)


type Match struct {
	id           string
	arena        framework.Arena
	world        *world.World
	players      []*player.Player
	alive        []*player.Player
	dead         []*player.Player
	spectators   []*player.Player
	spawnIndices map[uuid.UUID]int

	config       *Config
	fw           *framework.Framework
	pvpEnabled   bool
	state        framework.MatchState

	startTime    time.Time
	refillTime   time.Time
	totalKills   int

	mu           sync.RWMutex
	done         chan struct{}
}




type ghostMode struct{}

func (ghostMode) AllowsEditing() bool       { return false }
func (ghostMode) AllowsTakingDamage() bool  { return false }
func (ghostMode) CreativeInventory() bool   { return false }
func (ghostMode) HasCollision() bool        { return true }
func (ghostMode) AllowsFlying() bool        { return true }
func (ghostMode) AllowsInteraction() bool   { return false }
func (ghostMode) Visible() bool             { return false }
func (ghostMode) InstantPortalTravel() bool { return false }


func NewMatch(id string, arena framework.Arena, players []*player.Player, config *Config, fw *framework.Framework) *Match {
	m := &Match{
		id:           id,
		arena:        arena,
		world:        nil,
		players:      players,
		alive:        make([]*player.Player, len(players)),
		spawnIndices: make(map[uuid.UUID]int),
		config:       config,
		fw:           fw,
		state:        framework.MatchStateWaiting,
		done:         make(chan struct{}),
	}
	copy(m.alive, players)

	for i, p := range players {
		m.spawnIndices[p.UUID()] = i
	}

	return m
}


func (m *Match) AddPlayer(p *player.Player) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.state != framework.MatchStateWaiting {
		return fmt.Errorf("match %s is not accepting players", m.id)
	}
	if len(m.players) >= m.config.MaxPlayers {
		return fmt.Errorf("match is full")
	}
	for _, existing := range m.players {
		if existing == p {
			return nil
		}
	}
	m.players = append(m.players, p)
	m.alive = append(m.alive, p)
	m.spawnIndices[p.UUID()] = len(m.players) - 1
	p.Message(fmt.Sprintf("§eWaiting for %d more player(s)...", m.config.MinPlayers-len(m.players)))
	return nil
}


func (m *Match) ID() string { return m.id }


func (m *Match) State() framework.MatchState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}


func (m *Match) GameID() string { return "skywars" }


func (m *Match) Players() []*player.Player {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*player.Player, len(m.players))
	copy(out, m.players)
	return out
}


func (m *Match) Alive() []*player.Player {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*player.Player, len(m.alive))
	copy(out, m.alive)
	return out
}


func (m *Match) Start() {
	m.mu.Lock()
	m.state = framework.MatchStateCountdown
	m.mu.Unlock()

	go m.beginGame()
}

func (m *Match) beginGame() {
	w, err := m.fw.Worlds().AcquireWorld(m.arena)
	if err != nil {
		m.broadcast(fmt.Sprintf("§cFailed to load arena world: %v", err))
		m.endGame(framework.EndReasonForce)
		return
	}

	m.mu.Lock()
	m.state = framework.MatchStateStarting
	m.pvpEnabled = false
	m.world = w
	m.totalKills = 0
	m.startTime = time.Now()
	m.refillTime = m.startTime.Add(time.Duration(m.config.RefillInterval) * time.Second)
	m.mu.Unlock()

	w.Do(func(tx *world.Tx) {
		for i := range m.alive {
			spawn := m.arena.Spawns[i]
			cmdPos := cube.PosFromVec3(spawn.Sub(mgl64.Vec3{0, 1, 0}))
			tx.SetBlock(cmdPos, block.Glass{}, nil)
		}
	})

	for i, p := range m.alive {
		spawn := m.arena.Spawns[i]
		ctx := context.Background()
		handle, err := player.Call(ctx, p.H(), func(tx *world.Tx, pl *player.Player) (*world.EntityHandle, error) {
			h := tx.RemoveEntity(pl)
			return h, nil
		})
		if err != nil {
			continue
		}
		w.Do(func(tx *world.Tx) {
			np := tx.AddEntity(handle).(*player.Player)
			np.Inventory().Clear()
			np.Armour().Inventory().Clear()
			np.Teleport(spawn)
			np.SetGameMode(world.GameModeSurvival)
			np.SetMaxHealth(20)
			np.Heal(20, effect.InstantHealingSource{})
			np.SetFood(20)
			np.Handle(&handler{match: m})
			m.mu.RLock()
			remaining := time.Until(m.refillTime)
			m.mu.RUnlock()
			framework.SendMatchScoreboard(np, 0, len(m.alive), remaining)
		})
	}

	go m.scoreboardTicker()
	m.countdown()
}

func (m *Match) countdown() {
	stopCheck := make(chan struct{})
	defer close(stopCheck)

	
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()
		for {
			select {
			case <-stopCheck:
				return
			case <-ticker.C:
				if m.State() == framework.MatchStateClosed {
					return
				}
				for i, p := range m.alive {
					spawn := m.arena.Spawns[i]
					pos := p.Position()
					dx := pos.X() - spawn.X()
					dy := pos.Y() - spawn.Y()
					dz := pos.Z() - spawn.Z()
					if dx < -1.5 || dx > 1.5 || dy < -1.5 || dy > 1.5 || dz < -1.5 || dz > 1.5 {
						s := spawn
						p.Do(func(tx *world.Tx, pl *player.Player) {
							pl.Teleport(s)
						})
					}
				}
			}
		}
	}()

	sendTitle := func(p *player.Player, txt, sub string, dur time.Duration) {
		if txt == "" {
			return
		}
		t := title.New(txt)
		if sub != "" {
			t = t.WithSubtitle(sub)
		}
		t = t.WithFadeInDuration(0).WithDuration(dur).WithFadeOutDuration(time.Millisecond * 200)
		p.SendTitle(t)
	}

	for i := 30; i >= 0; i-- {
		if m.State() == framework.MatchStateClosed {
			return
		}

		txt := ""
		sub := ""
		chat := ""
		snd := ""
		dur := time.Second

		switch i {
		case 30:
			txt = "§6§lSKYWARS"
			sub = "§fStarting in §e30 seconds"
			chat = "§e[SkyWars] §fThe game is starting soon!"
			snd = "random.levelup"
			dur = 5 * time.Second
		case 20:
			snd = "note.pling"
		case 10:
			txt = "§e§l10"
			sub = "§fGet ready!"
			chat = "§c[SkyWars] §fStarting in §c10 seconds!"
			snd = "random.click"
			dur = 3 * time.Second
		case 9, 8, 7, 6:
			txt = fmt.Sprintf("§e§l%d", i)
			snd = "random.click"
		case 5:
			txt = "§c§l5"
			sub = "§fBattle begins soon!"
			chat = "§c[SkyWars] §fPrepare for battle!"
			snd = "note.hat"
			dur = 3 * time.Second
		case 4, 3, 2:
			txt = fmt.Sprintf("§c§l%d", i)
			snd = "note.hat"
		case 1:
			txt = "§c§l1"
			sub = "§fPrepare yourself!"
			snd = "note.hat"
		case 0:
			txt = "§a§lGO!"
			sub = "§fGood luck!"
			chat = "§a[SkyWars] §fThe battle has started!"
			snd = "random.levelup"
			dur = 5 * time.Second
		}

		for _, p := range m.alive {
			sendTitle(p, txt, sub, dur)
			if chat != "" {
				p.Message(chat)
			}
			if snd != "" {
				p.PlaySoundByName(snd, p.Position())
			}
			if i == 0 {
				p.PlaySoundByName("firework.launch", p.Position())
			}
		}

		if i == 0 {
			break
		}
		time.Sleep(time.Second)
	}

	
	w := m.world
	for _, p := range m.alive {
		spawn := m.arena.Spawns[m.spawnIndices[p.UUID()]]
		center := cube.PosFromVec3(spawn)
		w.Do(func(tx *world.Tx) {
			for dx := -3; dx <= 3; dx++ {
				for dy := -3; dy <= 3; dy++ {
					for dz := -3; dz <= 3; dz++ {
						bp := center.Add(cube.Pos{dx, dy, dz})
						if _, ok := tx.Block(bp).(block.Glass); ok {
							tx.SetBlock(bp, block.Air{}, nil)
						}
					}
				}
			}
		})
	}

	m.mu.Lock()
	m.state = framework.MatchStatePlaying
	m.pvpEnabled = true
	m.mu.Unlock()

	m.broadcast("§aGame started! Good luck!")

	go func() {
		time.Sleep(time.Duration(m.config.GameDurationSeconds) * time.Second)
		m.endGame(framework.EndReasonTimeLimit)
	}()
}


func (m *Match) handleDeath(p *player.Player) {
	m.mu.Lock()
	m.removeFromAlive(p)
	m.dead = append(m.dead, p)
	m.spectators = append(m.spectators, p)
	m.totalKills++
	aliveCount := len(m.alive)

	if m.config.DropLootOnDeath {
		m.dropPlayerLoot(p)
	}

	m.mu.Unlock()

	
	p.Inventory().Clear()
	p.Armour().Inventory().Clear()
	p.SetGameMode(ghostMode{})
	p.StartFlying()
	p.SetInvisible()
	p.SetNameTag("§7[SPECTATOR] " + p.Name())

	
	compass := item.NewStack(item.Compass{}, 1)
	p.Inventory().SetItem(4, compass)

	m.broadcast(fmt.Sprintf("§c%s has been eliminated! (%d remaining)", p.Name(), aliveCount-1))

	
	if aliveCount <= 1 {
		m.endGame(framework.EndReasonLastStanding)
	}
}


func (m *Match) handleQuit(p *player.Player) {
	m.mu.Lock()
	m.removeFromAlive(p)
	m.removeFromSpectators(p)
	aliveCount := len(m.alive)
	m.removeFromPlayers(p)
	m.mu.Unlock()

	m.broadcast(fmt.Sprintf("§e%s has left the game. (%d remaining)", p.Name(), aliveCount))

	if aliveCount <= 1 {
		m.endGame(framework.EndReasonLastStanding)
	}
}

func (m *Match) removeFromAlive(p *player.Player) {
	for i, ap := range m.alive {
		if ap == p {
			m.alive = append(m.alive[:i], m.alive[i+1:]...)
			return
		}
	}
}

func (m *Match) removeFromSpectators(p *player.Player) {
	for i, sp := range m.spectators {
		if sp == p {
			m.spectators = append(m.spectators[:i], m.spectators[i+1:]...)
			return
		}
	}
}

func (m *Match) removeFromPlayers(p *player.Player) {
	for i, mp := range m.players {
		if mp == p {
			m.players = append(m.players[:i], m.players[i+1:]...)
			return
		}
	}
}

func (m *Match) dropPlayerLoot(p *player.Player) {
	inv := p.Inventory()
	pos := p.Position()
	tx := p.Tx()
	if tx == nil {
		return
	}
	for _, slot := range inv.Slots() {
		if !slot.Empty() {
			dropItem(tx, slot, pos)
		}
	}
	inv.Clear()
}

func dropItem(tx *world.Tx, s item.Stack, pos mgl64.Vec3) {
	handle := entity.NewItem(world.EntitySpawnOpts{Position: pos}, s)
	tx.AddEntity(handle)
}


func (m *Match) endGame(reason framework.EndReason) {
	m.mu.Lock()
	if m.state == framework.MatchStateFinished || m.state == framework.MatchStateClosed {
		m.mu.Unlock()
		return
	}
	m.state = framework.MatchStateFinished
	m.mu.Unlock()

	switch reason {
	case framework.EndReasonLastStanding:
		if len(m.alive) == 1 {
			winner := m.alive[0]
			m.broadcast(fmt.Sprintf("§6§l%s wins the game!", winner.Name()))
			m.fw.Stats().IncrementWins(winner.UUID())
		}
	case framework.EndReasonTimeLimit:
		m.broadcast("§eTime's up! The game ends in a draw.")
		m.mu.Lock()
		for _, p := range m.alive {
			m.fw.Stats().IncrementWins(p.UUID())
		}
		m.mu.Unlock()
		reason = framework.EndReasonDraw
	default:
		m.broadcast("§cGame ended.")
	}

	
	go func() {
		time.Sleep(5 * time.Second)

		m.mu.Lock()
		m.state = framework.MatchStateClosed
		for _, p := range m.players {
			m.fw.Players().ReturnToLobby(p)
		}
		m.mu.Unlock()

		
		m.fw.Worlds().ReturnWorld(m.arena.WorldName)

		close(m.done)
		_ = reason
	}()
}


func (m *Match) openDeathMenu(p *player.Player) {
	menu := deathMenu{match: m}
	f := form.NewMenu(menu, "§6You Died!").WithBody("§fChoose an option below:")
	f = f.WithButtons(
		form.Button{Text: "§6Spectate"},
		form.Button{Text: "§6Play Again"},
		form.Button{Text: "§6Lobby"},
	)
	p.SendForm(f)
}


type deathMenu struct {
	match *Match
}

func (d deathMenu) Submit(submitter form.Submitter, pressed form.Button, _ *world.Tx) {
	p := submitter.(*player.Player)
	p.CloseForm()
	switch pressed.Text {
	case "§6Spectate":
	case "§6Play Again":
		d.match.fw.Players().ReturnToLobby(p)
		go func() {
			time.Sleep(time.Second)
			_ = d.match.fw.Games().JoinGame(p, d.match.GameID())
		}()
	case "§6Lobby":
		d.match.fw.Players().ReturnToLobby(p)
	}
}

func (d deathMenu) Close(submitter form.Submitter, _ *world.Tx) {
	submitter.CloseForm()
}


func (m *Match) scoreboardTicker() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-m.done:
			return
		case <-ticker.C:
			if m.State() == framework.MatchStateClosed || m.State() == framework.MatchStateFinished {
				return
			}
			m.mu.RLock()
			alive := make([]*player.Player, len(m.alive))
			copy(alive, m.alive)
			allPlayers := make([]*player.Player, len(m.players))
			copy(allPlayers, m.players)
			kills := m.totalKills
			aliveCount := len(alive)
			remaining := time.Until(m.refillTime)
			m.mu.RUnlock()

			for _, p := range allPlayers {
				framework.SendMatchScoreboard(p, kills, aliveCount, remaining)
			}
		}
	}
}


func (m *Match) broadcast(msg string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.players {
		p.Message(msg)
	}
}


func (m *Match) Done() <-chan struct{} {
	return m.done
}
