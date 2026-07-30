package block

import (
	"github.com/Origin-Net/FernMC/server/block/cube"
	"github.com/Origin-Net/FernMC/server/item"
	"github.com/Origin-Net/FernMC/server/item/inventory"
	"github.com/Origin-Net/FernMC/server/world"
)


type ContainerViewer interface {
	world.Viewer
	
	
	ViewSlotChange(slot int, newItem item.Stack)
}


type ContainerOpener interface {
	
	OpenBlockContainer(pos cube.Pos, tx *world.Tx)
}



type Container interface {
	AddViewer(v ContainerViewer, tx *world.Tx, pos cube.Pos)
	RemoveViewer(v ContainerViewer, tx *world.Tx, pos cube.Pos)
	Inventory(tx *world.Tx, pos cube.Pos) *inventory.Inventory
}
