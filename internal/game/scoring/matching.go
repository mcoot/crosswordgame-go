package scoring

import (
	"bufio"
	"github.com/cloudflare/ahocorasick"
	"os"
	"strings"
)

type dictMatcher struct {
	wordList []string
	matcher  *ahocorasick.Matcher
}

func newDictMatcher(filename string) (*dictMatcher, error) {
	wordList, err := loadDictionary(50000, filename)
	if err != nil {
		return nil, err
	}
	matcher := ahocorasick.NewStringMatcher(wordList)

	return &dictMatcher{
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
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			wordList = append(wordList, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return wordList, nil
}

func (d *dictMatcher) Match(line string) []string {
	matchIndices := d.matcher.MatchThreadSafe([]byte(line))
	words := make([]string, 0, len(matchIndices))
	for _, matchIdx := range matchIndices {
		word := d.wordList[matchIdx]
		words = append(words, word)
	}
	return words
}
