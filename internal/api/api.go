package api

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/mcoot/crosswordgame-go/internal/api/jsonapi"
	"github.com/mcoot/crosswordgame-go/internal/api/utils"
	"github.com/mcoot/crosswordgame-go/internal/api/webapi"
	"github.com/mcoot/crosswordgame-go/internal/game"
	"github.com/mcoot/crosswordgame-go/internal/game/scoring"
	"github.com/mcoot/crosswordgame-go/internal/lobby"
	"github.com/mcoot/crosswordgame-go/internal/player"
	"github.com/mcoot/crosswordgame-go/internal/store"
	"github.com/tomarrell/wrapcheck/v2/wrapcheck/testdata/ignore_pkg_errors/src/github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

func SetupAPI(
	logger *zap.SugaredLogger,
	db store.Store,
	sessionStore sessions.Store,
	schemaPath string,
	dictPath string,
) (http.Handler, error) {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	sessionManager := utils.NewSessionManager(sessionStore)

	logger.Infow("Loading dictionary")
	gameScorer, err := scoring.NewTxtDictScorer(dictPath)
	if err != nil {
		return nil, errors.Wrap(err, "error creating building dictionary scorer")
	}
	gameManager := game.NewGameManager(db, gameScorer)
	lobbyManager := lobby.NewLobbyManager(db)
	playerManager := player.NewPlayerManager(db)

	logger.Infow("Initialising APIs")

	jsonApi := jsonapi.NewCrosswordGameAPI(gameManager, lobbyManager, playerManager)
	err = jsonApi.AttachToRouter(apiRouter, logger, schemaPath)
	if err != nil {
		return nil, errors.Wrap(err, "error attaching JSON API to router")
	}
	webApi := webapi.NewCrosswordGameWebAPI(sessionManager, gameManager, lobbyManager, playerManager)
	err = webApi.AttachToRouter(router)
	if err != nil {
		return nil, errors.Wrap(err, "error attaching web API to router")
	}

	handler := SetupGlobalMiddleware(router, logger)

	return handler, nil
}
