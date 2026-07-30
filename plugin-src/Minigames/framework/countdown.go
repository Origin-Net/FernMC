package framework

import (
	"sync/atomic"
	"time"

	"github.com/Origin-Net/FernMC/server/player"
	"github.com/Origin-Net/FernMC/server/world"
)


type CountdownAction struct {
	AtSecond   int
	OnTick     func(remaining int)
}


type Countdown struct {
	seconds     int
	ticker      *time.Ticker
	running     atomic.Bool
	remaining   atomic.Int32
	actions     []CountdownAction
	onComplete  func()
	onCancel    func()
	broadcaster func(msg string)
	soundFunc   func(sound world.Sound)
	stopCh      chan struct{}
	logger      *LogWrapper
}


func NewCountdown(seconds int, logger *LogWrapper) *Countdown {
	return &Countdown{
		seconds: seconds,
		stopCh:  make(chan struct{}),
		logger:  logger,
	}
}


func (c *Countdown) WithActions(actions ...CountdownAction) *Countdown {
	c.actions = actions
	return c
}


func (c *Countdown) OnComplete(fn func()) *Countdown {
	c.onComplete = fn
	return c
}


func (c *Countdown) OnCancel(fn func()) *Countdown {
	c.onCancel = fn
	return c
}


func (c *Countdown) OnBroadcast(fn func(msg string)) *Countdown {
	c.broadcaster = fn
	return c
}


func (c *Countdown) Start() {
	if !c.running.CompareAndSwap(false, true) {
		return
	}
	c.remaining.Store(int32(c.seconds))

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				rem := int(c.remaining.Load())
				if rem <= 0 {
					c.stop()
					if c.onComplete != nil {
						c.onComplete()
					}
					return
				}

				rem--
				c.remaining.Store(int32(rem))

				for _, action := range c.actions {
					if action.AtSecond == rem || (action.AtSecond < 0 && rem <= -action.AtSecond) {
						action.OnTick(rem)
					}
				}

			case <-c.stopCh:
				if c.onCancel != nil {
					c.onCancel()
				}
				return
			}
		}
	}()
}


func (c *Countdown) Cancel() {
	if c.running.Load() {
		c.stop()
	}
}

func (c *Countdown) stop() {
	c.running.Store(false)
	select {
	case c.stopCh <- struct{}{}:
	default:
	}
}


func (c *Countdown) IsRunning() bool {
	return c.running.Load()
}


func (c *Countdown) Remaining() int {
	return int(c.remaining.Load())
}


func SetPlayerScale(players []*player.Player, remaining int, startMessage string) {
	for _, p := range players {
		if remaining <= 5 && remaining > 0 {
			p.SendTip(startMessage)
		}
	}
}
