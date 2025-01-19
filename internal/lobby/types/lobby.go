package types

type LobbyId string

type Lobby struct {
	Name        string
	Players     map[PlayerId]Player
	RunningGame *RunningGame
}
