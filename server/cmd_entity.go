package server

import (
	"strings"
	"time"

	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/entity/effect"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/world"
)

type killDamageSource struct{}

func (killDamageSource) ReducedByArmour() bool     { return false }
func (killDamageSource) ReducedByResistance() bool { return false }
func (killDamageSource) Fire() bool                { return false }
func (killDamageSource) IgnoreTotem() bool         { return true }

type cmdDamageSource struct{}

func (cmdDamageSource) ReducedByArmour() bool     { return true }
func (cmdDamageSource) ReducedByResistance() bool { return true }
func (cmdDamageSource) Fire() bool                { return false }
func (cmdDamageSource) IgnoreTotem() bool         { return false }

type killCmd struct {
	Target cmd.Optional[[]cmd.Target]
}

func (k killCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	targets, explicit := k.Target.Load()
	if !explicit {
		p, ok := src.(*player.Player)
		if !ok {
			o.Error("Please specify a target")
			return
		}
		targets = []cmd.Target{p}
	}
	for _, t := range targets {
		if p, ok := t.(*player.Player); ok {
			p.Hurt(p.MaxHealth(), killDamageSource{})
			o.Printf("Killed %s", p.Name())
		} else {
			o.Printf("Killed %s", targetName(t))
		}
	}
}


type EffectEnum string

func (EffectEnum) Type() string { return "Effect" }
func (EffectEnum) Options(cmd.Source) []string {
	return []string{
		"clear",
		"speed", "slowness", "haste", "mining_fatigue",
		"strength", "jump_boost", "nausea", "regeneration",
		"resistance", "fire_resistance", "water_breathing", "invisibility",
		"blindness", "night_vision", "hunger", "weakness",
		"poison", "wither", "health_boost", "absorption",
		"saturation", "levitation", "slow_falling", "conduit_power",
		"fatal_poison", "darkness",
	}
}

type effectCmd struct {
	Player        []cmd.Target
	Effect        EffectEnum
	Seconds       cmd.Optional[int]
	Amplifier     cmd.Optional[int]
	HideParticles cmd.Optional[Bool]
}

func (e effectCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	if strings.ToLower(string(e.Effect)) == "clear" {
		for _, t := range e.Player {
			p, ok := t.(*player.Player)
			if !ok {
				continue
			}
			for _, eff := range p.Effects() {
				p.RemoveEffect(eff.Type())
			}
			o.Printf("Took all effects from %s", p.Name())
		}
		return
	}
	for _, t := range e.Player {
		p, ok := t.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		sec := 30
		if s, ok := e.Seconds.Load(); ok && s > 0 {
			sec = s
		}
		amp := 0
		if a, ok := e.Amplifier.Load(); ok {
			amp = a
		}
		eff, ok := lastingEffectByName(string(e.Effect), amp, time.Duration(sec)*time.Second)
		if !ok {
			o.Error("Unknown effect: " + string(e.Effect))
			return
		}
		p.AddEffect(eff)
		o.Printf("Given %s * %d to %s for %d seconds", string(e.Effect), amp, p.Name(), sec)
	}
}

type enchantCmd struct {
	Player  []cmd.Target
	Enchant Enchant
	Level   cmd.Optional[int]
}

func (e enchantCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	ench := strings.ToLower(string(e.Enchant))
	ench = strings.ReplaceAll(ench, "_", " ")
	var found item.EnchantmentType
	for _, et := range item.Enchantments() {
		if strings.EqualFold(et.Name(), ench) {
			found = et
			break
		}
	}
	if found == nil {
		o.Error("Unknown enchantment: " + string(e.Enchant))
		return
	}
	lvl := 1
	if l, ok := e.Level.Load(); ok && l > 0 {
		lvl = l
	}
	if lvl > found.MaxLevel() {
		lvl = found.MaxLevel()
	}
	enc := item.NewEnchantment(found, lvl)
	for _, t := range e.Player {
		p, ok := t.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		main, off := p.HeldItems()
		_ = off
		p.SetHeldItems(main.WithEnchantments(enc), item.Stack{})
		o.Printf("Enchanting succeeded for %s", p.Name())
	}
}

type DamageCauseEnum string

func (DamageCauseEnum) Type() string { return "DamageCause" }
func (DamageCauseEnum) Options(cmd.Source) []string {
	return []string{
		"all", "anvil", "block_explosion", "charging", "contact", "drowning",
		"entity_attack", "entity_explosion", "fall", "falling_block", "fire",
		"fire_tick", "fly_into_wall", "freezing", "lava", "lightning", "magic",
		"magma", "none", "override", "piston", "projectile", "sonic_boom",
		"stalactite", "starvation", "suffocation", "suicide", "temperature",
		"thorns", "void", "wither",
	}
}

type damageCmd struct {
	Target []cmd.Target
	Amount int
	Cause  cmd.Optional[DamageCauseEnum]
}

func (d damageCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	for _, t := range d.Target {
		p, ok := t.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		p.Hurt(float64(d.Amount), cmdDamageSource{})
		o.Printf("Dealt %d damage to %s", d.Amount, p.Name())
	}
}

func lastingEffectByName(name string, lvl int, dur time.Duration) (effect.Effect, bool) {
	t := lastingEffectType(name)
	if t == nil {
		return effect.Effect{}, false
	}
	return effect.New(t, lvl, dur), true
}

func lastingEffectType(name string) effect.LastingType {
	switch strings.ToLower(name) {
	case "speed":
		return effect.Speed
	case "slowness":
		return effect.Slowness
	case "haste":
		return effect.Haste
	case "mining_fatigue":
		return effect.MiningFatigue
	case "strength":
		return effect.Strength
	case "jump_boost":
		return effect.JumpBoost
	case "nausea":
		return effect.Nausea
	case "regeneration":
		return effect.Regeneration
	case "resistance":
		return effect.Resistance
	case "fire_resistance":
		return effect.FireResistance
	case "water_breathing":
		return effect.WaterBreathing
	case "invisibility":
		return effect.Invisibility
	case "blindness":
		return effect.Blindness
	case "night_vision":
		return effect.NightVision
	case "hunger":
		return effect.Hunger
	case "weakness":
		return effect.Weakness
	case "poison":
		return effect.Poison
	case "wither":
		return effect.Wither
	case "health_boost":
		return effect.HealthBoost
	case "absorption":
		return effect.Absorption
	case "saturation":
		return effect.Saturation
	case "levitation":
		return effect.Levitation
	case "slow_falling":
		return effect.SlowFalling
	case "conduit_power":
		return effect.ConduitPower
	case "fatal_poison":
		return effect.FatalPoison
	case "darkness":
		return effect.Darkness
	}
	return nil
}

func init() {
	cmd.Register(cmd.New("kill", "Kills entities like players and mobs", []string{"suicide"}, killCmd{}))
	cmd.Register(cmd.New("effect", "Add or clear status effects", nil, effectCmd{}))
	cmd.Register(cmd.New("enchant", "Adds an enchantment to a player's selected item", nil, enchantCmd{}))
	cmd.Register(cmd.New("damage", "Apply damage to the specified entities", nil, damageCmd{}))
}
