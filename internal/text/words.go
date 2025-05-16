package text

import (
	"math/rand"
	"time"
)

// Common English words between 3-8 characters
var allWords = []string{
	// 3 letters
	"the", "and", "for", "are", "but", "not", "you", "all", "any", "can", "had", "her", "was", "one", "our", "out", "day", "get", "has", "him",

	// 4 letters
	"that", "with", "have", "this", "will", "your", "from", "they", "know", "want", "been", "good", "much", "some", "time", "very", "when", "come", "just", "into",

	// 5 letters
	"about", "there", "think", "would", "could", "people", "other", "first", "their", "these", "words", "which", "water", "write", "place", "sound", "great", "where", "help", "through",

	// 6 letters
	"should", "because", "through", "before", "little", "change", "around", "always", "between", "system", "number", "family", "second", "enough", "moment", "though", "person", "better", "really", "almost",

	// 7 letters
	"another", "thought", "example", "picture", "science", "measure", "product", "history", "position", "company", "quality", "service", "support", "network", "project", "process", "control", "current", "program", "problem",

	// 8 letters
	"language", "computer", "business", "research", "industry", "security", "software", "hardware", "database", "internet", "analysis", "solution", "strategy", "resource", "practice", "evidence", "approach", "function", "complete", "standard",
}

// GetRandomWords returns a random selection of words for the typing test
func GetRandomWords(count int) []string {
	// Create a local random generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create a copy of allWords to avoid modifying the original
	words := make([]string, len(allWords))
	copy(words, allWords)

	// Shuffle the words using the local generator
	r.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})

	// Return the first 'count' words
	if count > len(words) {
		count = len(words)
	}
	return words[:count]
}

// GetWords returns a fixed set of words for testing
func GetWords() []string {
	return GetRandomWords(50) // Default to 50 words per test
}
