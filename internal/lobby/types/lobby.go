package types

import playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"

type LobbyId string

type Lobby struct {
	Name        string
	Players     []playertypes.PlayerId
	RunningGame *RunningGame
}
