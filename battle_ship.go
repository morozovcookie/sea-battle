package seabattle

type BattleShipCellsIndexKey struct {
	row string
	col string
}

type BattleShip struct {
	leftTop     *BattleFieldCell
	rightBottom *BattleFieldCell

	cells map[BattleShipCellsIndexKey]*BattleFieldCell

	destroyed bool
}

func NewBattleShip(leftTop, rightBottom *BattleFieldCell) *BattleShip {
	return &BattleShip{
		leftTop:     leftTop,
		rightBottom: rightBottom,
		cells:       make(map[BattleShipCellsIndexKey]*BattleFieldCell),
	}
}

func (bs BattleShip) IsDestroyed() (destroy bool) {
	if bs.destroyed {
		return true
	}

	for _, cell := range bs.cells {
		if !cell.IsDestroyed() {
			return false
		}
	}

	return true
}

func (bs *BattleShip) Destroy() {
	if bs.destroyed {
		return
	}

	bs.destroyed = true
}

func (bs BattleShip) LeftTopCell() *BattleFieldCell {
	return bs.leftTop
}

func (bs BattleShip) RightBottomCell() *BattleFieldCell {
	return bs.rightBottom
}

func (bs *BattleShip) AttachCells(cells ...*BattleFieldCell) {
	for _, cell := range cells {
		cell.AttachShip(bs)
		bs.cells[BattleShipCellsIndexKey{row: cell.OriginalRowName(), col: cell.OriginalColumnName()}] = cell
	}
}
