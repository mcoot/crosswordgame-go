package types

type Player struct {
	Board *Board
}

func NewPlayer(boardDimension int) *Player {
	return &Player{
		Board: NewBoard(boardDimension),
	}
}
