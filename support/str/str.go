package str

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"unicode"
)

// Camel converts a string to camelCase.
func Camel(s string) string {
	words := splitIntoWords(s)
	if len(words) == 0 {
		return ""
	}
	
	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		result += title(words[i])
	}
	return result
}

// Studly converts a string to StudlyCase.
func Studly(s string) string {
	words := splitIntoWords(s)
	result := ""
	for _, word := range words {
		result += title(word)
	}
	return result
}

// Snake converts a string to snake_case.
func Snake(s string) string {
	words := splitIntoWords(s)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, "_")
}

// Kebab converts a string to kebab-case.
func Kebab(s string) string {
	words := splitIntoWords(s)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, "-")
}

// Random generates a random alphanumeric string.
func Random(length int) string {
	bytes := make([]byte, length/2+1)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)[:length]
}

// Limit truncates a string to a given length.
func Limit(s string, length int, end string) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + end
}

// Contains determines if a string contains another string.
func Contains(haystack string, needle string) bool {
	return strings.Contains(haystack, needle)
}

// splitIntoWords splits a string into words based on common delimiters and case changes.
func splitIntoWords(s string) []string {
	var words []string
	var currentWord strings.Builder

	for i, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			continue
		}

		if unicode.IsUpper(r) && currentWord.Len() > 0 {
			// Check if previous character was lower or next is lower (e.g., XMLHttp -> XML Http)
			prev := rune(s[i-1])
			if unicode.IsLower(prev) {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		}

		currentWord.WriteRune(r)
	}

	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}

func title(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

