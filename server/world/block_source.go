package world

import "github.com/Origin-Net/FernMC/server/block/cube"


type BlockSource interface {
	
	Block(cube.Pos) Block
}


type worldSource struct{ tx *Tx }

func (w worldSource) Block(pos cube.Pos) Block { return w.tx.block(pos) }
