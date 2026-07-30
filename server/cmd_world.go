package server

import (
	"strings"
	"time"

	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/world"
)

type DifficultyEnum string

func (DifficultyEnum) Type() string { return "Difficulty" }
func (DifficultyEnum) Options(cmd.Source) []string {
	return []string{"peaceful", "easy", "normal", "hard"}
}

type timeQuery struct {
	Sub cmd.SubCommand `cmd:"query"`
}

func (timeQuery) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	o.Printf("Time is %d", tx.World().Time())
}

type timeStart struct {
	Sub cmd.SubCommand `cmd:"start"`
}

func (timeStart) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	tx.World().StartTime()
	o.Print("Restarted the time")
}

type timeStop struct {
	Sub cmd.SubCommand `cmd:"stop"`
}

func (timeStop) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	tx.World().StopTime()
	o.Print("Stopped the time")
}

type TimeOfDayEnum string

func (TimeOfDayEnum) Type() string { return "TimeOfDay" }
func (TimeOfDayEnum) Options(cmd.Source) []string {
	return []string{"day", "noon", "night", "midnight"}
}

type timeSetNamed struct {
	Sub  cmd.SubCommand `cmd:"set"`
	Time TimeOfDayEnum
}

func (t timeSetNamed) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	w := tx.World()
	switch strings.ToLower(string(t.Time)) {
	case "day":
		w.SetTime(1000)
		o.Print("Set the time to day")
	case "noon":
		w.SetTime(6000)
		o.Print("Set the time to noon")
	case "night":
		w.SetTime(13000)
		o.Print("Set the time to night")
	case "midnight":
		w.SetTime(18000)
		o.Print("Set the time to midnight")
	}
}

type timeSetInt struct {
	Sub  cmd.SubCommand `cmd:"set"`
	Time int
}

func (t timeSetInt) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	tx.World().SetTime(t.Time)
	o.Printf("Set the time to %d", t.Time)
}

type timeAdd struct {
	Sub   cmd.SubCommand `cmd:"add"`
	Value int
}

func (t timeAdd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	w := tx.World()
	w.SetTime(w.Time() + t.Value)
	o.Printf("Added %d to the time", t.Value)
}

type weatherClear struct {
	Sub      cmd.SubCommand  `cmd:"clear"`
	Duration cmd.Optional[int]
}

func (w weatherClear) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	tx.World().StopRaining()
	tx.World().StopThundering()
	o.Print("Set weather to clear")
}

type weatherRain struct {
	Sub      cmd.SubCommand  `cmd:"rain"`
	Duration cmd.Optional[int]
}

func (w weatherRain) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	dur := 1 * time.Hour
	if d, ok := w.Duration.Load(); ok && d > 0 {
		dur = time.Duration(d) * time.Second
	}
	tx.World().StartRaining(dur)
	o.Print("Set weather to rain")
}

type weatherThunder struct {
	Sub      cmd.SubCommand  `cmd:"thunder"`
	Duration cmd.Optional[int]
}

func (w weatherThunder) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	dur := 1 * time.Hour
	if d, ok := w.Duration.Load(); ok && d > 0 {
		dur = time.Duration(d) * time.Second
	}
	tx.World().StartRaining(dur)
	tx.World().StartThundering(dur)
	o.Print("Set weather to thunder")
}

type difficultyCmd struct {
	Difficulty DifficultyEnum
}

func (d difficultyCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	var diff world.Difficulty
	switch strings.ToLower(string(d.Difficulty)) {
	case "peaceful":
		diff = world.DifficultyPeaceful
	case "easy":
		diff = world.DifficultyEasy
	case "normal":
		diff = world.DifficultyNormal
	case "hard":
		diff = world.DifficultyHard
	}
	tx.World().SetDifficulty(diff)
	o.Printf("Set game difficulty to %s", string(d.Difficulty))
}

type gameruleCmd struct {
	Rule  GameRule
	Value cmd.Optional[Bool]
}

func (g gameruleCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	w := tx.World()
	rule := strings.ToLower(string(g.Rule))

	if v, ok := g.Value.Load(); ok {
		val := strings.ToLower(string(v))
		switch rule {
		case "dodaylightcycle":
			if val == "true" {
				w.StartTime()
			} else {
				w.StopTime()
			}
			o.Printf("Set game rule %s to %s", string(g.Rule), string(v))
		case "doweathercycle":
			if val == "true" {
				w.StartWeatherCycle()
			} else {
				w.StopWeatherCycle()
			}
			o.Printf("Set game rule %s to %s", string(g.Rule), string(v))
		case "showcoordinates":
			for e := range tx.Players() {
				if val == "true" {
					e.(*player.Player).ShowCoordinates()
				} else {
					e.(*player.Player).HideCoordinates()
				}
			}
			o.Printf("Set game rule %s to %s", string(g.Rule), string(v))
		case "doimmediaterespawn":
			for e := range tx.Players() {
				if val == "true" {
					e.(*player.Player).EnableInstantRespawn()
				} else {
					e.(*player.Player).DisableInstantRespawn()
				}
			}
			o.Printf("Set game rule %s to %s", string(g.Rule), string(v))
		case "keepinventory", "mobgriefing", "domobspawning", "dofiretick",
			"doentitydrops", "dotiledrops", "falldamage", "firedamage",
			"drowningdamage", "freezedamage", "naturalregeneration",
			"pvp", "tntexplodes", "commandblockoutput", "commandblocksenabled",
			"sendcommandfeedback", "randomtickspeed", "maxcommandchainlength",
			"playerssleepingpercentage", "functioncommandlimit", "spawnradius",
			"showbordereffect", "showtags", "showdeathmessages",
			"showdaysplayed", "showrecipemessages", "respawnblocksexplode",
			"recipesunlock", "projectilescanbreakblocks",
			"doinsomnia", "dolimitedcrafting", "tntexplosiondropdecay":
			for e := range tx.Players() {
				p, ok := e.(*player.Player)
				if ok {
					boolVal := val == "true"
					p.SendGameRule(rule, boolVal)
				}
			}
			o.Printf("Set game rule %s to %s", string(g.Rule), string(v))
		default:
			o.Error("Unknown game rule: " + string(g.Rule))
		}
		return
	}
	o.Printf("Game rule %s is currently not queryable from here", string(g.Rule))
}

type toggleDownfallCmd struct{}

func (toggleDownfallCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	w := tx.World()
	if w.Raining() {
		w.StopRaining()
		o.Print("Toggled downfall off")
	} else {
		w.StartRaining(30 * time.Second)
		o.Print("Toggled downfall on")
	}
}

type daylockCmd struct {
	Lock cmd.Optional[Bool]
}

func (d daylockCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	w := tx.World()
	lock := true
	if l, ok := d.Lock.Load(); ok {
		lock = strings.ToLower(string(l)) == "true"
	}
	if lock {
		w.StopTime()
		o.Print("Day cycle locked")
	} else {
		w.StartTime()
		o.Print("Day cycle unlocked")
	}
}

func init() {
	cmd.Register(cmd.New("time", "Changes or queries the world's game time", nil,
		timeQuery{}, timeStart{}, timeStop{}, timeSetNamed{}, timeSetInt{}, timeAdd{},
	))
	cmd.Register(cmd.New("difficulty", "Sets the difficulty level", nil, difficultyCmd{}))
	cmd.Register(cmd.New("gamerule", "Sets or queries a game rule value", nil, gameruleCmd{}))
	cmd.Register(cmd.New("toggledownfall", "Toggles the weather", nil, toggleDownfallCmd{}))
	cmd.Register(cmd.New("weather", "Sets the weather in the environment", nil,
		weatherClear{}, weatherRain{}, weatherThunder{},
	))
	cmd.Register(cmd.New("daylock", "Locks and unlocks the day-night cycle", nil, daylockCmd{}))
}
