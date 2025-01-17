package game

import (
	"bufio"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	"os"
)

type Scorer interface {
	Score(player *types.Player) (int, []*types.ScoredWord)
}

type TxtDictScorer struct {
	wordList map[string]bool
}

func NewTxtDictScorer() *TxtDictScorer {
	return &TxtDictScorer{
		wordList: make(map[string]bool),
	}
}

func (s *TxtDictScorer) LoadDictionary(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s.wordList[scanner.Text()] = true
	}

	return scanner.Err()
}

func (s *TxtDictScorer) Score(player *types.Player) (int, []*types.ScoredWord) {
	words := s.findScoringWords(player.Board)
	total := 0
	for _, word := range words {
		total += word.Score
	}
	return total, words
}

func (s *TxtDictScorer) findScoringWords(board *types.Board) []*types.ScoredWord {
	return []*types.ScoredWord{}
}
