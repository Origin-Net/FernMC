package entity

import (
	"github.com/Origin-Net/FernMC/server/world"
	"time"
)


type SwingArmAction struct{ action }



type HurtAction struct{ action }



type CriticalHitAction struct {
	action
	
	Count int
}



type EnchantedHitAction struct {
	action
	
	Count int
}



type DeathAction struct{ action }



type EatAction struct{ action }


type ArrowShakeAction struct {
	
	Duration time.Duration

	action
}



type PickedUpAction struct {
	
	Collector world.Entity

	action
}


type FireworkExplosionAction struct{ action }


type TotemUseAction struct{ action }



type action struct{}

func (action) EntityAction() {}
