package server

import (
	"strings"

	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/entity"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/world"
)

type BlockName string

func (BlockName) Type() string { return "BlockType" }
func (BlockName) Options(cmd.Source) []string {
	blocks := world.Blocks()
	names := make([]string, 0, len(blocks))
	for _, b := range blocks {
		name, _ := b.EncodeBlock()
		name = strings.TrimPrefix(name, "minecraft:")
		names = append(names, name)
	}
	return names
}

type ItemName string

func (ItemName) Type() string { return "ItemType" }
func (ItemName) Options(cmd.Source) []string {
	items := world.Items()
	names := make([]string, 0, len(items))
	seen := map[string]bool{}
	for _, it := range items {
		name, _ := it.EncodeItem()
		name = strings.TrimPrefix(name, "minecraft:")
		if !seen[name] {
			names = append(names, name)
			seen[name] = true
		}
	}
	return names
}

type MobName string

func (MobName) Type() string { return "EntityType" }
func (MobName) Options(cmd.Source) []string {
	mobs := entity.MobTypes()
	n := make([]string, len(mobs))
	for i, m := range mobs {
		n[i] = strings.TrimPrefix(m.EncodeEntity(), "minecraft:")
	}
	return n
}

type CmdName string

func (CmdName) Type() string { return "CmdName" }
func (CmdName) Options(src cmd.Source) []string {
	cmds := cmd.Commands()
	names := make([]string, 0, len(cmds))
	for n := range cmds {
		names = append(names, n)
	}
	return names
}

type Enchant string

func (Enchant) Type() string { return "EnchantType" }
func (Enchant) Options(cmd.Source) []string {
	enchants := item.Enchantments()
	names := make([]string, len(enchants))
	for i, e := range enchants {
		names[i] = strings.ReplaceAll(strings.ToLower(e.Name()), " ", "_")
	}
	return names
}

type GameRule string

func (GameRule) Type() string { return "GameRule" }
func (GameRule) Options(cmd.Source) []string {
	return []string{
		"commandblockoutput", "commandblocksenabled",
		"dodaylightcycle", "doentitydrops", "dofiretick",
		"doimmediaterespawn", "doinsomnia", "dolimitedcrafting",
		"domobloot", "domobspawning", "dotiledrops",
		"doweathercycle", "drowningdamage", "falldamage",
		"firedamage", "freezedamage", "functioncommandlimit",
		"keepinventory", "maxcommandchainlength", "mobgriefing",
		"naturalregeneration", "playerssleepingpercentage",
		"projectilescanbreakblocks", "pvp", "randomtickspeed",
		"recipesunlock", "respawnblocksexplode",
		"sendcommandfeedback", "showbordereffect", "showcoordinates",
		"showdaysplayed", "showdeathmessages", "showrecipemessages",
		"showtags", "spawnradius", "tntexplodes",
		"tntexplosiondropdecay",
	}
}

type Bool string

func (Bool) Type() string { return "Bool" }
func (Bool) Options(cmd.Source) []string {
	return []string{"true", "false"}
}

type OldBlockHandling string

func (OldBlockHandling) Type() string { return "OldBlockHandling" }
func (OldBlockHandling) Options(cmd.Source) []string {
	return []string{"destroy", "replace", "keep"}
}

type MaskMode string

func (MaskMode) Type() string { return "MaskMode" }
func (MaskMode) Options(cmd.Source) []string {
	return []string{"masked", "replace", "filtered"}
}

type ParticleType string

func (ParticleType) Type() string { return "ParticleType" }
func (ParticleType) Options(cmd.Source) []string {
	return []string{
		"dragon_breath", "block", "bubble", "critical", "damage_indicator",
		"drip", "droplet", "dripstone_lava", "dripstone_water",
		"electric_spark", "enchanting_table", "end_rod", "explosion",
		"falling_dust", "firework", "flame", "glow",
		"huge_explosion", "icon_crack", "ink", "item_splash",
		"large_explosion", "lava", "lava_drip", "llama_spit",
		"mob_flame", "mob_spell", "mob_spell_instantaneous", "note",
		"potion", "rain_splash", "redstone", "rising_border_dust",
		"sculk_soul", "shulker_bullet", "slime", "smoke",
		"snowball_poof", "soul", "spell", "splash",
		"splash_spell", "spore_blossom_air", "squid_ink", "sweep_attack",
		"terrain", "tint", "town_aura", "vibration",
		"vibrate", "village_anger", "void", "water_drip",
		"water_splash", "water_wake", "wax_off", "wax_on",
		"white_ash", "wind_explosion", "witch_spell",
	}
}

type SoundType string

func (SoundType) Type() string { return "SoundType" }
func (SoundType) Options(cmd.Source) []string {
	return []string{
		"ambient.cave", "ambient.underwater.loop",
		"block.anvil.land", "block.anvil.use", "block.anvil.break",
		"block.beacon.activate", "block.beacon.ambient", "block.beacon.deactivate",
		"block.chest.open", "block.chest.close",
		"block.dispenser.dispense", "block.dispenser.fail", "block.dispenser.launch",
		"block.door.open", "block.door.close",
		"block.fence_gate.open", "block.fence_gate.close",
		"block.fire.ambient", "block.furnace.lit",
		"block.lava.ambient", "block.lava.extinguish",
		"block.portal.ambient", "block.water.ambient",
		"block.trapdoor.open", "block.trapdoor.close",
		"enchant.thorns.hit",
		"entity.arrow.hit", "entity.arrow.shoot",
		"entity.blaze.hurt", "entity.blaze.shoot",
		"entity.cat.ambient", "entity.chicken.ambient", "entity.cow.ambient",
		"entity.creeper.primed", "entity.creeper.death",
		"entity.enderman.ambient", "entity.enderman.death",
		"entity.experience_orb.pickup",
		"entity.firework.rocket.blast",
		"entity.ghast.ambient", "entity.ghast.shoot",
		"entity.item.pickup",
		"entity.player.death", "entity.player.hurt", "entity.player.levelup",
		"entity.sheep.ambient", "entity.skeleton.ambient", "entity.skeleton.shoot",
		"entity.slime.squish",
		"entity.spider.ambient",
		"entity.tnt.primed",
		"entity.villager.ambient",
		"entity.witch.ambient", "entity.witch.hurt",
		"entity.wither.shoot",
		"entity.wolf.ambient", "entity.wolf.death", "entity.wolf.hurt",
		"entity.zombie.ambient", "entity.zombie.death", "entity.zombie.hurt",
		"item.armor.equip_diamond", "item.armor.equip_iron", "item.armor.equip_gold",
		"item.armor.equip_leather", "item.armor.equip_chain",
		"item.axe.strip", "item.flintandsteel.use",
		"item.shield.block", "item.shield.break",
		"item.totem.use",
		"item.trident.throw", "item.trident.hit",
		"jukebox.play", "jukebox.stop",
		"music.game", "music.game.creative", "music.game.end", "music.game.nether", "music.menu",
		"note.harp", "note.bass", "note.bd", "note.snare", "note.hat", "note.pling",
		"portal.ambient", "portal.travel",
		"random.explode", "random.fuse", "random.burp", "random.eat", "random.drink",
		"random.levelup", "random.orb", "random.pop", "random.splash",
		"record.13", "record.cat", "record.chirp", "record.far", "record.mall", "record.mellohi",
		"record.stal", "record.strad", "record.ward", "record.11", "record.wait", "record.pigstep",
		"step.grass", "step.gravel", "step.stone", "step.sand", "step.wood", "step.ladder",
		"weather.rain",
	}
}

type Setting string

func (Setting) Type() string { return "Setting" }
func (Setting) Options(cmd.Source) []string {
	return []string{
		"allow-cheats", "difficulty",
	}
}

type PermissionAction string

func (PermissionAction) Type() string { return "PermissionAction" }
func (PermissionAction) Options(cmd.Source) []string {
	return []string{"add", "remove", "list", "reload"}
}

type Permission string

func (Permission) Type() string { return "Permission" }
func (Permission) Options(cmd.Source) []string {
	return []string{"operator", "member", "visitor"}
}

type EntityEvent string

func (EntityEvent) Type() string { return "EntityEvent" }
func (EntityEvent) Options(cmd.Source) []string {
	return []string{
		"acknowledge_player", "become_aggressive", "be_interested",
		"be_interval_untamed", "be_nervous", "eat_block", "eating",
		"enter_critical", "enter_environment", "idle",
		"idle_angry", "jump", "love", "panic",
		"play_dead", "sneeze", "step_height", "step_height_force",
		"tame", "trust",
	}
}

type MusicAction string

func (MusicAction) Type() string { return "MusicAction" }
func (MusicAction) Options(cmd.Source) []string {
	return []string{"play", "stop", "queue", "volume"}
}

type StructureSaveMode string

func (StructureSaveMode) Type() string { return "StructureSaveMode" }
func (StructureSaveMode) Options(cmd.Source) []string {
	return []string{"memory", "disk"}
}

type StructureRotation string

func (StructureRotation) Type() string { return "StructureRotation" }
func (StructureRotation) Options(cmd.Source) []string {
	return []string{"0", "90", "180", "270"}
}

func srcName(s cmd.Source) string {
	if n, ok := s.(cmd.NamedTarget); ok {
		return n.Name()
	}
	return "Console"
}

func targetName(t cmd.Target) string {
	if n, ok := t.(cmd.NamedTarget); ok {
		return n.Name()
	}
	return "Unknown"
}
