package entity

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/world"
)

func mobIdent(name string) string { return "minecraft:" + name }

func mobSpawnEggItemName(entityIdent string) string {
	name := entityIdent[10:]
	if len(name) > 3 && name[len(name)-3:] == "_v2" {
		name = name[:len(name)-3]
	}
	return "minecraft:" + name + "_spawn_egg"
}

func init() {
	for _, t := range MobTypes() {
		t := t
		name := t.encode[10:]
		world.RegisterSpawnEggHandler(name, func(tx *world.Tx, pos cube.Pos) bool {
			handle := world.EntitySpawnOpts{Position: pos.Vec3Centre()}.New(t, MobBehaviourConfig)
			tx.AddEntity(handle)
			return true
		})
		world.RegisterItem(item.SpawnEgg{ItemName: mobSpawnEggItemName(t.encode)})
	}
}

func mobBBOX(w, h float64) cube.BBox {
	half := w / 2
	return cube.Box(-half, 0, -half, half, h, half)
}

var MobBehaviourConfig = PassiveBehaviourConfig{Gravity: 0.08, Drag: 0.02}

type MobType struct {
	encode string
	bbox   cube.BBox
}

func (t MobType) Open(tx *world.Tx, handle *world.EntityHandle, data *world.EntityData) world.Entity {
	return &Mob{Ent: Open(tx, handle, data)}
}

func (t MobType) EncodeEntity() string          { return t.encode }
func (t MobType) BBox(world.Entity) cube.BBox   { return t.bbox }
func (t MobType) DecodeNBT(_ map[string]any, data *world.EntityData) {
	data.Data = MobBehaviourConfig.New()
}
func (t MobType) EncodeNBT(_ *world.EntityData) map[string]any { return nil }

type Mob struct{ *Ent }

var (
	AllayType           = MobType{encode: mobIdent("allay"), bbox: mobBBOX(0.35, 0.6)}
	ArmadilloType       = MobType{encode: mobIdent("armadillo"), bbox: mobBBOX(0.7, 0.65)}
	AxolotlType         = MobType{encode: mobIdent("axolotl"), bbox: mobBBOX(0.75, 0.42)}
	BatType             = MobType{encode: mobIdent("bat"), bbox: mobBBOX(0.5, 0.9)}
	BeeType             = MobType{encode: mobIdent("bee"), bbox: mobBBOX(0.55, 0.5)}
	BlazeType           = MobType{encode: mobIdent("blaze"), bbox: mobBBOX(0.6, 1.8)}
	BoggedType          = MobType{encode: mobIdent("bogged"), bbox: mobBBOX(0.6, 1.99)}
	BreezeType          = MobType{encode: mobIdent("breeze"), bbox: mobBBOX(0.6, 1.7)}
	CamelType           = MobType{encode: mobIdent("camel"), bbox: mobBBOX(1.7, 2.375)}
	CatType             = MobType{encode: mobIdent("cat"), bbox: mobBBOX(0.6, 0.7)}
	CaveSpiderType      = MobType{encode: mobIdent("cave_spider"), bbox: mobBBOX(0.7, 0.5)}
	ChickenType         = MobType{encode: mobIdent("chicken"), bbox: mobBBOX(0.4, 0.7)}
	CodType             = MobType{encode: mobIdent("cod"), bbox: mobBBOX(0.5, 0.3)}
	CowType             = MobType{encode: mobIdent("cow"), bbox: mobBBOX(0.9, 1.3)}
	CreakingType        = MobType{encode: mobIdent("creaking"), bbox: mobBBOX(0.6, 2.0)}
	CreeperType         = MobType{encode: mobIdent("creeper"), bbox: mobBBOX(0.6, 1.7)}
	DolphinType         = MobType{encode: mobIdent("dolphin"), bbox: mobBBOX(0.9, 0.6)}
	DonkeyType          = MobType{encode: mobIdent("donkey"), bbox: mobBBOX(1.4, 1.6)}
	DrownedType         = MobType{encode: mobIdent("drowned"), bbox: mobBBOX(0.6, 1.8)}
	ElderGuardianType   = MobType{encode: mobIdent("elder_guardian"), bbox: mobBBOX(1.9975, 1.9975)}
	EnderDragonType     = MobType{encode: mobIdent("ender_dragon"), bbox: mobBBOX(16.0, 8.0)}
	EndermanType        = MobType{encode: mobIdent("enderman"), bbox: mobBBOX(0.6, 2.9)}
	EndermiteType       = MobType{encode: mobIdent("endermite"), bbox: mobBBOX(0.4, 0.3)}
	EvokerType          = MobType{encode: mobIdent("evoker"), bbox: mobBBOX(0.6, 1.95)}
	FoxType             = MobType{encode: mobIdent("fox"), bbox: mobBBOX(0.6, 0.7)}
	FrogType            = MobType{encode: mobIdent("frog"), bbox: mobBBOX(0.5, 0.5)}
	GhastType           = MobType{encode: mobIdent("ghast"), bbox: mobBBOX(4.0, 4.0)}
	GlowSquidType       = MobType{encode: mobIdent("glow_squid"), bbox: mobBBOX(0.95, 0.95)}
	GoatType            = MobType{encode: mobIdent("goat"), bbox: mobBBOX(0.9, 1.3)}
	GuardianType        = MobType{encode: mobIdent("guardian"), bbox: mobBBOX(0.85, 0.85)}
	HoglinType          = MobType{encode: mobIdent("hoglin"), bbox: mobBBOX(1.3965, 1.3965)}
	HorseType           = MobType{encode: mobIdent("horse"), bbox: mobBBOX(1.4, 1.6)}
	HuskType            = MobType{encode: mobIdent("husk"), bbox: mobBBOX(0.6, 1.8)}
	IronGolemType       = MobType{encode: mobIdent("iron_golem"), bbox: mobBBOX(1.4, 2.9)}
	LlamaType           = MobType{encode: mobIdent("llama"), bbox: mobBBOX(0.9, 1.87)}
	MagmaCubeType       = MobType{encode: mobIdent("magma_cube"), bbox: mobBBOX(2.04, 2.04)}
	MooshroomType       = MobType{encode: mobIdent("mooshroom"), bbox: mobBBOX(0.9, 1.3)}
	MuleType            = MobType{encode: mobIdent("mule"), bbox: mobBBOX(1.4, 1.6)}
	OcelotType          = MobType{encode: mobIdent("ocelot"), bbox: mobBBOX(0.6, 0.7)}
	PandaType           = MobType{encode: mobIdent("panda"), bbox: mobBBOX(1.3, 1.25)}
	ParrotType          = MobType{encode: mobIdent("parrot"), bbox: mobBBOX(0.5, 0.9)}
	PhantomType         = MobType{encode: mobIdent("phantom"), bbox: mobBBOX(0.9, 0.5)}
	PigType             = MobType{encode: mobIdent("pig"), bbox: mobBBOX(0.9, 0.9)}
	PiglinBruteType     = MobType{encode: mobIdent("piglin_brute"), bbox: mobBBOX(0.6, 1.95)}
	PiglinType          = MobType{encode: mobIdent("piglin"), bbox: mobBBOX(0.6, 1.95)}
	PillagerType        = MobType{encode: mobIdent("pillager"), bbox: mobBBOX(0.6, 1.95)}
	PolarBearType       = MobType{encode: mobIdent("polar_bear"), bbox: mobBBOX(1.3, 1.4)}
	PufferfishType      = MobType{encode: mobIdent("pufferfish"), bbox: mobBBOX(0.35, 0.35)}
	RabbitType          = MobType{encode: mobIdent("rabbit"), bbox: mobBBOX(0.4, 0.5)}
	RavagerType         = MobType{encode: mobIdent("ravager"), bbox: mobBBOX(1.95, 2.2)}
	SalmonType          = MobType{encode: mobIdent("salmon"), bbox: mobBBOX(0.7, 0.4)}
	SheepType           = MobType{encode: mobIdent("sheep"), bbox: mobBBOX(0.9, 1.3)}
	ShulkerType         = MobType{encode: mobIdent("shulker"), bbox: mobBBOX(1.0, 1.0)}
	SilverfishType      = MobType{encode: mobIdent("silverfish"), bbox: mobBBOX(0.4, 0.3)}
	SkeletonHorseType   = MobType{encode: mobIdent("skeleton_horse"), bbox: mobBBOX(1.4, 1.6)}
	SkeletonType        = MobType{encode: mobIdent("skeleton"), bbox: mobBBOX(0.6, 1.99)}
	SlimeType           = MobType{encode: mobIdent("slime"), bbox: mobBBOX(2.04, 2.04)}
	SnifferType         = MobType{encode: mobIdent("sniffer"), bbox: mobBBOX(1.9, 1.75)}
	SnowGolemType       = MobType{encode: mobIdent("snow_golem"), bbox: mobBBOX(0.7, 1.9)}
	SpiderType          = MobType{encode: mobIdent("spider"), bbox: mobBBOX(1.4, 0.9)}
	SquidType           = MobType{encode: mobIdent("squid"), bbox: mobBBOX(0.95, 0.95)}
	StrayType           = MobType{encode: mobIdent("stray"), bbox: mobBBOX(0.6, 1.99)}
	StriderType         = MobType{encode: mobIdent("strider"), bbox: mobBBOX(0.9, 1.7)}
	TadpoleType         = MobType{encode: mobIdent("tadpole"), bbox: mobBBOX(0.35, 0.3)}
	TraderLlamaType     = MobType{encode: mobIdent("trader_llama"), bbox: mobBBOX(0.9, 1.87)}
	TropicalFishType    = MobType{encode: mobIdent("tropical_fish"), bbox: mobBBOX(0.5, 0.4)}
	TurtleType          = MobType{encode: mobIdent("turtle"), bbox: mobBBOX(1.2, 0.4)}
	VexType             = MobType{encode: mobIdent("vex"), bbox: mobBBOX(0.4, 0.8)}
	VillagerType        = MobType{encode: mobIdent("villager_v2"), bbox: mobBBOX(0.6, 1.95)}
	VindicatorType      = MobType{encode: mobIdent("vindicator"), bbox: mobBBOX(0.6, 1.95)}
	WanderingTraderType = MobType{encode: mobIdent("wandering_trader"), bbox: mobBBOX(0.6, 1.95)}
	WardenType          = MobType{encode: mobIdent("warden"), bbox: mobBBOX(0.9, 2.9)}
	WitchType           = MobType{encode: mobIdent("witch"), bbox: mobBBOX(0.6, 1.95)}
	WitherType          = MobType{encode: mobIdent("wither"), bbox: mobBBOX(0.9, 3.5)}
	WitherSkeletonType  = MobType{encode: mobIdent("wither_skeleton"), bbox: mobBBOX(0.7, 2.4)}
	WolfType            = MobType{encode: mobIdent("wolf"), bbox: mobBBOX(0.6, 0.85)}
	ZoglinType          = MobType{encode: mobIdent("zoglin"), bbox: mobBBOX(1.3965, 1.3965)}
	ZombieHorseType     = MobType{encode: mobIdent("zombie_horse"), bbox: mobBBOX(1.4, 1.6)}
	ZombiePigmanType    = MobType{encode: mobIdent("zombie_pigman"), bbox: mobBBOX(0.6, 1.8)}
	ZombieType          = MobType{encode: mobIdent("zombie"), bbox: mobBBOX(0.6, 1.8)}
	ZombieVillagerType  = MobType{encode: mobIdent("zombie_villager_v2"), bbox: mobBBOX(0.6, 1.8)}
)

func MobTypes() []MobType {
	return []MobType{
		AllayType, ArmadilloType, AxolotlType, BatType, BeeType, BlazeType,
		BoggedType, BreezeType, CamelType, CatType, CaveSpiderType,
		ChickenType, CodType, CowType, CreakingType, CreeperType,
		DolphinType, DonkeyType, DrownedType, ElderGuardianType,
		EnderDragonType, EndermanType, EndermiteType, EvokerType,
		FoxType, FrogType, GhastType, GlowSquidType, GoatType,
		GuardianType, HoglinType, HorseType, HuskType, IronGolemType,
		LlamaType, MagmaCubeType, MooshroomType, MuleType, OcelotType,
		PandaType, ParrotType, PhantomType, PigType, PiglinBruteType,
		PiglinType, PillagerType, PolarBearType, PufferfishType,
		RabbitType, RavagerType, SalmonType, SheepType, ShulkerType,
		SilverfishType, SkeletonHorseType, SkeletonType, SlimeType,
		SnifferType, SnowGolemType, SpiderType, SquidType, StrayType,
		StriderType, TadpoleType, TraderLlamaType, TropicalFishType,
		TurtleType, VexType, VillagerType, VindicatorType,
		WanderingTraderType, WardenType, WitchType, WitherType,
		WitherSkeletonType, WolfType, ZoglinType, ZombieHorseType,
		ZombiePigmanType, ZombieType, ZombieVillagerType,
	}
}

var mobByName map[string]MobType

func init() {
	mobByName = make(map[string]MobType, len(MobTypes()))
	for _, t := range MobTypes() {
		mobByName[t.encode[10:]] = t
	}
}

func MobByName(name string) (MobType, bool) {
	t, ok := mobByName[name]
	return t, ok
}

func mobEntityTypes() []world.EntityType {
	all := MobTypes()
	out := make([]world.EntityType, len(all))
	for i, t := range all {
		out[i] = t
	}
	return out
}
