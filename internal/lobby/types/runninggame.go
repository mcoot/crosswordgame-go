package types

import (
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
)

type RunningGame struct {
	PlayerIdToIdx map[playertypes.PlayerId]int
	GameId        types.GameId
}
