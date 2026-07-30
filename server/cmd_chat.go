package server

import (
	"fmt"
	"time"

	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/player/title"
	"github.com/Origin-Net/FernMC/server/world"
)

type meCmd struct {
	Action cmd.Varargs
}

func (m meCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("No world context")
		return
	}
	srcName := srcName(src)
	msg := fmt.Sprintf("* %s %s", srcName, string(m.Action))
	for p := range tx.Players() {
		p.(*player.Player).Message(msg)
	}
	o.Print(msg)
}

type titleClear struct {
	Sub    cmd.SubCommand `cmd:"clear"`
	Player []cmd.Target
}

func (t titleClear) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	for _, target := range t.Player {
		p, ok := target.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		p.SendTitle(title.New(""))
		o.Printf("Cleared title for %s", p.Name())
	}
}

type titleReset struct {
	Sub    cmd.SubCommand `cmd:"reset"`
	Player []cmd.Target
}

func (t titleReset) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	for _, target := range t.Player {
		p, ok := target.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		p.SendTitle(title.New(""))
		o.Printf("Reset title for %s", p.Name())
	}
}

type titleTitle struct {
	Sub   cmd.SubCommand `cmd:"title"`
	Player []cmd.Target
	Text  cmd.Varargs
}

func (t titleTitle) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	for _, target := range t.Player {
		p, ok := target.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		p.SendTitle(title.New(string(t.Text)))
		o.Printf("Sent title to %s", p.Name())
	}
}

type titleSubtitle struct {
	Sub      cmd.SubCommand `cmd:"subtitle"`
	Player   []cmd.Target
	Text cmd.Varargs
}

func (t titleSubtitle) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	for _, target := range t.Player {
		p, ok := target.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		t := title.New("").WithSubtitle(string(t.Text))
		p.SendTitle(t)
		o.Printf("Sent subtitle to %s", p.Name())
	}
}

type titleActionbar struct {
	Sub      cmd.SubCommand `cmd:"actionbar"`
	Player   []cmd.Target
	Text cmd.Varargs
}

func (t titleActionbar) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	for _, target := range t.Player {
		p, ok := target.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		t := title.New("").WithActionText(string(t.Text))
		p.SendTitle(t)
		o.Printf("Sent action bar to %s", p.Name())
	}
}

type titleTimes struct {
	Sub     cmd.SubCommand `cmd:"times"`
	Player  []cmd.Target
	FadeIn  int
	Stay    int
	FadeOut int
}

func (t titleTimes) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	for _, target := range t.Player {
		p, ok := target.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		t := title.New("").
			WithFadeInDuration(time.Duration(t.FadeIn)*time.Second/20).
			WithDuration(time.Duration(t.Stay)*time.Second/20).
			WithFadeOutDuration(time.Duration(t.FadeOut)*time.Second/20)
		p.SendTitle(t)
		o.Printf("Set title times for %s", p.Name())
	}
}

func init() {
	cmd.Register(cmd.New("me", "Displays a message about yourself", nil, meCmd{}))
	cmd.Register(cmd.New("title", "Controls screen titles", nil,
		titleClear{}, titleReset{}, titleTitle{}, titleSubtitle{}, titleActionbar{}, titleTimes{},
	))
}
