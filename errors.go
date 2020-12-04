package seabattle

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidBattleFieldSize        = errors.New("invalid battle field size")
	ErrInvalidCoordinates            = errors.New("invalid coordinates")
	ErrCellAlreadyDestroyed          = errors.New("ship already destroyed")
	ErrAllShipsAlreadyDestroyed      = errors.New("all ships already destroyed")
	ErrClearFieldBeforePlaceNewShips = errors.New("destroy all ships or clear battlefield before place new " +
		"ships")
	ErrGameDidNotStarted  = errors.New("game did not started")
	ErrCellAlreadyTaken   = errors.New("cell already taken")
	ErrGameAlreadyStarted = errors.New("game already started")
	ErrNoShipWasNotPlaced = errors.New("no ship was not placed")
)
