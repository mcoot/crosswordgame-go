package matching

import "github.com/cloudflare/ahocorasick"

type AhoCorasickMatcher struct {
	wordList []string
	matcher  *ahocorasick.Matcher
}

func NewAhoCorasickMatcher(filename string) (*AhoCorasickMatcher, error) {
	wordList, err := LoadDictionary(50000, filename)
	if err != nil {
		return nil, err
	}
	matcher := ahocorasick.NewStringMatcher(wordList)

	return &AhoCorasickMatcher{
		wordList: wordList,
		matcher:  matcher,
	}, nil
}

func (d *AhoCorasickMatcher) Match(line string) []string {
	matchIndices := d.matcher.MatchThreadSafe([]byte(line))
	words := make([]string, 0, len(matchIndices))
	for _, matchIdx := range matchIndices {
		word := d.wordList[matchIdx]
		words = append(words, word)
	}
	return words
}
