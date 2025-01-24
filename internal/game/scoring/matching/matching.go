package matching

import (
	"bufio"
	"os"
	"strings"
)

type Matcher interface {
	Match(line string) []string
}

func LoadDictionary(preallocatedCapacity int, filename string) ([]string, error) {
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
