package seabattle

import (
	"strconv"
)

const (
	MinBattleFieldSize = 1
	MaxBattleFieldSize = 10

	ShipPlacingMargin = 1
)

type BattleFieldCreator func() (*BattleField, error)

type Event func()

type BattleField struct {
	size int

	cells [][]BattleFieldCell

	ShotEvent       Event
	DestroyEvent    Event
	KnockEvent      Event
	PlaceShipEvent  Event
	ClearFieldEvent Event
}

func NewBattleField(size int) (bf *BattleField, err error) {
	if size < MinBattleFieldSize || size > MaxBattleFieldSize {
		return nil, ErrInvalidBattleFieldSize
	}

	bf = &BattleField{
		size:  size,
		cells: make([][]BattleFieldCell, size),
	}

	for i := 0; i < size; i++ {
		bf.cells[i] = make([]BattleFieldCell, size)

		for j := 0; j < size; j++ {
			bf.cells[i][j] = BattleFieldCell{
				originalRowName:    string(byte(i + 1 + '0')),
				originalColumnName: string(byte('A' + j)),
				rowIndex:           i,
				columnIndex:        j,
			}
		}
	}

	return bf, nil
}

func (bf BattleField) Size() int {
	return bf.size
}

func (bf *BattleField) PlaceBattleShips(ships ...*BattleShip) (err error) {
	defer func(err *error) {
		if *err == nil {
			return
		}

		bf.Clear()
	}(&err)

	for _, ship := range ships {
		var (
			leftTop     = ship.LeftTopCell()
			rightBottom = ship.RightBottomCell()

			gridRowStartIndex    = getStartGridIndex(leftTop.RowIndex(), ShipPlacingMargin)
			gridColumnStartIndex = getStartGridIndex(leftTop.ColumnIndex(), ShipPlacingMargin)
			gridRowEndIndex      = getEndGridIndex(rightBottom.RowIndex(), bf.size, ShipPlacingMargin)
			gridColumnEndIndex   = getEndGridIndex(rightBottom.ColumnIndex(), bf.size, ShipPlacingMargin)
		)

		for i := gridRowStartIndex; i <= gridRowEndIndex; i++ {
			for j := gridColumnStartIndex; j <= gridColumnEndIndex; j++ {
				if !bf.cells[i][j].IsEmpty() {
					return ErrCellAlreadyTaken
				}
			}
		}

		cells := make(
			[]*BattleFieldCell,
			0,
			(rightBottom.RowIndex()-leftTop.RowIndex()+1)*(rightBottom.ColumnIndex()-leftTop.ColumnIndex()+1))

		for i := leftTop.RowIndex(); i <= rightBottom.RowIndex(); i++ {
			for j := leftTop.ColumnIndex(); j <= rightBottom.ColumnIndex(); j++ {
				cells = append(cells, &bf.cells[i][j])
			}
		}

		ship.AttachCells(cells...)
		bf.PlaceShipEvent()
	}

	return nil
}

func getStartGridIndex(i, margin int) int {
	if i >= margin {
		return i - margin
	}

	return i
}

func getEndGridIndex(i, size, margin int) int {
	if i < size-1 {
		return i + margin
	}

	return i
}

func (bf *BattleField) Clear() {
	for i := 0; i < bf.size; i++ {
		for j := 0; j < bf.size; j++ {
			bf.cells[i][j].Clear()
		}
	}

	bf.ClearFieldEvent()
}

func (bf *BattleField) MakeShot(cell *BattleFieldCell) (destroyed, knocked bool, err error) {
	if cell.IsDestroyed() {
		return false, false, ErrCellAlreadyDestroyed
	}

	cell.Destroy()
	bf.ShotEvent()

	if cell.IsEmpty() {
		return false, false, nil
	}

	bf.KnockEvent()

	ship := cell.BattleShip()

	if destroyed = ship.IsDestroyed(); destroyed {
		ship.Destroy()
		bf.DestroyEvent()
	}

	return destroyed, true, nil
}

func (bf *BattleField) TakeCell(coordinates string) (cell *BattleFieldCell, err error) {
	i, j, err := convertToIndexes(coordinates)
	if err != nil {
		return nil, err
	}

	return &bf.cells[i][j], nil
}

func convertToIndexes(coordinates string) (i, j int, err error) {
	row := coordinates[:len(coordinates)-1]
	if i, err = strconv.Atoi(row); err != nil {
		return 0, 0, err
	}

	return i - 1, int(coordinates[len(row):][0] - 'A'), nil
}
