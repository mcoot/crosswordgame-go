package scoring

import (
	"github.com/mcoot/crosswordgame-go/internal/game/scoring/matching"
	"github.com/mcoot/crosswordgame-go/internal/game/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type ScoringSuite struct {
	suite.Suite
}

func TestScoringSuite(t *testing.T) {
	suite.Run(t, new(ScoringSuite))
}

type testCase struct {
	name       string
	dictionary []string
	line       string
	direction  types.ScoringDirection
	row        int
	column     int
	expect     ScoredWordsExpectation
}

func (s *ScoringSuite) Test_AhoCorasickMatcher_scoreWordsForLine() {
	cases := []testCase{
		{
			name:       "when empty dict, no matches",
			dictionary: []string{},
			line:       "hello",
			expect:     expectExactly([]*types.ScoredWord{}),
		},
		{
			name: "when empty line, no matches",
			dictionary: []string{
				"go",
				"car",
				"cargo",
			},
			line:      "",
			direction: types.ScoringDirectionHorizontal,
			expect:    expectExactly([]*types.ScoredWord{}),
		},
		{
			name: "does not match single letter words",
			dictionary: []string{
				"a",
				"b",
				"c",
			},
			line:      "dabce",
			direction: types.ScoringDirectionHorizontal,
			expect:    expectWordsToBe([]string{}),
		},
		{
			name: "partial match at end",
			dictionary: []string{
				"go",
				"car",
				"cargo",
			},
			line:      "caago",
			direction: types.ScoringDirectionHorizontal,
			expect: expectWordsToBe([]string{
				"go",
			}),
		},
		{
			name: "partial match at beginning",
			dictionary: []string{
				"go",
				"car",
				"cargo",
			},
			line:      "cargn",
			direction: types.ScoringDirectionHorizontal,
			expect: expectWordsToBe([]string{
				"car",
			}),
		},
		{
			name: "partial match in the middle",
			dictionary: []string{
				"go",
				"car",
				"cargo",
			},
			line:      "ccarg",
			direction: types.ScoringDirectionHorizontal,
			expect: expectWordsToBe([]string{
				"car",
			}),
		},
		{
			name: "4 letter match",
			dictionary: []string{
				"cash",
			},
			line:      "acash",
			direction: types.ScoringDirectionHorizontal,
			expect: expectWordsToBe([]string{
				"cash",
			}),
		},
		{
			name: "matches two non-overlapping segments",
			dictionary: []string{
				"yes",
				"in",
			},
			line:      "yesin",
			direction: types.ScoringDirectionHorizontal,
			expect: expectWordsToBe([]string{
				"yes",
				"in",
			}),
		},
		{
			name: "does not match overlap, prioritises longer match",
			dictionary: []string{
				"ca",
				"ango",
			},
			line:      "cango",
			direction: types.ScoringDirectionHorizontal,
			expect: expectWordsToBe([]string{
				"ango",
			}),
		},
		{
			name: "prioritises first match in overlap if equal length",
			dictionary: []string{
				"can",
				"ngo",
			},
			line:      "cango",
			direction: types.ScoringDirectionHorizontal,
			expect: expectWordsToBe([]string{
				"can",
			}),
		},
		{
			name: "matches the same word twice if non-overlapping",
			dictionary: []string{
				"to",
			},
			line:      "totoa",
			direction: types.ScoringDirectionHorizontal,
			expect: expectWordsToBe([]string{
				"to",
				"to",
			}),
		},
		{
			name: "does not re-match the same word with overlap",
			dictionary: []string{
				"tt",
			},
			line:      "tttaa",
			direction: types.ScoringDirectionHorizontal,
			expect: expectWordsToBe([]string{
				"tt",
			}),
		},

		{
			name: "prioritises total match over two segments",
			dictionary: []string{
				"go",
				"car",
				"cargo",
			},
			line:      "cargo",
			direction: types.ScoringDirectionHorizontal,
			expect: expectWordsToBe([]string{
				"cargo",
			}),
		},
		{
			name: "prioritises longer total length over individual longer word",
			dictionary: []string{
				"to",
				"tot",
			},
			line:      "totoa",
			direction: types.ScoringDirectionHorizontal,
			expect: expectWordsToBe([]string{
				"to",
				"to",
			}),
		},

		{
			name: "for horizontal words, gives the correct row and column",
			dictionary: []string{
				"car",
			},
			line:      "xcarx",
			direction: types.ScoringDirectionHorizontal,
			row:       3,
			expect: expectExactly([]*types.ScoredWord{
				{
					Word:        "car",
					Score:       3,
					Direction:   types.ScoringDirectionHorizontal,
					StartRow:    3,
					StartColumn: 1,
				},
			}),
		},
		{
			name: "for vertical words, gives the correct row and column",
			dictionary: []string{
				"car",
			},
			line:      "xcarx",
			direction: types.ScoringDirectionVertical,
			column:    2,
			expect: expectExactly([]*types.ScoredWord{
				{
					Word:        "car",
					Score:       3,
					Direction:   types.ScoringDirectionVertical,
					StartRow:    1,
					StartColumn: 2,
				},
			}),
		},
		{
			name: "scores a full match as double points",
			dictionary: []string{
				"car",
				"go",
				"cargo",
			},
			line:      "cargo",
			direction: types.ScoringDirectionVertical,
			column:    2,
			expect: expectExactly([]*types.ScoredWord{
				{
					Word:        "cargo",
					Score:       10,
					Direction:   types.ScoringDirectionVertical,
					StartRow:    0,
					StartColumn: 2,
				},
			}),
		},
		{
			name: "scores words in appearance order",
			dictionary: []string{
				"car",
				"to",
			},
			line:      "tocar",
			direction: types.ScoringDirectionHorizontal,
			row:       4,
			expect: expectExactly([]*types.ScoredWord{
				{
					Word:        "to",
					Score:       2,
					Direction:   types.ScoringDirectionHorizontal,
					StartRow:    4,
					StartColumn: 0,
				},
				{
					Word:        "car",
					Score:       3,
					Direction:   types.ScoringDirectionHorizontal,
					StartRow:    4,
					StartColumn: 2,
				},
			}),
		},
	}

	for _, c := range cases {
		s.T().Run(c.name, func(t *testing.T) {
			words := make([]string, 0, len(c.dictionary))
			for _, word := range c.dictionary {
				words = append(words, strings.ToUpper(word))
			}
			matcher := matching.NewAhoCorasickMatcher(words)
			scorer := NewTxtDictScorer(matcher)
			got := scorer.scoreWordsForLine(lineScoreInput{
				Line:      c.line,
				Direction: c.direction,
				Row:       c.row,
				Column:    c.column,
			})
			c.expect(t, got)
		})
	}
}

type ScoredWordsExpectation func(t *testing.T, sw []*types.ScoredWord)

func expectWordsToBe(expectedWords []string) ScoredWordsExpectation {
	return func(t *testing.T, sw []*types.ScoredWord) {
		t.Helper()
		require.Len(t, sw, len(expectedWords))
		for i, word := range expectedWords {
			require.Equal(t, strings.ToUpper(word), sw[i].Word)
		}
	}
}

func expectExactly(expected []*types.ScoredWord) ScoredWordsExpectation {
	return func(t *testing.T, sw []*types.ScoredWord) {
		t.Helper()
		require.Len(t, sw, len(expected))
		for i, word := range expected {
			word.Word = strings.ToUpper(word.Word)
			require.Equal(t, word, sw[i])
		}
	}
}
