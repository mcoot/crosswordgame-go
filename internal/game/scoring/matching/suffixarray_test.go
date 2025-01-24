package matching

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type SuffixArraySuite struct {
	suite.Suite
}

func TestSuffixArraySuite(t *testing.T) {
	suite.Run(t, new(SuffixArraySuite))
}

func (s *SuffixArraySuite) Test_SuffixArrayMatcher_Match() {
	type testCase struct {
		name       string
		dictionary []string
		line       string
		expect     []string
	}

	cases := []testCase{
		{
			name:       "when empty dict, no matches",
			dictionary: []string{},
			line:       "hello",
			expect:     []string{},
		},
		{
			name: "when word not in dict, no matches",
			dictionary: []string{
				"apple",
				"banana",
				"cherry",
				"donut",
				"eggplant",
			},
			line:   "hello",
			expect: []string{},
		},

		{
			name: "whole word match, start of dict",
			dictionary: []string{
				"apple",
				"banana",
				"cherry",
				"donut",
				"eggplant",
				"pineapple",
			},
			line: "apple",
			expect: []string{
				"apple",
			},
		},
		{
			name: "whole word match, mid-dict",
			dictionary: []string{
				"apple",
				"banana",
				"cherry",
				"donut",
				"eggplant",
				"pineapple",
			},
			line: "cherry",
			expect: []string{
				"cherry",
			},
		},
		{
			name: "whole word match, end of dict",
			dictionary: []string{
				"apple",
				"banana",
				"cherry",
				"donut",
				"pineapple",
				"eggplant",
			},
			line: "eggplant",
			expect: []string{
				"eggplant",
			},
		},
		{
			name: "whole and partial match",
			dictionary: []string{
				"apple",
				"banana",
				"cherry",
				"donut",
				"eggplant",
				"pineapple",
			},
			line: "pineapple",
			expect: []string{
				"pineapple",
				"apple",
			},
		},

		{
			name: "prefix match",
			dictionary: []string{
				"apple",
				"banana",
				"cherry",
				"donut",
				"eggplant",
				"pineapple",
			},
			line: "cherryooo",
			expect: []string{
				"cherry",
			},
		},
		{
			name: "suffix match",
			dictionary: []string{
				"apple",
				"banana",
				"cherry",
				"donut",
				"eggplant",
				"pineapple",
			},
			line: "ooodonut",
			expect: []string{
				"donut",
			},
		},
		{
			name: "internal match",
			dictionary: []string{
				"apple",
				"banana",
				"cherry",
				"donut",
				"eggplant",
				"pineapple",
			},
			line: "xxxbananaxxx",
			expect: []string{
				"banana",
			},
		},
		{
			name: "multiple sub-matches",
			dictionary: []string{
				"apple",
				"banana",
				"cherry",
				"donut",
				"eggplant",
				"pineapple",
			},
			line: "applexbananaxcherryxeggplantxpineapple",
			expect: []string{
				"apple",
				"banana",
				"cherry",
				"eggplant",
				"pineapple",
				// Apple appears a second time as a substring of pineapple
				// (our scoring logic deduplicates these kinds of matches)
				"apple",
			},
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			matcher := NewSuffixArrayMatcher(tc.dictionary)
			actual := matcher.Match(tc.line)
			s.Equal(tc.expect, actual)
		})
	}
}
