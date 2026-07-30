package server

import (
	"fmt"
	"strings"

	"github.com/Origin-Net/FernMC/server/cmd"
	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/world"
)

type stopCmd struct{}

func (stopCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	srv := GetServer()
	if srv == nil {
		o.Error("Server not available")
		return
	}
	o.Print("Stopping the server")
	go srv.Close()
}

type listCmd struct{}

func (listCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("Cannot list players from console")
		return
	}
	var names []string
	for e := range tx.Players() {
		p := e.(*player.Player)
		names = append(names, p.Name())
	}
	o.Printf("There are %d/%d players online:\n%s", len(names), -1, strings.Join(names, ", "))
}

type kickCmd struct {
	Player []cmd.Target
	Reason cmd.Optional[cmd.Varargs]
}

func (k kickCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	for _, t := range k.Player {
		p, ok := t.(*player.Player)
		if !ok {
			o.Error("That player cannot be found")
			return
		}
		reason := ""
		if r, ok := k.Reason.Load(); ok {
			reason = string(r)
		}
		p.Disconnect(reason)
		if reason != "" {
			o.Printf("Kicked %s from the game: '%s'", p.Name(), reason)
		} else {
			o.Printf("Kicked %s from the game", p.Name())
		}
	}
}

type opCmd struct {
	Player []cmd.Target
}

func (o opCmd) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	if tx == nil {
		out.Error("No world context")
		return
	}
	for _, t := range o.Player {
		p, ok := t.(*player.Player)
		if !ok {
			out.Error("Target must be a player")
			return
		}
		out.Printf("Opped %s", p.Name())
	}
}

type deopCmd struct {
	Player []cmd.Target
}

func (d deopCmd) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	if tx == nil {
		out.Error("No world context")
		return
	}
	for _, t := range d.Player {
		p, ok := t.(*player.Player)
		if !ok {
			out.Error("Target must be a player")
			return
		}
		out.Printf("De-opped %s", p.Name())
	}
}

type sayCmd struct {
	Message cmd.Varargs
}

func (s sayCmd) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if tx == nil {
		o.Error("Cannot use say from console")
		return
	}
	msg := fmt.Sprintf("[%s] %s", srcName(src), string(s.Message))
	for e := range tx.Players() {
		e.(*player.Player).Message(msg)
	}
	o.Print(msg)
}

type tellCmd struct {
	Target  []cmd.Target
	Message cmd.Varargs
}

func (t tellCmd) Run(src cmd.Source, o *cmd.Output, _ *world.Tx) {
	for _, target := range t.Target {
		p, ok := target.(*player.Player)
		if !ok {
			o.Error("Target must be a player")
			return
		}
		p.Message(fmt.Sprintf("%s whispers to you: %s", srcName(src), string(t.Message)))
	}
	if len(t.Target) > 0 {
		o.Printf("You whisper to %s: %s", targetName(t.Target[0]), string(t.Message))
	}
}

func init() {
	cmd.Register(cmd.New("stop", "Stops the server", []string{"shutdown"}, stopCmd{}))
	cmd.Register(cmd.New("list", "Lists players on the server", nil, listCmd{}))
	cmd.Register(cmd.New("kick", "Kicks a player from the server", nil, kickCmd{}))
	cmd.Register(cmd.New("op", "Grants operator status to a player", nil, opCmd{}))
	cmd.Register(cmd.New("deop", "Revokes operator status from a player", nil, deopCmd{}))
	cmd.Register(cmd.New("say", "Sends a message in the chat to other players", nil, sayCmd{}))
	cmd.Register(cmd.New("tell", "Sends a private message to one or more players", []string{"msg", "w"}, tellCmd{}))
}
