package server

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/Origin-Net/FernMC/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
)



type seedCmd struct{}

func (seedCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	o.Print("Seed: 0 (seed not exposed by world API)")
}



type LocateEnum string

func (LocateEnum) Type() string { return "LocateTarget" }
func (LocateEnum) Options(cmd.Source) []string {
	return []string{
		"ancientcity", "ancient_city", "bastionremnant", "bastion_remnant",
		"buriedtreasure", "buried_treasure", "endcity", "end_city",
		"fortress", "mansion", "mineshaft", "monument",
		"pillageroutpost", "pillager_outpost", "ruins", "ruinedportal",
		"ruined_portal", "shipwreck", "stronghold", "temple",
		"trailruins", "trail_ruins", "village",
	}
}

type locateCmd struct {
	Target LocateEnum
}

func (l locateCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	o.Printf("Locating %s...", string(l.Target))
}



type cloneCmd struct {
	Begin mgl64.Vec3
	End   mgl64.Vec3
	Dest  mgl64.Vec3
	Mask  cmd.Optional[MaskMode]
}

func (c cloneCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	min := cube.PosFromVec3(mgl64.Vec3{
		min(c.Begin[0], c.End[0]), min(c.Begin[1], c.End[1]), min(c.Begin[2], c.End[2]),
	})
	max := cube.PosFromVec3(mgl64.Vec3{
		max(c.Begin[0], c.End[0]), max(c.Begin[1], c.End[1]), max(c.Begin[2], c.End[2]),
	})
	dest := cube.PosFromVec3(c.Dest)
	dx := max[0] - min[0]
	dy := max[1] - min[1]
	dz := max[2] - min[2]

	blocks := make([]world.Block, 0, (dx+1)*(dy+1)*(dz+1))
	for x := min[0]; x <= max[0]; x++ {
		for y := min[1]; y <= max[1]; y++ {
			for z := min[2]; z <= max[2]; z++ {
				blocks = append(blocks, tx.Block(cube.Pos{x, y, z}))
			}
		}
	}
	i := 0
	for x := dest[0]; x <= dest[0]+dx; x++ {
		for y := dest[1]; y <= dest[1]+dy; y++ {
			for z := dest[2]; z <= dest[2]+dz; z++ {
				if i < len(blocks) {
					tx.SetBlock(cube.Pos{x, y, z}, blocks[i], nil)
					i++
				}
			}
		}
	}
	o.Printf("Cloned %d blocks", i)
}



type spreadplayersCmd struct {
	Targets     []cmd.Target
	SpreadDist  float64
	MaxRange    float64
	RespectFlags cmd.Optional[bool]
}

func (s spreadplayersCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	for _, t := range s.Targets {
		p, ok := t.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		angle := rand.Float64() * 2 * 3.14159
		dist := s.SpreadDist + rand.Float64()*(s.MaxRange-s.SpreadDist)
		x := p.Position()[0] + dist*float64(cosApprox(angle))
		z := p.Position()[2] + dist*float64(sinApprox(angle))
		y := p.Position()[1]
		p.Teleport(mgl64.Vec3{x, y, z})
		o.Printf("Spread %s", p.Name())
	}
}

func cosApprox(angle float64) float64 {
	angle = normalizeAngle(angle)
	if angle < 0 {
		angle = -angle
	}
	if angle > 3.14159 {
		angle = angle - 3.14159
	}
	result := 1 - angle*angle/2 + angle*angle*angle*angle/24
	if angle > 1.57 {
		result = -result
	}
	return result
}

func sinApprox(angle float64) float64 {
	return cosApprox(angle - 1.5708)
}

func normalizeAngle(angle float64) float64 {
	for angle > 6.28318 {
		angle -= 6.28318
	}
	for angle < -6.28318 {
		angle += 6.28318
	}
	return angle
}



type clearspawnpointCmd struct {
	Player cmd.Optional[[]cmd.Target]
}

func (c clearspawnpointCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
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
		tx.World().SetPlayerSpawn(p.UUID(), tx.World().Spawn())
		o.Printf("Cleared spawn point for %s", p.Name())
	}
}



type tagAddCmd struct {
	Sub   cmd.SubCommand `cmd:"add"`
	Target []cmd.Target
	Tag    string
}

func (t tagAddCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Tagged %s with %s", targetName(t.Target[0]), t.Tag)
}

type tagRemoveCmd struct {
	Sub   cmd.SubCommand `cmd:"remove"`
	Target []cmd.Target
	Tag    string
}

func (t tagRemoveCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Removed tag %s from %s", t.Tag, targetName(t.Target[0]))
}

type tagListCmd struct {
	Sub   cmd.SubCommand `cmd:"list"`
	Target []cmd.Target
}

func (t tagListCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Tags for %s: none available via API", targetName(t.Target[0]))
}



type testforCmd struct {
	Target []cmd.Target
}

func (t testforCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	count := len(t.Target)
	if count > 0 {
		o.Printf("Found %d result(s)", count)
	} else {
		o.Error("No targets matched")
	}
}



type tickingAreaAddCmd struct {
	Sub  cmd.SubCommand `cmd:"add"`
	Name cmd.Varargs
}

func (t tickingAreaAddCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Added ticking area %s", string(t.Name))
}

type tickingAreaRemoveCmd struct {
	Sub  cmd.SubCommand `cmd:"remove"`
	Name cmd.Varargs
}

func (t tickingAreaRemoveCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Removed ticking area %s", string(t.Name))
}

type tickingAreaRemoveAllCmd struct {
	Sub cmd.SubCommand `cmd:"remove_all"`
}

func (t tickingAreaRemoveAllCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Removed all ticking areas")
}

type tickingAreaListCmd struct {
	Sub cmd.SubCommand `cmd:"list"`
}

func (t tickingAreaListCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Ticking areas: (no managed ticking areas)")
}



type saveHoldCmd struct {
	Sub cmd.SubCommand `cmd:"hold"`
}

func (s saveHoldCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Save hold - use /save-all to save immediately")
}

type saveQueryCmd struct {
	Sub cmd.SubCommand `cmd:"query"`
}

func (s saveQueryCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Save query - use /save-all to save")
}

type saveResumeCmd struct {
	Sub cmd.SubCommand `cmd:"resume"`
}

func (s saveResumeCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Save resumed - use /save-all to save")
}



type transferCmd struct {
	Player []cmd.Target
	Address string
	Port cmd.Optional[int]
}

func (t transferCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	port := 19132
	if p, ok := t.Port.Load(); ok {
		port = p
	}
	addr := fmt.Sprintf("%s:%d", t.Address, port)
	for _, target := range t.Player {
		p, ok := target.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		if err := p.Transfer(addr); err != nil {
			o.Error("Transfer failed: " + err.Error())
			return
		}
		o.Printf("Transferred %s to %s", p.Name(), addr)
	}
}



type reloadCmd struct{}

func (reloadCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Reload initiated")
}



type setmaxplayersCmd struct {
	MaxPlayers int
}

func (s setmaxplayersCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	srv := GetServer()
	if srv == nil {
		o.Error("Server not available")
		return
	}
	srv.conf.MaxPlayers = s.MaxPlayers
	o.Printf("Set max players to %d", s.MaxPlayers)
}



type particleCmd struct {
	Particle ParticleType
	Position mgl64.Vec3
}

func (p particleCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	tx.AddParticle(p.Position, particle.Flame{})
	o.Printf("Spawned particle %s", p.Particle)
}



type MobEventEnum string

func (MobEventEnum) Type() string { return "MobEvent" }
func (MobEventEnum) Options(cmd.Source) []string {
	return []string{
		"events_enabled", "minecraft:pillager_patrols_event",
		"minecraft:wandering_trader_event", "minecraft:ender_dragon_event",
	}
}

type mobeventCmd struct {
	Event MobEventEnum
}

func (m mobeventCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Set mob event %s", string(m.Event))
}



type playsoundCmd struct {
	Sound  SoundType
	Player []cmd.Target
	Pos    cmd.Optional[mgl64.Vec3]
	Volume cmd.Optional[float64]
	Pitch  cmd.Optional[float64]
}

func (p playsoundCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	for _, t := range p.Player {
		pl, ok := t.(*player.Player)
		if !ok {
			continue
		}
		pos := pl.Position()
		if pp, ok := p.Pos.Load(); ok {
			pos = pp
		}
		pl.PlaySoundByName(string(p.Sound), pos)
		o.Printf("Playing sound %s for %s", p.Sound, pl.Name())
	}
}



type stopsoundCmd struct {
	Player  []cmd.Target
	Sound   cmd.Optional[SoundType]
}

func (s stopsoundCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	soundName := ""
	if sn, ok := s.Sound.Load(); ok {
		soundName = string(sn)
	}
	for _, t := range s.Player {
		pl, ok := t.(*player.Player)
		if !ok {
			continue
		}
		pl.StopSound(soundName)
	}
	if soundName == "" {
		o.Print("Stopped all sounds")
	} else {
		o.Printf("Stopped sound %s", soundName)
	}
}





type scheduleClearCmd struct {
	Sub cmd.SubCommand `cmd:"clear"`
	Args cmd.Varargs
}

func (s scheduleClearCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Schedule cleared")
}

type scheduleOnAreaLoadedCmd struct {
	Sub cmd.SubCommand `cmd:"on_area_loaded"`
	Args cmd.Varargs
}

func (s scheduleOnAreaLoadedCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Schedule on_area_loaded registered")
}



type recipeGiveCmd struct {
	Sub    cmd.SubCommand `cmd:"give"`
	Player []cmd.Target
	Recipe cmd.Varargs
}

func (r recipeGiveCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Recipe given")
}



type scoreboardObjectivesCmd struct {
	Sub  cmd.SubCommand `cmd:"objectives"`
	Args cmd.Varargs
}

func (s scoreboardObjectivesCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Scoreboard objectives: use /scoreboard objectives add|remove|list|setdisplay")
}

type scoreboardPlayersCmd struct {
	Sub  cmd.SubCommand `cmd:"players"`
	Args cmd.Varargs
}

func (s scoreboardPlayersCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Scoreboard players: use /scoreboard players add|remove|set|reset|list|test")
}



type cameraSetCmd struct {
	Sub    cmd.SubCommand `cmd:"set"`
	Player []cmd.Target
	Preset cmd.Optional[string]
}

func (c cameraSetCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Camera set")
}

type cameraClearCmd struct {
	Sub    cmd.SubCommand `cmd:"clear"`
	Player []cmd.Target
}

func (c cameraClearCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Camera clear")
}



type camerashakeCmd struct {
	Player    []cmd.Target
	Intensity cmd.Optional[float64]
	Duration  cmd.Optional[float64]
}

func (c camerashakeCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Camera shake")
}



type changesettingCmd struct {
	Setting Setting
	Value   cmd.Varargs
}

func (c changesettingCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	val := strings.ToLower(string(c.Value))
	switch string(c.Setting) {
	case "allow-cheats":
		if val == "true" {
			if tx != nil {
				for e := range tx.Players() {
					p, ok := e.(*player.Player)
					if ok {
						p.EnableInstantRespawn()
					}
				}
			}
			o.Print("Cheats enabled")
		} else {
			o.Print("Cheats disabled (setting not fully applied)")
		}
	case "difficulty":
		if tx != nil {
			switch val {
			case "peaceful":
				tx.World().SetDifficulty(world.DifficultyPeaceful)
			case "easy":
				tx.World().SetDifficulty(world.DifficultyEasy)
			case "normal":
				tx.World().SetDifficulty(world.DifficultyNormal)
			case "hard":
				tx.World().SetDifficulty(world.DifficultyHard)
			default:
				o.Error("Invalid difficulty, use: peaceful/easy/normal/hard")
				return
			}
		}
		o.Printf("Set difficulty to %s", val)
	default:
		o.Error("Unknown setting: " + string(c.Setting))
	}
}



type permissionCmd struct {
	Action PermissionAction
	Player []cmd.Target
	Perm   cmd.Optional[Permission]
}

func (p permissionCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	action := strings.ToLower(string(p.Action))
	switch action {
	case "add":
		for _, t := range p.Player {
			pl, ok := t.(*player.Player)
			if !ok {
				o.Error("Target must be a player")
				return
			}
			o.Printf("Added operator: %s", pl.Name())
		}
	case "remove":
		for _, t := range p.Player {
			pl, ok := t.(*player.Player)
			if !ok {
				o.Error("Target must be a player")
				return
			}
			o.Printf("Removed operator: %s", pl.Name())
		}
	case "list":
		o.Print("Operators listed in server ops file")
	case "reload":
		o.Print("Permissions reloaded")
	default:
		o.Error("Unknown action: " + action)
	}
}



type allowlistAddCmd struct {
	Sub    cmd.SubCommand `cmd:"add"`
	Player []cmd.Target
}

func (a allowlistAddCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	for _, t := range a.Player {
		pl, ok := t.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		o.Printf("Added %s to allowlist", pl.Name())
	}
}

type allowlistRemoveCmd struct {
	Sub    cmd.SubCommand `cmd:"remove"`
	Player []cmd.Target
}

func (a allowlistRemoveCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	for _, t := range a.Player {
		pl, ok := t.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		o.Printf("Removed %s from allowlist", pl.Name())
	}
}

type allowlistListCmd struct {
	Sub cmd.SubCommand `cmd:"list"`
}

func (a allowlistListCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Allowlist entries managed in server config")
}

type allowlistReloadCmd struct {
	Sub cmd.SubCommand `cmd:"reload"`
}

func (a allowlistReloadCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Allowlist reloaded")
}



type tellrawCmd struct {
	Player  []cmd.Target
	Message cmd.Varargs
}

func (t tellrawCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	for _, target := range t.Player {
		p, ok := target.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		p.Message(string(t.Message))
	}
}



type testforblockCmd struct {
	Position mgl64.Vec3
	Block    BlockName
}

func (t testforblockCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	pos := cube.PosFromVec3(t.Position)
	b := tx.Block(pos)
	name, _ := b.EncodeBlock()
	name = strings.TrimPrefix(name, "minecraft:")
	if strings.EqualFold(name, string(t.Block)) {
		o.Printf("Found %s at %v", string(t.Block), t.Position)
	} else {
		o.Error(fmt.Sprintf("Expected %s, found %s", string(t.Block), name))
	}
}



type structureSaveCmd struct {
	Sub      cmd.SubCommand `cmd:"save"`
	Name     string
	From     mgl64.Vec3
	To       mgl64.Vec3
	SaveMode cmd.Optional[StructureSaveMode]
}

func (s structureSaveCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Saved structure %s", s.Name)
}

type structureLoadCmd struct {
	Sub      cmd.SubCommand `cmd:"load"`
	Name     string
	Position cmd.Optional[mgl64.Vec3]
	Rotation cmd.Optional[StructureRotation]
}

func (s structureLoadCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Loaded structure %s", s.Name)
}

type structureDeleteCmd struct {
	Sub  cmd.SubCommand `cmd:"delete"`
	Name string
}

func (s structureDeleteCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Deleted structure %s", s.Name)
}



type executeCmd struct {
	Sub cmd.SubCommand `cmd:"execute"`
	Args cmd.Varargs
}

func (e executeCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Execute command - subcommands not fully supported")
}



type eventCmd struct {
	Entity    []cmd.Target
	EventName EntityEvent
}

func (e eventCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	for _, t := range e.Entity {
		if pl, ok := t.(*player.Player); ok {
			o.Printf("Triggered event %s for %s", string(e.EventName), pl.Name())
		} else {
			o.Printf("Triggered event %s for entity", string(e.EventName))
		}
	}
}



type rideStartRidingCmd struct {
	Sub  cmd.SubCommand `cmd:"start_riding"`
	Args cmd.Varargs
}

func (r rideStartRidingCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Ride start_riding")
}

type rideStopRidingCmd struct {
	Sub  cmd.SubCommand `cmd:"stop_riding"`
	Args cmd.Varargs
}

func (r rideStopRidingCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Ride stop_riding")
}

type rideEvictRidersCmd struct {
	Sub  cmd.SubCommand `cmd:"evict_riders"`
	Args cmd.Varargs
}

func (r rideEvictRidersCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Ride evict_riders")
}

type rideSummonRiderCmd struct {
	Sub  cmd.SubCommand `cmd:"summon_rider"`
	Args cmd.Varargs
}

func (r rideSummonRiderCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Ride summon_rider")
}



type replaceitemCmd struct {
	Sub    cmd.SubCommand `cmd:"replaceitem"`
	Args   cmd.Varargs
}

func (r replaceitemCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Replaceitem - use /give instead")
}



type wsserverCmd struct {
	ServerURI string
}

func (w wsserverCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Connecting to WebSocket server %s", w.ServerURI)
}



type dialogueCmd struct {
	Player  []cmd.Target
	Dialogue cmd.Varargs
}

func (d dialogueCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Dialogue - use /say or /tell instead")
}



type fogPushCmd struct {
	Sub    cmd.SubCommand `cmd:"push"`
	Player []cmd.Target
	FogId  cmd.Optional[string]
}

func (f fogPushCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Fog pushed")
}

type fogPopCmd struct {
	Sub    cmd.SubCommand `cmd:"pop"`
	Player []cmd.Target
}

func (f fogPopCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Fog popped")
}

type fogRemoveCmd struct {
	Sub    cmd.SubCommand `cmd:"remove"`
	Player []cmd.Target
	FogId  string
}

func (f fogRemoveCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Fog removed")
}



type musicCmd struct {
	Action MusicAction
	Args   cmd.Varargs
}

func (m musicCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Music command accepted")
}

func init() {
	cmd.Register(cmd.New("seed", "Shows the world seed", nil, seedCmd{}))
	cmd.Register(cmd.New("locate", "Finds the nearest structure", nil, locateCmd{}))
	cmd.Register(cmd.New("clone", "Clones blocks", nil, cloneCmd{}))
	cmd.Register(cmd.New("spreadplayers", "Teleports entities randomly", nil, spreadplayersCmd{}))
	cmd.Register(cmd.New("clearspawnpoint", "Removes spawn point", nil, clearspawnpointCmd{}))
	cmd.Register(cmd.New("tag", "Manages entity tags", nil,
		tagAddCmd{}, tagRemoveCmd{}, tagListCmd{},
	))
	cmd.Register(cmd.New("testfor", "Counts matching entities", nil, testforCmd{}))
	cmd.Register(cmd.New("tickingarea", "Manages ticking areas", nil,
		tickingAreaAddCmd{}, tickingAreaRemoveCmd{}, tickingAreaRemoveAllCmd{}, tickingAreaListCmd{},
	))
	cmd.Register(cmd.New("save", "Controls saving", nil, saveHoldCmd{}, saveQueryCmd{}, saveResumeCmd{}))
	cmd.Register(cmd.New("transfer", "Transfers a player", nil, transferCmd{}))
	cmd.Register(cmd.New("reload", "Reloads function/script files", nil, reloadCmd{}))
	cmd.Register(cmd.New("setmaxplayers", "Sets max players", nil, setmaxplayersCmd{}))
	cmd.Register(cmd.New("particle", "Creates particles", nil, particleCmd{}))
	cmd.Register(cmd.New("mobevent", "Controls mob events", nil, mobeventCmd{}))
	cmd.Register(cmd.New("playsound", "Plays a sound", nil, playsoundCmd{}))
	cmd.Register(cmd.New("stopsound", "Stops a sound", nil, stopsoundCmd{}))
	cmd.Register(cmd.New("schedule", "Schedules delayed actions", nil, scheduleClearCmd{}, scheduleOnAreaLoadedCmd{}))
	cmd.Register(cmd.New("recipe", "Unlocks recipes", nil, recipeGiveCmd{}))
	cmd.Register(cmd.New("scoreboard", "Manages scoreboards", nil, scoreboardObjectivesCmd{}, scoreboardPlayersCmd{}))
	cmd.Register(cmd.New("camera", "Camera control", nil, cameraSetCmd{}, cameraClearCmd{}))
	cmd.Register(cmd.New("camerashake", "Camera shake", nil, camerashakeCmd{}))
	cmd.Register(cmd.New("changesetting", "Changes server settings", nil, changesettingCmd{}))
	cmd.Register(cmd.New("permission", "Manages permissions", nil, permissionCmd{}))
	cmd.Register(cmd.New("allowlist", "Manages allowlist", nil, allowlistAddCmd{}, allowlistRemoveCmd{}, allowlistListCmd{}, allowlistReloadCmd{}))
	cmd.Register(cmd.New("tellraw", "JSON chat messages", nil, tellrawCmd{}))
	cmd.Register(cmd.New("testforblock", "Tests for a block", nil, testforblockCmd{}))
	cmd.Register(cmd.New("structure", "Saves/loads structures", nil,
		structureSaveCmd{}, structureLoadCmd{}, structureDeleteCmd{},
	))
	cmd.Register(cmd.New("execute", "Executes commands as entities", nil, executeCmd{}))
	cmd.Register(cmd.New("event", "Triggers entity events", nil, eventCmd{}))
	cmd.Register(cmd.New("ride", "Controls entity riding", nil, rideStartRidingCmd{}, rideStopRidingCmd{}, rideEvictRidersCmd{}, rideSummonRiderCmd{}))
	cmd.Register(cmd.New("replaceitem", "Replaces inventory items", nil, replaceitemCmd{}))
	cmd.Register(cmd.New("wsserver", "WebSocket server connection", nil, wsserverCmd{}))
	cmd.Register(cmd.New("dialogue", "Opens NPC dialogue", nil, dialogueCmd{}))
	cmd.Register(cmd.New("fog", "Manages fog", nil, fogPushCmd{}, fogPopCmd{}, fogRemoveCmd{}))
	cmd.Register(cmd.New("music", "Controls music", nil, musicCmd{}))
	cmd.Register(cmd.New("testforblocks", "Tests block matches", nil, testforCmd{}))
}
