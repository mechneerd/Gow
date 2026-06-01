package str

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"strings"
	"time"
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

func splitIntoWords(s string) []string {
	var words []string
	var currentWord strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			continue
		}
		if unicode.IsUpper(r) && currentWord.Len() > 0 {
			if i > 0 {
				prev := runes[i-1]
				if unicode.IsLower(prev) {
					words = append(words, currentWord.String())
					currentWord.Reset()
				}
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

var timeRand = time.Now().UnixNano()

// After returns the substring after the first occurrence of a value.
func After(s, search string) string {
	idx := strings.Index(s, search)
	if idx == -1 {
		return ""
	}
	return s[idx+len(search):]
}

// AfterLast returns the substring after the last occurrence of a value.
func AfterLast(s, search string) string {
	idx := strings.LastIndex(s, search)
	if idx == -1 {
		return ""
	}
	return s[idx+len(search):]
}

// Before returns the substring before the first occurrence of a value.
func Before(s, search string) string {
	idx := strings.Index(s, search)
	if idx == -1 {
		return ""
	}
	return s[:idx]
}

// BeforeLast returns the substring before the last occurrence of a value.
func BeforeLast(s, search string) string {
	idx := strings.LastIndex(s, search)
	if idx == -1 {
		return ""
	}
	return s[:idx]
}

// Between returns the substring between two values.
func Between(s, from, to string) string {
	start := strings.Index(s, from)
	if start == -1 {
		return ""
	}
	start += len(from)
	end := strings.Index(s[start:], to)
	if end == -1 {
		return ""
	}
	return s[start : start+end]
}

// Finish appends a suffix to the string if it doesn't already end with it.
func Finish(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s
	}
	return s + suffix
}

// IsEmpty determines if the string is empty.
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsNotEmpty determines if the string is not empty.
func IsNotEmpty(s string) bool {
	return len(s) > 0
}

// IsBlank determines if the string is blank (empty or whitespace).
func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// IsNotBlank determines if the string is not blank.
func IsNotBlank(s string) bool {
	return !IsBlank(s)
}

// StartsWith determines if the string starts with a given prefix.
func StartsWith(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

// EndsWith determines if the string ends with a given suffix.
func EndsWith(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

// Plural appends "s" to the string if count is not 1.
func Plural(s string, count int) string {
	if count == 1 {
		return s
	}
	return s + "s"
}

// PluralStudly converts the string to StudlyCase and pluralizes it.
func PluralStudly(s string, count int) string {
	return Plural(Studly(s), count)
}

// PluralSnake converts the string to snake_case and pluralizes it.
func PluralSnake(s string, count int) string {
	return Plural(Snake(s), count)
}

// Title converts the string to title case.
func Title(s string) string {
	return strings.Title(s)
}

// Ucfirst capitalizes the first character of the string.
func Ucfirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// Lcfirst lowercases the first character of the string.
func Lcfirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// Lower converts the string to lowercase.
func Lower(s string) string {
	return strings.ToLower(s)
}

// Upper converts the string to uppercase.
func Upper(s string) string {
	return strings.ToUpper(s)
}

// Replace replaces all occurrences of a search with a replace.
func Replace(search, replace, subject string) string {
	return strings.ReplaceAll(subject, search, replace)
}

// ReplaceFirst replaces the first occurrence of a search with a replace.
func ReplaceFirst(search, replace, subject string) string {
	return strings.Replace(subject, search, replace, 1)
}

// ReplaceLast replaces the last occurrence of a search with a replace.
func ReplaceLast(search, replace, subject string) string {
	idx := strings.LastIndex(subject, search)
	if idx == -1 {
		return subject
	}
	return subject[:idx] + replace + subject[idx+len(search):]
}

// Mask masks a portion of the string with a character.
func Mask(s string, character string, keepStart, keepEnd int) string {
	runes := []rune(s)
	totalLen := len(runes)
	if keepStart+keepEnd >= totalLen {
		return s
	}

	maskLen := totalLen - keepStart - keepEnd
	result := make([]rune, 0, totalLen)
	result = append(result, runes[:keepStart]...)
	for i := 0; i < maskLen; i++ {
		result = append(result, []rune(character)...)
	}
	result = append(result, runes[totalLen-keepEnd:]...)
	return string(result)
}

// Words limits the string to a given number of words.
func Words(s string, limit int, end string) string {
	words := strings.Fields(s)
	if len(words) <= limit {
		return s
	}
	return strings.Join(words[:limit], " ") + end
}

// StrContains checks if the string contains the given value.
func StrContains(haystack string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}

// StrSlug generates a URL-friendly slug from the string.
func StrSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' || r == '_' {
			return r
		}
		if unicode.IsSpace(r) {
			return '-'
		}
		return -1
	}, s)
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}

// StrLimit truncates the string to the given length, preserving whole words.
func StrLimit(s string, length int, end string) string {
	if len(s) <= length {
		return s
	}
	truncated := s[:length]
	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}
	return truncated + end
}

// Random generates a random string of the given length.
func Random(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}

// Uuid generates a UUID v4 string.
func Uuid() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// Is checks if the string matches a given pattern.
func Is(pattern, value string) bool {
	matched, _ := filepath.Match(pattern, value)
	return matched
}

// Start appends a prefix to the string if it doesn't already start with it.
func Start(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s
	}
	return prefix + s
}

// Repeat repeats the string n times.
func Repeat(s string, n int) string {
	return strings.Repeat(s, n)
}

// Reverse reverses the string.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Trim removes characters from both ends of the string.
func Trim(s string, characters ...string) string {
	if len(characters) == 0 {
		return strings.TrimSpace(s)
	}
	return strings.Trim(s, characters[0])
}

// Ltrim removes characters from the left side of the string.
func Ltrim(s string, characters ...string) string {
	if len(characters) == 0 {
		return strings.TrimLeftFunc(s, unicode.IsSpace)
	}
	return strings.TrimLeft(s, characters[0])
}

// Rtrim removes characters from the right side of the string.
func Rtrim(s string, characters ...string) string {
	if len(characters) == 0 {
		return strings.TrimRightFunc(s, unicode.IsSpace)
	}
	return strings.TrimRight(s, characters[0])
}

// Lpad pads the left side of the string with another string until it reaches the given length.
func Lpad(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(pad, (length-len(s))/len(pad)+1)
	return padding[:length-len(s)] + s
}

// Rpad pads the right side of the string with another string until it reaches the given length.
func Rpad(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(pad, (length-len(s))/len(pad)+1)
	return s + padding[:length-len(s)]
}

// SwapCase swaps the case of each character.
func SwapCase(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsUpper(r) {
			return unicode.ToLower(r)
		}
		return unicode.ToUpper(r)
	}, s)
}

// IsAlpha checks if the string contains only letters.
func IsAlpha(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return len(s) > 0
}

// IsAlphaNumeric checks if the string contains only letters and numbers.
func IsAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return len(s) > 0
}

// IsNumeric checks if the string contains only numbers.
func IsNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsNumber(r) {
			return false
		}
	}
	return len(s) > 0
}

// IsEmail checks if the string is a valid email address.
func IsEmail(s string) bool {
	parts := strings.Split(s, "@")
	if len(parts) != 2 {
		return false
	}
	return len(parts[0]) > 0 && strings.Contains(parts[1], ".")
}

// IsURL checks if the string is a valid URL.
func IsURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

// IsUUID checks if the string is a valid UUID.
func IsUUID(s string) bool {
	return len(s) == 36 && strings.Count(s, "-") == 4
}

// IsIP checks if the string is a valid IP address.
func IsIP(s string) bool {
	return net.ParseIP(s) != nil
}

// IsIPv4 checks if the string is a valid IPv4 address.
func IsIPv4(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil && ip.To4() != nil
}

// IsIPv6 checks if the string is a valid IPv6 address.
func IsIPv6(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil && ip.To4() == nil
}

// IsJSON checks if the string is valid JSON.
func IsJSON(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}

// AfterFunc returns the substring after the first occurrence of a value using a function.
func AfterFunc(s string, f func(rune) bool) string {
	for i, r := range s {
		if f(r) {
			return s[i+len(string(r)):]
		}
	}
	return ""
}

// BeforeFunc returns the substring before the first occurrence of a value using a function.
func BeforeFunc(s string, f func(rune) bool) string {
	for i, r := range s {
		if f(r) {
			return s[:i]
		}
	}
	return s
}

// ContainsAll determines if the string contains all substrings.
func ContainsAll(haystack string, needles []string) bool {
	for _, needle := range needles {
		if !strings.Contains(haystack, needle) {
			return false
		}
	}
	return true
}

// ContainsAny determines if the string contains any of the substrings.
func ContainsAny(haystack string, needles []string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}
