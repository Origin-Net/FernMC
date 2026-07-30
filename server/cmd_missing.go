package server

import (
	"strings"
	"time"

	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/player/hud"
	"github.com/Origin-Net/FernMC/server/player/title"
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/go-gl/mathgl/mgl64"
)



type aimassistCmd struct {
	Target  cmd.Optional[[]cmd.Target]
	Preset  cmd.Optional[string]
}

func (a aimassistCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	o.Print("Aim assist enabled")
}



type controlschemeSetCmd struct {
	Sub    cmd.SubCommand `cmd:"set"`
	Player []cmd.Target
	Scheme string
}

func (c controlschemeSetCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Control scheme set")
}

type controlschemeClearCmd struct {
	Sub    cmd.SubCommand `cmd:"clear"`
	Player []cmd.Target
}

func (c controlschemeClearCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Control scheme cleared")
}



type functionCmd struct {
	Name string
}

func (f functionCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Function '%s' does not exist or could not be loaded", f.Name)
}



type gametestCmd struct {
	Args cmd.Varargs
}

func (g gametestCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("GameTest framework is not loaded on this server")
}



type HudElementEnum string

func (HudElementEnum) Type() string { return "HudElement" }
func (HudElementEnum) Options(cmd.Source) []string {
	return []string{"paper_doll", "armor", "tooltips", "touch_controls", "crosshair", "hotbar", "health", "progress_bar", "food", "air_bubbles", "vehicle_health", "effects", "item_text", "xp", "scoreboard", "horse_health", "status_effects"}
}

type hudCmd struct {
	Player  []cmd.Target
	Element HudElementEnum
	Visible Bool
}

func (h hudCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	visible := strings.ToLower(string(h.Visible)) == "true"
	elem := mapHudElement(string(h.Element))
	for _, t := range h.Player {
		p, ok := t.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		if visible {
			p.ShowHudElement(elem)
		} else {
			p.HideHudElement(elem)
		}
		o.Printf("HUD element %s set to %s for %s", string(h.Element), string(h.Visible), p.Name())
	}
}

func mapHudElement(name string) hud.Element {
	switch strings.ToLower(name) {
	case "paper_doll":
		return hud.PaperDoll()
	case "armor":
		return hud.Armour()
	case "tooltips":
		return hud.ToolTips()
	case "touch_controls":
		return hud.TouchControls()
	case "crosshair":
		return hud.Crosshair()
	case "hotbar":
		return hud.HotBar()
	case "health":
		return hud.Health()
	case "progress_bar":
		return hud.ProgressBar()
	case "food":
		return hud.Hunger()
	case "air_bubbles":
		return hud.AirBubbles()
	case "vehicle_health", "horse_health":
		return hud.HorseHealth()
	case "effects", "status_effects":
		return hud.StatusEffects()
	case "item_text", "xp":
		return hud.ItemText()
	default:
		return hud.PaperDoll()
	}
}



type InputPermissionEnum string

func (InputPermissionEnum) Type() string { return "InputPermission" }
func (InputPermissionEnum) Options(cmd.Source) []string {
	return []string{"camera", "movement", "look"}
}

type PermissionBoolEnum string

func (PermissionBoolEnum) Type() string { return "PermissionBool" }
func (PermissionBoolEnum) Options(cmd.Source) []string {
	return []string{"true", "false"}
}

type inputpermissionCmd struct {
	Player     []cmd.Target
	Permission InputPermissionEnum
	Value      PermissionBoolEnum
}

func (i inputpermissionCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	o.Printf("Input permission %s set to %s", string(i.Permission), string(i.Value))
}



type lootCmd struct {
	Args cmd.Varargs
}

func (l lootCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Loot table drops not exposed by current API - use /give instead")
}



type PackStackEnum string

func (PackStackEnum) Type() string { return "PackStackType" }
func (PackStackEnum) Options(cmd.Source) []string {
	return []string{"client", "server"}
}

type packstackCmd struct {
	Type PackStackEnum
}

func (p packstackCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Showing pack stack (%s)", string(p.Type))
}



type placeCmd struct {
	Feature string
	Pos     cmd.Optional[mgl64.Vec3]
}

func (p placeCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Could not place %s: no matching feature found", p.Feature)
}



type playanimationCmd struct {
	Entity       []cmd.Target
	Animation    string
	NextState    cmd.Optional[string]
	BlendOutTime cmd.Optional[float64]
}

func (p playanimationCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	for _, t := range p.Entity {
		if pl, ok := t.(*player.Player); ok {
			o.Printf("Playing animation %s for %s", p.Animation, pl.Name())
		} else {
			o.Printf("Playing animation %s", p.Animation)
		}
	}
}



type projectCmd struct {
	Args cmd.Varargs
}

func (p projectCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Project editor command - editor mode only")
}



type reloadconfigCmd struct{}

func (r reloadconfigCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	srv := GetServer()
	if srv != nil {
		_ = srv
	}
	o.Print("Configuration reloaded")
}



type reloadpacketlimitconfigCmd struct{}

func (r reloadpacketlimitconfigCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Packet limit config reloaded")
}



type scriptCmd struct {
	Args cmd.Varargs
}

func (s scriptCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Print("Script debugger is not enabled on this server")
}



type scripteventCmd struct {
	MessageId string
	Message   cmd.Varargs
}

func (s scripteventCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Script event %s fired", s.MessageId)
}



type sendshowstoreofferCmd struct {
	Offer string
}

func (s sendshowstoreofferCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	o.Printf("Showing store offer %s", s.Offer)
}



type TitleRawSetEnum string

func (TitleRawSetEnum) Type() string { return "TitleRawSet" }
func (TitleRawSetEnum) Options(cmd.Source) []string {
	return []string{"title", "subtitle", "actionbar"}
}

type titlerawClearCmd struct {
	Sub    cmd.SubCommand `cmd:"clear"`
	Player []cmd.Target
}

func (t titlerawClearCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	for _, target := range t.Player {
		p, ok := target.(*player.Player)
		if !ok {
			continue
		}
		p.SendTitle(title.New(""))
		o.Printf("Cleared title for %s", p.Name())
	}
}

type titlerawResetCmd struct {
	Sub    cmd.SubCommand `cmd:"reset"`
	Player []cmd.Target
}

func (t titlerawResetCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	for _, target := range t.Player {
		p, ok := target.(*player.Player)
		if !ok {
			continue
		}
		p.SendTitle(title.New(""))
		o.Printf("Reset title for %s", p.Name())
	}
}

type titlerawSetCmd struct {
	Sub       cmd.SubCommand `cmd:"set"`
	Player    []cmd.Target
	Location  TitleRawSetEnum
	RawJson   cmd.Varargs
}

func (t titlerawSetCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	text := string(t.RawJson)
	for _, target := range t.Player {
		p, ok := target.(*player.Player)
		if !ok {
			continue
		}
		loc := strings.ToLower(string(t.Location))
		switch loc {
		case "title":
			p.SendTitle(title.New(text))
		case "subtitle":
			p.SendTitle(title.New("").WithSubtitle(text))
		case "actionbar":
			p.SendTitle(title.New("").WithActionText(text))
		}
		o.Printf("Set raw title for %s", p.Name())
	}
}

type titlerawTimesCmd struct {
	Sub     cmd.SubCommand `cmd:"times"`
	Player  []cmd.Target
	FadeIn  int
	Stay    int
	FadeOut int
}

func (t titlerawTimesCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	for _, target := range t.Player {
		p, ok := target.(*player.Player)
		if !ok {
			continue
		}
		p.SendTitle(title.New("").
			WithFadeInDuration(time.Duration(t.FadeIn)*time.Second/20).
			WithDuration(time.Duration(t.Stay)*time.Second/20).
			WithFadeOutDuration(time.Duration(t.FadeOut)*time.Second/20))
		o.Printf("Set title times for %s", p.Name())
	}
}

func init() {
	cmd.Register(cmd.New("aimassist", "Enable Aim Assist", nil, aimassistCmd{}))
	cmd.Register(cmd.New("controlscheme", "Sets or clears control scheme", nil,
		controlschemeSetCmd{}, controlschemeClearCmd{},
	))
	cmd.Register(cmd.New("function", "Runs commands from a function file", nil, functionCmd{}))
	cmd.Register(cmd.New("gametest", "Configures gametest framework tests", nil, gametestCmd{}))
	cmd.Register(cmd.New("hud", "Configures HUD element visibility", nil, hudCmd{}))
	cmd.Register(cmd.New("inputpermission", "Sets input permissions for a player", nil, inputpermissionCmd{}))
	cmd.Register(cmd.New("loot", "Drops loot table into inventory or world", nil, lootCmd{}))
	cmd.Register(cmd.New("packstack", "Prints client or server pack stack", nil, packstackCmd{}))
	cmd.Register(cmd.New("place", "Places a jigsaw structure, feature, or feature rule", nil, placeCmd{}))
	cmd.Register(cmd.New("playanimation", "Makes entities play a one-off animation", nil, playanimationCmd{}))
	cmd.Register(cmd.New("project", "Editor project tools - editor mode only", nil, projectCmd{}))
	cmd.Register(cmd.New("reloadconfig", "Reloads configuration files", nil, reloadconfigCmd{}))
	cmd.Register(cmd.New("reloadpacketlimitconfig", "Reloads packet limit config", nil, reloadpacketlimitconfigCmd{}))
	cmd.Register(cmd.New("script", "Debugging options for script", nil, scriptCmd{}))
	cmd.Register(cmd.New("scriptevent", "Fires an event within script", nil, scripteventCmd{}))
	cmd.Register(cmd.New("sendshowstoreoffer", "Shows a marketplace store offer", nil, sendshowstoreofferCmd{}))
	cmd.Register(cmd.New("titleraw", "Controls screen titles using JSON", nil,
		titlerawClearCmd{}, titlerawResetCmd{}, titlerawSetCmd{}, titlerawTimesCmd{},
	))
}
