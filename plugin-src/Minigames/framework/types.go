package framework

import (
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/go-gl/mathgl/mgl64"
)


type MatchState int

const (
	MatchStateWaiting   MatchState = iota 
	MatchStateCountdown                   
	MatchStateStarting                    
	MatchStatePlaying                     
	MatchStateFinished                    
	MatchStateClosed                      
)

func (s MatchState) String() string {
	switch s {
	case MatchStateWaiting:
		return "waiting"
	case MatchStateCountdown:
		return "countdown"
	case MatchStateStarting:
		return "starting"
	case MatchStatePlaying:
		return "playing"
	case MatchStateFinished:
		return "finished"
	case MatchStateClosed:
		return "closed"
	default:
		return "unknown"
	}
}


type Arena struct {
	Name       string
	WorldName  string
	Spawns     []mgl64.Vec3
	Chests     []ChestPosition
	MaxPlayers int
}


type ChestPosition struct {
	Pos   mgl64.Vec3
	Type  ChestType
}


type ChestType int

const (
	ChestTypeNormal ChestType = iota
	ChestTypeCenter
)


type EndReason string

const (
	EndReasonLastStanding EndReason = "last_standing"
	EndReasonTimeLimit    EndReason = "time_limit"
	EndReasonDraw         EndReason = "draw"
	EndReasonForce        EndReason = "force"
)


type RemoveReason string

const (
	RemoveReasonQuit     RemoveReason = "quit"
	RemoveReasonDeath    RemoveReason = "death"
	RemoveReasonKicked   RemoveReason = "kicked"
	RemoveReasonTransfer RemoveReason = "transfer"
)


type WorldHandler interface {
	HandleChestLoot(tx *world.Tx, pos mgl64.Vec3, chestType ChestType)
}
