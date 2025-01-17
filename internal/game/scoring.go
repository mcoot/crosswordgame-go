package game

import "github.com/mcoot/crosswordgame-go/internal/game/types"

func determineScore(player *types.Player) (int, []*types.ScoredWord) {
	words := findScoringWords(player.Board)
	total := 0
	for _, word := range words {
		total += word.Score
	}
	return total, words
}

func findScoringWords(board *types.Board) []*types.ScoredWord {
	return []*types.ScoredWord{}
}
