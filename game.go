package seabattle

type Game struct {
	bf *BattleField

	gs *GameState
}

func NewGame(creator BattleFieldCreator) (g *Game, err error) {
	g = &Game{
		gs: &GameState{},
	}

	if g.bf, err = creator(); err != nil {
		return nil, err
	}

	g.bf.ShotEvent = func() {
		g.gs.IncShotCount()
	}
	g.bf.DestroyEvent = func() {
		g.gs.IncDestroyed()
	}
	g.bf.KnockEvent = func() {
		g.gs.IncKnocked()
	}
	g.bf.PlaceShipEvent = func() {
		g.gs.IncShipCount()
	}
	g.bf.ClearFieldEvent = func() {
		g.gs.Clear()
	}

	return g, nil
}

func (g *Game) BattleField() *BattleField {
	return g.bf
}

func (g *Game) IsOver() bool {
	if g.gs.ShipCount() == 0 {
		return false
	}

	return g.gs.ShipCount() == g.gs.Destroyed()
}

func (g *Game) State() *GameState {
	return g.gs
}
