package seabattle

import (
	"sync/atomic"
)

type GameState struct {
	shipCount int32
	destroyed int32
	knocked   int32
	shotCount int32
}

func (gs *GameState) IncShipCount() {
	atomic.AddInt32(&gs.shipCount, 1)
}

func (gs GameState) ShipCount() int32 {
	return atomic.LoadInt32(&gs.shipCount)
}

func (gs *GameState) IncDestroyed() {
	atomic.AddInt32(&gs.destroyed, 1)
}

func (gs *GameState) Destroyed() int32 {
	return atomic.LoadInt32(&gs.destroyed)
}

func (gs *GameState) IncKnocked() {
	atomic.AddInt32(&gs.knocked, 1)
}

func (gs *GameState) Knocked() int32 {
	return atomic.LoadInt32(&gs.knocked)
}

func (gs *GameState) IncShotCount() {
	atomic.AddInt32(&gs.shotCount, 1)
}

func (gs *GameState) ShotCount() int32 {
	return atomic.LoadInt32(&gs.shotCount)
}

func (gs *GameState) Clear() {
	atomic.StoreInt32(&gs.shipCount, 0)
	atomic.StoreInt32(&gs.destroyed, 0)
	atomic.StoreInt32(&gs.knocked, 0)
	atomic.StoreInt32(&gs.shotCount, 0)
}
