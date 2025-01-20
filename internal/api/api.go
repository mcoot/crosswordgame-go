package api

import (
	"github.com/gorilla/mux"
	"github.com/mcoot/crosswordgame-go/internal/api/jsonapi"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/game/scoring"
	"github.com/mcoot/crosswordgame-go/internal/lobby"
	"github.com/mcoot/crosswordgame-go/internal/middleware"
	"github.com/mcoot/crosswordgame-go/internal/store"
	"github.com/tomarrell/wrapcheck/v2/wrapcheck/testdata/ignore_pkg_errors/src/github.com/pkg/errors"
	"go.uber.org/zap"
)

func SetupAPI(logger *zap.SugaredLogger, db store.Store, schemaPath string, dictPath string) (*mux.Router, error) {
	router := mux.NewRouter()

	err := middleware.SetupMiddleware(router, logger, schemaPath)
	if err != nil {
		return nil, errors.Wrap(err, "error setting up middleware")
	}

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	gameScorer, err := scoring.NewTxtDictScorer(dictPath)
	if err != nil {
		return nil, errors.Wrap(err, "error creating building dictionary scorer")
	}
	gameManager := game.NewGameManager(db, gameScorer)
	lobbyManager := lobby.NewLobbyManager(db)

	jsonApi := jsonapi.NewCrosswordGameAPI(gameManager, lobbyManager)
	err = jsonApi.AttachToRouter(apiRouter)
	if err != nil {
		return nil, errors.Wrap(err, "error attaching JSON API to router")
	}
	webApi := webapi.NewCrosswordGameWebAPI(gameManager, lobbyManager)
	err = webApi.AttachToRouter(router)
	if err != nil {
		return nil, errors.Wrap(err, "error attaching web API to router")
	}

	return router, nil
}
