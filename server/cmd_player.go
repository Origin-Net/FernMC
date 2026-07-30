package server

import (
	"strconv"
	"strings"

	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type giveCmd struct {
	Player []cmd.Target
	Item   ItemName
	Amount cmd.Optional[int]
}

func (g giveCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	name := string(g.Item)
	if strings.HasPrefix(name, "minecraft:") {
		name = name[10:]
	}
	it, ok := world.ItemByName(name, 0)
	if !ok {
		o.Error("Unknown item: " + string(g.Item))
		return
	}
	count := 1
	if a, ok := g.Amount.Load(); ok && a > 0 {
		count = a
	}
	stack := item.NewStack(it, count)
	for _, t := range g.Player {
		p, ok := t.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		p.Inventory().AddItem(stack)
		o.Printf("Given %s * %d to %s", name, count, p.Name())
	}
}

type clearCmd struct {
	Player   cmd.Optional[[]cmd.Target]
	Item     cmd.Optional[ItemName]
	MaxCount cmd.Optional[int]
}

func (c clearCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	players, explicit := c.Player.Load()
	if !explicit {
		p, ok := src.(*player.Player)
		if !ok {
			o.Error("Please specify a player")
			return
		}
		players = []cmd.Target{p}
	}
	for _, t := range players {
		p, ok := t.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		if name, hasItem := c.Item.Load(); hasItem {
			n := strings.TrimPrefix(string(name), "minecraft:")
			it, ok := world.ItemByName(n, 0)
			if !ok {
				o.Error("Unknown item: " + string(name))
				return
			}
			count, _ := c.MaxCount.Load()
			_ = p.Inventory().RemoveItem(item.NewStack(it, count))
			o.Printf("Cleared %d %s from %s", count, n, p.Name())
		} else {
			p.Inventory().Clear()
			o.Printf("Cleared the inventory of %s, removing all items", p.Name())
		}
	}
}

type xpCmd struct {
	Amount string
	Player cmd.Optional[[]cmd.Target]
}

func (x xpCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	players, explicit := x.Player.Load()
	if !explicit {
		p, ok := src.(*player.Player)
		if !ok {
			o.Error("Please specify a player")
			return
		}
		players = []cmd.Target{p}
	}
	isLevels := strings.HasSuffix(strings.ToUpper(x.Amount), "L")
	amountStr := strings.TrimSuffix(strings.TrimSuffix(x.Amount, "L"), "l")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		o.Error("Invalid amount: " + x.Amount)
		return
	}
	for _, t := range players {
		p, ok := t.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		if isLevels {
			if amount >= 0 {
				p.AddExperienceLevels(amount)
				o.Printf("Gave %d experience levels to %s", amount, p.Name())
			} else {
				p.AddExperienceLevels(amount)
				o.Printf("Taken %d levels from %s", -amount, p.Name())
			}
		} else {
			if amount >= 0 {
				p.AddExperience(amount)
				o.Printf("Gave %d experience to %s", amount, p.Name())
			} else {
				p.RemoveExperience(-amount)
				o.Printf("Taken %d experience from %s", -amount, p.Name())
			}
		}
	}
}

type tpCmd struct {
	Targets     []cmd.Target
	Destination cmd.Optional[mgl64.Vec3]
	Target      cmd.Optional[[]cmd.Target]
}

func (t tpCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	teleport := func(target cmd.Target, pos mgl64.Vec3) {
		switch e := target.(type) {
		case *player.Player:
			e.Teleport(pos)
		case interface{ Teleport(mgl64.Vec3) }:
			e.Teleport(pos)
		default:
			o.Error("Target cannot be teleported")
		}
	}
	if dest, ok := t.Destination.Load(); ok {
		for _, target := range t.Targets {
			teleport(target, dest)
			o.Printf("Teleported %s to %.1f %.1f %.1f", targetName(target), dest[0], dest[1], dest[2])
		}
		return
	}
	if target, ok := t.Target.Load(); ok && len(target) > 0 {
		pos := target[0].Position()
		for _, target := range t.Targets {
			teleport(target, pos)
		}
		o.Printf("Teleported %s to %s", targetName(t.Targets[0]), targetName(target[0]))
	}
}

type spawnpointCmd struct {
	Player cmd.Optional[[]cmd.Target]
	Pos    cmd.Optional[mgl64.Vec3]
}

func (s spawnpointCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	players, explicit := s.Player.Load()
	if !explicit {
		p, ok := src.(*player.Player)
		if !ok {
			o.Error("Please specify a player")
			return
		}
		players = []cmd.Target{p}
	}
	pos := src.Position()
	if p, ok := s.Pos.Load(); ok {
		pos = p
	}
	for _, t := range players {
		p, ok := t.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		tx.World().SetPlayerSpawn(p.UUID(), cube.PosFromVec3(pos))
		o.Printf("Set %s's spawn to %.1f %.1f %.1f", p.Name(), pos[0], pos[1], pos[2])
	}
}

type setworldspawnCmd struct {
	Pos cmd.Optional[mgl64.Vec3]
}

func (s setworldspawnCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	w := tx.World()
	pos := src.Position()
	if p, ok := s.Pos.Load(); ok {
		pos = p
	}
	w.SetSpawn(cube.PosFromVec3(pos))
	o.Printf("Set the world spawn point to (%.1f, %.1f, %.1f)", pos[0], pos[1], pos[2])
}

func init() {
	cmd.Register(cmd.New("give", "Gives an item to a player", nil, giveCmd{}))
	cmd.Register(cmd.New("clear", "Clears items from player inventory", nil, clearCmd{}))
	cmd.Register(cmd.New("xp", "Adds or removes player experience", nil, xpCmd{}))
	cmd.Register(cmd.New("teleport", "Teleports entities to specific locations", []string{"tp"}, tpCmd{}))
	cmd.Register(cmd.New("spawnpoint", "Sets the spawn point for a player", nil, spawnpointCmd{}))
	cmd.Register(cmd.New("setworldspawn", "Sets the world spawn location", nil, setworldspawnCmd{}))
}
