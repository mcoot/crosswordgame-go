package types

import "github.com/mcoot/crosswordgame-go/internal/game/types"

type RunningGame struct {
	PlayerIdToIdx map[PlayerId]int
	GameId        types.GameId
}
