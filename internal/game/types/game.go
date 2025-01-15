package types

type GameId string

type Status string

const (
	StatusAwaitingAnnouncement Status = "awaiting_announcement"
	StatusAwaitingPlacement    Status = "awaiting_placement"
	StatusFinished             Status = "finished"
)

type GameState struct {
	Status                  Status
	SquaresFilled           int
	CurrentAnnouncingPlayer int
	PlayerCount             int
}

type Game struct {
	GameState
	Players []*Player
}

func NewGame(playerCount int) *Game {
	players := make([]*Player, playerCount)
	for i := range players {
		players[i] = NewPlayer()
	}

	return &Game{
		GameState: GameState{
			Status:                  StatusAwaitingAnnouncement,
			PlayerCount:             playerCount,
			SquaresFilled:           0,
			CurrentAnnouncingPlayer: 0,
		},
		Players: players,
	}
}
