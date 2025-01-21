package utils

import (
	"fmt"
	"github.com/gorilla/sessions"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"net/http"
)

type Session struct {
	*sessions.Session
	PlayerId playertypes.PlayerId
}

func (s *Session) IsLoggedIn() bool {
	return s.PlayerId != ""
}

func GetSessionDetails(sessionStore sessions.Store, r *http.Request) (*Session, error) {
	session, err := sessionStore.Get(r, "session")
	if err != nil {
		return nil, err
	}

	rawPlayerId, ok := session.Values["player_id"]
	if !ok {
		// If session is missing player_id, we just aren't logged in
		return &Session{
			Session: session,
		}, nil
	}
	strPlayerId, ok := rawPlayerId.(string)
	if !ok {
		return nil, fmt.Errorf("malformed session: player_id is not a string")
	}
	playerId := playertypes.PlayerId(strPlayerId)

	return &Session{
		Session:  session,
		PlayerId: playerId,
	}, nil
}

func SetSession(store sessions.Store, session *Session, w http.ResponseWriter, r *http.Request) error {
	session.Session.Values["player_id"] = session.PlayerId

	err := store.Save(r, w, session.Session)
	if err != nil {
		return err
	}
	return nil
}
