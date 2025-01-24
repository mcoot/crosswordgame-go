package matching

import "github.com/cloudflare/ahocorasick"

type AhoCorasickMatcher struct {
	wordList []string
	matcher  *ahocorasick.Matcher
}

func NewAhoCorasickMatcher(wordList []string) *AhoCorasickMatcher {
	result := &AhoCorasickMatcher{}
	result.wordList = getFilteredDictionary(wordList)

	result.matcher = ahocorasick.NewStringMatcher(result.wordList)
	return result
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
