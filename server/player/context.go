package player

import (
	"github.com/Origin-Net/FernMC/server/world"
)




type Context struct {
	*world.Context
	p *Player
}


func newContext(p *Player) *Context {
	return &Context{Context: p.tx.Event(), p: p}
}



func (ctx *Context) Player() *Player { return ctx.p }





func (ctx *Context) Defer(f func(ctx *Context)) *world.Task {
	return ctx.DeferErr(func(ctx *Context) error {
		f(ctx)
		return nil
	})
}


func (ctx *Context) DeferErr(f func(ctx *Context) error) *world.Task {
	h := ctx.p.H()
	return ctx.Context.DeferErr(func(tx *world.Tx) error {
		if e, ok := h.Entity(tx); ok {
			return f(newContext(e.(*Player)))
		}
		if h.Closed() {
			return world.ErrEntityClosed
		}
		return world.ErrEntityNotInWorld
	})
}
