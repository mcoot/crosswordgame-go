package types

import (
	"github.com/hashicorp/go-uuid"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
)

type LobbyId string

type Lobby struct {
	Id          LobbyId
	Name        string
	Players     []playertypes.PlayerId
	RunningGame *RunningGame
}

func NewLobby(name string) (*Lobby, error) {
	rawId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	id := LobbyId(rawId)

	return &Lobby{
		Id:          id,
		Name:        name,
		Players:     make([]playertypes.PlayerId, 0),
		RunningGame: nil,
	}, nil
}
