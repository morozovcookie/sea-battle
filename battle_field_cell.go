package seabattle

type BattleFieldCell struct {
	bs *BattleShip

	destroyed bool

	originalRowName    string
	originalColumnName string

	rowIndex    int
	columnIndex int
}

func (cell BattleFieldCell) BattleShip() *BattleShip {
	return cell.bs
}

func (cell BattleFieldCell) IsDestroyed() bool {
	return cell.destroyed
}

func (cell *BattleFieldCell) Destroy() {
	cell.destroyed = true
}

func (cell BattleFieldCell) OriginalRowName() string {
	return cell.originalRowName
}

func (cell BattleFieldCell) OriginalColumnName() string {
	return cell.originalColumnName
}

func (cell BattleFieldCell) RowIndex() int {
	return cell.rowIndex
}

func (cell BattleFieldCell) ColumnIndex() int {
	return cell.columnIndex
}

func (cell BattleFieldCell) IsEmpty() bool {
	return cell.bs == nil
}

func (cell *BattleFieldCell) Clear() {
	cell.bs = nil
	cell.destroyed = false
}

func (cell *BattleFieldCell) AttachShip(ship *BattleShip) {
	cell.bs = ship
}
