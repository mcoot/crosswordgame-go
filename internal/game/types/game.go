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
	PlayerCount             int
	SquaresFilled           int
	BoardDimension          int
	CurrentAnnouncingPlayer int
	CurrentAnnouncedLetter  string
}

type Game struct {
	GameState
	Players []*Player
}

func NewGame(playerCount int, boardDimension int) *Game {
	players := make([]*Player, playerCount)
	for i := range players {
		players[i] = NewPlayer(boardDimension)
	}

	return &Game{
		GameState: GameState{
			Status:                  StatusAwaitingAnnouncement,
			PlayerCount:             playerCount,
			SquaresFilled:           0,
			BoardDimension:          boardDimension,
			CurrentAnnouncingPlayer: 0,
			CurrentAnnouncedLetter:  "",
		},
		Players: players,
	}
}

func (g *Game) TotalSquares() int {
	return g.BoardDimension * g.BoardDimension
}
