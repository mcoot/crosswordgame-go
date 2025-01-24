package matching

import (
	"index/suffixarray"
)

type SuffixArrayMatcher struct {
	wordList   []string
	wordBuffer []byte
	index      *suffixarray.Index
}

func NewSuffixArrayMatcher(wordList []string) *SuffixArrayMatcher {
	result := &SuffixArrayMatcher{
		wordList: getFilteredDictionary(wordList),
	}

	result.buildIndex()
	return result
}

func (m *SuffixArrayMatcher) buildIndex() {
	// Calculate total buffer size needed
	// We keep an extra null terminator at the start
	// because we want to match on \0word\0
	totalLen := 1
	for _, w := range m.wordList {
		totalLen += len(w) + 1 // +1 for null separator
	}

	// Create buffer with null separators between words
	buf := make([]byte, 0, totalLen)
	for _, w := range m.wordList {
		buf = append(buf, 0) // null separator
		buf = append(buf, []byte(w)...)
	}
	buf = append(buf, 0)

	m.wordBuffer = buf
	m.index = suffixarray.New(m.wordBuffer)
}

func (m *SuffixArrayMatcher) Match(line string) []string {
	lineBytes := []byte(line)

	results := make([]string, 0)

	// Search over all substrings of the line
	for i := 0; i < len(lineBytes); i++ {
		for j := i + 1; j <= len(lineBytes); j++ {
			// We want to match a whole word, so use the null terminators to find word boundaries
			substr := append([]byte{0}, lineBytes[i:j]...)
			substr = append(substr, 0)

			// Check if the substring is a word
			indices := m.index.Lookup(substr, -1)
			if indices != nil {
				results = append(results, string(substr[1:len(substr)-1]))
			}
		}
	}

	return results
}
