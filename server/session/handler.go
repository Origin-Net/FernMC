package session

import (
	"github.com/Origin-Net/FernMC/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)


type packetHandler interface {
	
	
	Handle(p packet.Packet, s *Session, tx *world.Tx, c Controllable) error
}
