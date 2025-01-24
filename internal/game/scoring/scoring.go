package scoring

import (
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	"strings"
)

type Scorer interface {
	Score(board [][]string) *types.ScoreResult
}

type TxtDictScorer struct {
	matcher Matcher
}

func NewTxtDictScorer(matcher Matcher) *TxtDictScorer {
	return &TxtDictScorer{
		matcher: matcher,
	}
}

func (s *TxtDictScorer) Score(board [][]string) *types.ScoreResult {
	words := s.findScoringWords(board)
	total := 0
	for _, word := range words {
		total += word.Score
	}
	return &types.ScoreResult{
		TotalScore: total,
		Words:      words,
	}
}

func (s *TxtDictScorer) findScoringWords(board [][]string) []*types.ScoredWord {
	words := make([]*types.ScoredWord, 0)

	// Horizontal words
	for r := range len(board) {
		line := strings.Join(board[r], "")
		words = append(words, s.scoreWordsForLine(lineScoreInput{
			Line:      line,
			Direction: types.ScoringDirectionHorizontal,
			Row:       r,
			Column:    0,
		})...)
	}

	// Vertical words
	for c := range len(board) {
		var sb strings.Builder
		for r := range len(board) {
			sb.WriteString(board[r][c])
		}
		line := sb.String()
		words = append(words, s.scoreWordsForLine(lineScoreInput{
			Line:      line,
			Direction: types.ScoringDirectionVertical,
			Row:       0,
			Column:    c,
		})...)
	}

	return words
}

type lineScoreInput struct {
	Line      string
	Direction types.ScoringDirection
	Row       int
	Column    int
}

func (s *TxtDictScorer) scoreWordsForLine(
	input lineScoreInput,
) []*types.ScoredWord {
	// Ensure the line is uppercase since we store our dictionary that way
	input.Line = strings.ToUpper(input.Line)
	matchedWords := s.matcher.Match(input.Line)
	matchedLineIndices := matchedWordsToLineIndices(input.Line, matchedWords)
	bestScoringWords := getBestScoringWordCombination(input, matchedLineIndices)
	return bestScoringWords
}

type wordWithRange struct {
	Word     string
	StartIdx int
	// Ultimately we're checking for overlap, and with the board size limited to 7,
	// we may as well just keep a bit mask of used indices
	IncludedIndices uint32
}

func matchedWordsToLineIndices(line string, matchedWords []string) []wordWithRange {
	ranges := make([]wordWithRange, 0, len(matchedWords))
	for _, word := range matchedWords {
		i := 0
		for {
			relativeOccurrenceIndex := strings.Index(line[i:], word)
			if relativeOccurrenceIndex == -1 {
				// No more occurrences of the word
				break
			}
			occurrenceIdx := i + relativeOccurrenceIndex

			// Generate a bitmask for the indices of the word in the line
			includedIndices := uint32(0)
			for j := 0; j < len(word); j++ {
				includedIndices |= 1 << uint32(occurrenceIdx+j)
			}

			ranges = append(ranges, wordWithRange{
				Word:            word,
				StartIdx:        occurrenceIdx,
				IncludedIndices: includedIndices,
			})

			// The next place a copy of the word could occur is the character after the current occurrence
			// (e.g. consider a word `aaa` in a string `aaaaaa`)
			i = occurrenceIdx + 1
		}
	}
	return ranges
}

// The possible word combinations are, essentially, the powerset of our matched words
// But restricted to subsets which do not overlap in index
// We only care about the best scoring one
// Doing it recursively; with the line/board size limited to 7, this should be fine
func getBestScoringWordCombination(input lineScoreInput, words []wordWithRange) []*types.ScoredWord {
	var bestScoringSubset []*types.ScoredWord
	bestScore := 0
	getBestScoringWordCombinationRec(
		input,
		words,
		0,
		make([]wordWithRange, 0),
		&bestScoringSubset,
		&bestScore,
	)
	return bestScoringSubset
}

func getBestScoringWordCombinationRec(
	input lineScoreInput,
	inputSet []wordWithRange,
	currentIdx int,
	currentSubset []wordWithRange,
	out *[]*types.ScoredWord,
	bestScore *int,
) {
	// If the subset we've built is invalid due to overlap, discard this branch entirely
	if !isSubsetNonOverlapping(currentSubset) {
		return
	}

	// Base case, we're past the end of the input set
	if currentIdx == len(inputSet) {
		// Score this word combination
		total, scoredWords := scoreWordCombination(currentSubset, input)
		// Save if this is our new best score
		if total > *bestScore {
			*bestScore = total
			*out = scoredWords
		}
		return
	}

	// TODO: Commonality between recursive cases, calculate once

	// Recursive case 1: where the current element is in the set
	// (backtracking logic to avoid copying the subset here)
	currentSubset = append(currentSubset, inputSet[currentIdx])
	getBestScoringWordCombinationRec(input, inputSet, currentIdx+1, currentSubset, out, bestScore)
	currentSubset = currentSubset[:len(currentSubset)-1]

	// Recursive case 2: where the current element is not in the set
	getBestScoringWordCombinationRec(input, inputSet, currentIdx+1, currentSubset, out, bestScore)

}

func isSubsetNonOverlapping(subset []wordWithRange) bool {
	// Track which indices have been used, and fail if there is overlap
	usedIndices := uint32(0)
	for _, word := range subset {
		if word.IncludedIndices&usedIndices != 0 {
			return false
		}
		usedIndices |= word.IncludedIndices
	}
	return true
}

func scoreWordCombination(words []wordWithRange, input lineScoreInput) (int, []*types.ScoredWord) {
	total := 0
	scoredWords := make([]*types.ScoredWord, 0, len(words))
	for _, word := range words {
		wordScore := scoreWord(word.Word, len(input.Line))
		total += wordScore
		currentScoredWord := types.ScoredWord{
			Word:      word.Word,
			Score:     wordScore,
			Direction: input.Direction,
		}
		if currentScoredWord.Direction == types.ScoringDirectionHorizontal {
			currentScoredWord.StartRow = input.Row
			currentScoredWord.StartColumn = word.StartIdx
		} else {
			currentScoredWord.StartRow = word.StartIdx
			currentScoredWord.StartColumn = input.Column
		}

		scoredWords = append(scoredWords, &currentScoredWord)
	}
	return total, scoredWords
}

func scoreWord(word string, boardDimension int) int {
	if len(word) == boardDimension {
		return len(word) * 2
	}
	return len(word)
}
