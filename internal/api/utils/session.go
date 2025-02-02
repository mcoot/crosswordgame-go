package utils

import (
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/mcoot/crosswordgame-go/internal/errors"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	"github.com/mcoot/crosswordgame-go/internal/player"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"github.com/mcoot/crosswordgame-go/internal/utils"
	"net/http"
)

const (
	ContextKeySession utils.ContextKey = "session"
)

// Session is the baseSession enriched with lookups from the database
type Session struct {
	*sessions.Session
	Player *playertypes.Player
	Lobby  *lobbytypes.Lobby
}

type SessionManager struct {
	sessionStore sessions.Store
}

func (s *Session) IsLoggedIn() bool {
	return s.Player != nil
}

func (s *Session) IsInLobby() bool {
	return s.Lobby != nil
}

func NewSessionManager(sessionStore sessions.Store) *SessionManager {
	return &SessionManager{
		sessionStore: sessionStore,
	}
}

func (sm *SessionManager) GetSession(r *http.Request, playerManager *player.Manager) (*Session, error) {
	session, err := sm.sessionStore.Get(r, "session")
	if err != nil {
		return nil, err
	}

	rawPlayerId, ok := session.Values["player_id"]
	if !ok {
		// If session is missing player_id, we just aren't logged in
		return &Session{
			Session: session,
			Player:  nil,
			Lobby:   nil,
		}, nil
	}
	strPlayerId, ok := rawPlayerId.(string)
	if !ok {
		return nil, fmt.Errorf("malformed session: player_id is not a string")
	}
	playerId := playertypes.PlayerId(strPlayerId)

	p, err := playerManager.LookupPlayer(playerId)
	if err != nil {
		// If the session has an invalid player_id, treat it as just not being logged in
		if errors.IsNotFoundError(err) {
			return &Session{
				Session: session,
				Player:  nil,
				Lobby:   nil,
			}, nil
		} else {
			return nil, err
		}
	}

	lobby, err := playerManager.GetLobbyForPlayer(playerId)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return &Session{
				Session: session,
				Player:  p,
				Lobby:   nil,
			}, nil
		} else {
			return nil, err
		}
	}

	return &Session{
		Session: session,
		Player:  p,
		Lobby:   lobby,
	}, nil
}

func (sm *SessionManager) SaveLoggedInPlayer(
	w http.ResponseWriter,
	r *http.Request,
	playerId playertypes.PlayerId,
) error {
	session, err := sm.sessionStore.Get(r, "session")
	if err != nil {
		return err
	}
	session.Values["player_id"] = string(playerId)

	return sm.sessionStore.Save(r, w, session)
}

func AddSessionToContext(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, ContextKeySession, session)
}

func GetSessionFromContext(ctx context.Context) (*Session, error) {
	session, ok := ctx.Value(ContextKeySession).(*Session)
	if !ok {
		return nil, fmt.Errorf("session not found in context")
	}
	return session, nil
}
