package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi"
	seabattle "github.com/morozovcookie/sea-battle"
	"go.uber.org/zap"
)

const (
	SeaBattleSvcPathPrefix             = "/"
	SeaBattleSvcCreateMatrixPathPrefix = "/create-matrix"
	SeaBattleSvcShipPathPrefix         = "/ship"
	SeaBattleSvcShotPathPrefix         = "/shot"
	SeaBattleSvcClearPathPrefix        = "/clear"
	SeaBattleSvcStatePathPrefix        = "/state"
)

type SeaBattleService struct {
	router chi.Router

	logger *zap.Logger

	game *seabattle.Game
}

func NewSeaBattleService(logger *zap.Logger) (svc *SeaBattleService) {
	svc = &SeaBattleService{
		router: chi.NewRouter(),

		logger: logger,
	}

	svc.router.Post(SeaBattleSvcCreateMatrixPathPrefix, svc.startGameChecker(svc.CreateMatrixHandler))
	svc.router.Post(SeaBattleSvcShipPathPrefix, svc.nilGameChecker(svc.endGameChecker(svc.ShipHandler)))
	svc.router.Post(SeaBattleSvcShotPathPrefix, svc.nilGameChecker(svc.shipPlacingChecker(svc.ShotHandler)))
	svc.router.Post(SeaBattleSvcClearPathPrefix, svc.nilGameChecker(svc.ClearHandler))
	svc.router.Get(SeaBattleSvcStatePathPrefix, svc.nilGameChecker(svc.StateHandler))

	return svc
}

func (svc *SeaBattleService) CreateMatrixHandler(w http.ResponseWriter, r *http.Request) {
	var (
		req = &struct {
			Range int `json:"range"`
		}{}

		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		svc.writeError(w, http.StatusBadRequest, err)

		return
	}

	creator := func() (*seabattle.BattleField, error) {
		return seabattle.NewBattleField(req.Range)
	}

	if svc.game, err = seabattle.NewGame(creator); err != nil {
		svc.writeError(w, http.StatusBadRequest, err)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (svc *SeaBattleService) ShipHandler(w http.ResponseWriter, r *http.Request) {
	if svc.game.IsOver() {
		svc.game.BattleField().Clear()
	}

	var (
		req = &struct {
			Coordinates []string `json:"Coordinates"`
		}{}
		validator = func(coordinates string) error {
			var (
				size  = svc.game.BattleField().Size()
				regex = fmt.Sprintf(
					`^[1-%d][A-%c]\s[1-%d][A-%c]$`, size, byte('A'+size-1), size, byte('A'+size-1))
			)
			if size == seabattle.MaxBattleFieldSize {
				regex = `^((10)|([1-9]))[A-J]\s((10)|([1-9]))[A-J]$`
			}

			if !regexp.MustCompile(regex).MatchString(coordinates) {
				return seabattle.ErrInvalidCoordinates
			}

			return nil
		}
	)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		svc.writeError(w, http.StatusBadRequest, err)

		return
	}

	ships := make([]*seabattle.BattleShip, 0, len(req.Coordinates))

	for _, coord := range req.Coordinates {
		if err := validator(coord); err != nil {
			svc.writeError(w, http.StatusBadRequest, err)

			return
		}

		cells := strings.Split(coord, " ")

		leftTop, err := svc.game.BattleField().TakeCell(cells[0])
		if err != nil {
			svc.writeError(w, http.StatusBadRequest, err)

			return
		}

		rightBottom, err := svc.game.BattleField().TakeCell(cells[1])
		if err != nil {
			svc.writeError(w, http.StatusBadRequest, err)

			return
		}

		if leftTop.RowIndex() > rightBottom.RowIndex() || leftTop.ColumnIndex() > rightBottom.ColumnIndex() {
			svc.writeError(w, http.StatusBadRequest, seabattle.ErrInvalidCoordinates)

			return
		}

		ships = append(ships, seabattle.NewBattleShip(leftTop, rightBottom))
	}

	if err := svc.game.BattleField().PlaceBattleShips(ships...); err != nil {
		svc.writeError(w, http.StatusBadRequest, err)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (svc *SeaBattleService) ShotHandler(w http.ResponseWriter, r *http.Request) {
	var (
		req = &struct {
			Coordinates string `json:"сoord"`
		}{}
		validator = func(coordinates string) error {
			var (
				size  = svc.game.BattleField().Size()
				regex = fmt.Sprintf(`^[1-%d][A-%c]$`, size, byte('A'+size-1))
			)

			if size == seabattle.MaxBattleFieldSize {
				regex = `^((10)|([1-9]))[A-J]$`
			}

			if !regexp.MustCompile(regex).MatchString(coordinates) {
				return seabattle.ErrInvalidCoordinates
			}

			return nil
		}

		err error
	)

	if err = json.NewDecoder(r.Body).Decode(req); err != nil {
		svc.writeError(w, http.StatusBadRequest, err)

		return
	}

	if err = validator(req.Coordinates); err != nil {
		svc.writeError(w, http.StatusBadRequest, err)

		return
	}

	resp := &struct {
		Destroyed bool `json:"destroy"`
		Knocked   bool `json:"knock"`
		Ended     bool `json:"end"`
	}{}

	if svc.game.IsOver() {
		svc.writeError(w, http.StatusBadRequest, seabattle.ErrAllShipsAlreadyDestroyed)

		return
	}

	cell, err := svc.game.BattleField().TakeCell(req.Coordinates)
	if err != nil {
		svc.writeError(w, http.StatusBadRequest, err)

		return
	}

	if resp.Destroyed, resp.Knocked, err = svc.game.BattleField().MakeShot(cell); err != nil {
		svc.writeError(w, http.StatusBadRequest, err)

		return
	}

	resp.Ended = svc.game.IsOver()

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		svc.writeError(w, http.StatusInternalServerError, err)
	}
}

func (svc *SeaBattleService) ClearHandler(w http.ResponseWriter, _ *http.Request) {
	svc.game = nil

	w.WriteHeader(http.StatusNoContent)
}

func (svc *SeaBattleService) StateHandler(w http.ResponseWriter, _ *http.Request) {
	resp := &struct {
		ShipCount int `json:"ship_count"`
		Destroyed int `json:"destroyed"`
		Knocked   int `json:"knocked"`
		ShotCount int `json:"shot_count"`
	}{
		ShipCount: int(svc.game.State().ShipCount()),
		Destroyed: int(svc.game.State().Destroyed()),
		Knocked:   int(svc.game.State().Knocked()),
		ShotCount: int(svc.game.State().ShotCount()),
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (svc *SeaBattleService) nilGameChecker(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if svc.game == nil {
			svc.writeError(w, http.StatusBadRequest, seabattle.ErrGameDidNotStarted)

			return
		}

		next.ServeHTTP(w, r)
	}
}

func (svc *SeaBattleService) startGameChecker(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if svc.game != nil {
			svc.writeError(w, http.StatusBadRequest, seabattle.ErrGameAlreadyStarted)

			return
		}

		next.ServeHTTP(w, r)
	}
}

func (svc *SeaBattleService) endGameChecker(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if svc.game.State().ShipCount() > 0 && !svc.game.IsOver() {
			svc.writeError(w, http.StatusBadRequest, seabattle.ErrClearFieldBeforePlaceNewShips)

			return
		}

		next.ServeHTTP(w, r)
	}
}

func (svc *SeaBattleService) shipPlacingChecker(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if svc.game.State().ShipCount() == 0 {
			svc.writeError(w, http.StatusBadRequest, seabattle.ErrNoShipWasNotPlaced)

			return
		}

		next.ServeHTTP(w, r)
	}
}

func (svc *SeaBattleService) writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	if _, wErr := w.Write(append([]byte{}, err.Error()...)); wErr != nil {
		svc.logger.Error("write response error", zap.Error(wErr))
	}
}

func (svc *SeaBattleService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	svc.router.ServeHTTP(w, r)
}
