package types

type Player struct {
	Board *Board
}

func NewPlayer() *Player {
	return &Player{
		Board: NewBoard(5),
	}
}
