package matching

import (
	"bufio"
	"github.com/cloudflare/ahocorasick"
	"os"
	"strings"
)

type Matcher interface {
	Match(line string) []string
}

type AhoCorasickMatcher struct {
	wordList []string
	matcher  *ahocorasick.Matcher
}

func NewAhoCorasickMatcher(filename string) (*AhoCorasickMatcher, error) {
	wordList, err := loadDictionary(50000, filename)
	if err != nil {
		return nil, err
	}
	matcher := ahocorasick.NewStringMatcher(wordList)

	return &AhoCorasickMatcher{
		wordList: wordList,
		matcher:  matcher,
	}, nil
}

func loadDictionary(preallocatedCapacity int, filename string) ([]string, error) {
	wordList := make([]string, 0, preallocatedCapacity)

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// Trim space and convert all words to uppercase as that's how we store boards
		word := strings.ToUpper(strings.TrimSpace(scanner.Text()))
		// Skip empty lines and one-letter words, since they don't score
		if len(word) > 1 {
			wordList = append(wordList, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return wordList, nil
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
