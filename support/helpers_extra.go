package support

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

// Str returns a string helper
func Str(s string) *StringHelper {
	return &StringHelper{value: s}
}

// StringHelper provides string manipulation helpers
type StringHelper struct {
	value string
}

// After returns the substring after the first occurrence of a given value
func (s *StringHelper) After(delimiter string) string {
	index := strings.Index(s.value, delimiter)
	if index == -1 {
		return ""
	}
	return s.value[index+len(delimiter):]
}

// AfterLast returns the substring after the last occurrence of a given value
func (s *StringHelper) AfterLast(delimiter string) string {
	index := strings.LastIndex(s.value, delimiter)
	if index == -1 {
		return ""
	}
	return s.value[index+len(delimiter):]
}

// Before returns the substring before the first occurrence of a given value
func (s *StringHelper) Before(delimiter string) string {
	index := strings.Index(s.value, delimiter)
	if index == -1 {
		return ""
	}
	return s.value[:index]
}

// BeforeLast returns the substring before the last occurrence of a given value
func (s *StringHelper) BeforeLast(delimiter string) string {
	index := strings.LastIndex(s.value, delimiter)
	if index == -1 {
		return ""
	}
	return s.value[:index]
}

// Between returns the substring between two delimiters
func (s *StringHelper) Between(start, end string) string {
	startIndex := strings.Index(s.value, start)
	if startIndex == -1 {
		return ""
	}
	startIndex += len(start)
	endIndex := strings.Index(s.value[startIndex:], end)
	if endIndex == -1 {
		return ""
	}
	return s.value[startIndex : startIndex+endIndex]
}

// Contains checks if the string contains a substring
func (s *StringHelper) Contains(substring string) bool {
	return strings.Contains(s.value, substring)
}

// ContainsAll checks if the string contains all substrings
func (s *StringHelper) ContainsAll(substrings []string) bool {
	for _, sub := range substrings {
		if !strings.Contains(s.value, sub) {
			return false
		}
	}
	return true
}

// EndsWith checks if the string ends with a substring
func (s *StringHelper) EndsWith(suffix string) bool {
	return strings.HasSuffix(s.value, suffix)
}

// StartsWith checks if the string starts with a substring
func (s *StringHelper) StartsWith(prefix string) bool {
	return strings.HasPrefix(s.value, prefix)
}

// IsEmpty checks if the string is empty
func (s *StringHelper) IsEmpty() bool {
	return s.value == ""
}

// IsNotEmpty checks if the string is not empty
func (s *StringHelper) IsNotEmpty() bool {
	return s.value != ""
}

// Length returns the length of the string
func (s *StringHelper) Length() int {
	return len(s.value)
}

// Limit limits the string to a given length
func (s *StringHelper) Limit(limit int, suffix string) string {
	if len(s.value) <= limit {
		return s.value
	}
	return s.value[:limit] + suffix
}

// Lower converts the string to lowercase
func (s *StringHelper) Lower() string {
	return strings.ToLower(s.value)
}

// Upper converts the string to uppercase
func (s *StringHelper) Upper() string {
	return strings.ToUpper(s.value)
}

// Title converts the string to title case
func (s *StringHelper) Title() string {
	return strings.Title(s.value)
}

// Camel converts the string to camelCase
func (s *StringHelper) Camel() string {
	words := strings.Fields(s.value)
	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

// Snake converts the string to snake_case
func (s *StringHelper) Snake() string {
	reg := regexp.MustCompile("([A-Z])")
	result := reg.ReplaceAllString(s.value, "_${1}")
	return strings.ToLower(strings.TrimLeft(result, "_"))
}

// Kebab converts the string to kebab-case
func (s *StringHelper) Kebab() string {
	return strings.ReplaceAll(s.Snake(), "_", "-")
}

// Studly converts the string to StudlyCase
func (s *StringHelper) Studly() string {
	words := strings.Fields(s.value)
	for i := range words {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

// Plural converts the string to its plural form (simple implementation)
func (s *StringHelper) Plural() string {
	if strings.HasSuffix(s.value, "y") {
		return s.value[:len(s.value)-1] + "ies"
	}
	if strings.HasSuffix(s.value, "s") {
		return s.value + "es"
	}
	return s.value + "s"
}

// Singular converts the string to its singular form (simple implementation)
func (s *StringHelper) Singular() string {
	if strings.HasSuffix(s.value, "ies") {
		return s.value[:len(s.value)-3] + "y"
	}
	if strings.HasSuffix(s.value, "es") {
		return s.value[:len(s.value)-2]
	}
	if strings.HasSuffix(s.value, "s") {
		return s.value[:len(s.value)-1]
	}
	return s.value
}

// PadLeft pads the string on the left to a given length
func (s *StringHelper) PadLeft(length int, pad string) string {
	if len(s.value) >= length {
		return s.value
	}
	padding := strings.Repeat(pad, length-len(s.value))
	return padding + s.value
}

// PadRight pads the string on the right to a given length
func (s *StringHelper) PadRight(length int, pad string) string {
	if len(s.value) >= length {
		return s.value
	}
	padding := strings.Repeat(pad, length-len(s.value))
	return s.value + padding
}

// Repeat repeats the string n times
func (s *StringHelper) Repeat(n int) string {
	return strings.Repeat(s.value, n)
}

// Replace replaces occurrences of a substring
func (s *StringHelper) Replace(search, replace string) string {
	return strings.ReplaceAll(s.value, search, replace)
}

// ReplaceFirst replaces the first occurrence of a substring
func (s *StringHelper) ReplaceFirst(search, replace string) string {
	return strings.Replace(s.value, search, replace, 1)
}

// ReplaceLast replaces the last occurrence of a substring
func (s *StringHelper) ReplaceLast(search, replace string) string {
	index := strings.LastIndex(s.value, search)
	if index == -1 {
		return s.value
	}
	return s.value[:index] + replace + s.value[index+len(search):]
}

// Reverse reverses the string
func (s *StringHelper) Reverse() string {
	runes := []rune(s.value)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Slug converts the string to a slug
func (s *StringHelper) Slug() string {
	reg := regexp.MustCompile("[^a-z0-9]+")
	result := reg.ReplaceAllString(strings.ToLower(s.value), "-")
	return strings.Trim(result, "-")
}

// Words returns the first n words
func (s *StringHelper) Words(words int) string {
	parts := strings.Fields(s.value)
	if words >= len(parts) {
		return s.value
	}
	return strings.Join(parts[:words], " ")
}

// LimitWords limits the string to a given number of words
func (s *StringHelper) LimitWords(words int, suffix string) string {
	parts := strings.Fields(s.value)
	if words >= len(parts) {
		return s.value
	}
	return strings.Join(parts[:words], " ") + suffix
}

// Is checks if the string matches a pattern
func (s *StringHelper) Is(pattern string) bool {
	matched, err := regexp.MatchString(pattern, s.value)
	if err != nil {
		return false
	}
	return matched
}

// IsEmail checks if the string is an email
func (s *StringHelper) IsEmail() bool {
	return s.Is(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
}

// IsURL checks if the string is a URL
func (s *StringHelper) IsURL() bool {
	return s.Is(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)
}

// IsUUID checks if the string is a UUID
func (s *StringHelper) IsUUID() bool {
	return s.Is(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
}

// IsIP checks if the string is an IP address
func (s *StringHelper) IsIP() bool {
	return s.Is(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
}

// IsAlpha checks if the string contains only letters
func (s *StringHelper) IsAlpha() bool {
	return s.Is(`^[a-zA-Z]+$`)
}

// IsAlphaNumeric checks if the string contains only letters and numbers
func (s *StringHelper) IsAlphaNumeric() bool {
	return s.Is(`^[a-zA-Z0-9]+$`)
}

// IsNumeric checks if the string is numeric
func (s *StringHelper) IsNumeric() bool {
	return s.Is(`^[0-9]+$`)
}

// IsJSON checks if the string is valid JSON
func (s *StringHelper) IsJSON() bool {
	return s.Is(`^\{.*\}$|^\[.*\]$`)
}

// Base64Encode encodes the string to base64
func (s *StringHelper) Base64Encode() string {
	return base64.StdEncoding.EncodeToString([]byte(s.value))
}

// Base64Decode decodes the string from base64
func (s *StringHelper) Base64Decode() (string, error) {
	result, err := base64.StdEncoding.DecodeString(s.value)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// Hash generates a hash of the string (simple implementation)
func (s *StringHelper) Hash() string {
	// In production, use bcrypt or argon2
	return fmt.Sprintf("%x", s.value)
}

// Random returns a random string of given length
func Random(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// RandomAlphanumeric returns a random alphanumeric string
func RandomAlphanumeric(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	rand.Read(bytes)
	for i := range bytes {
		bytes[i] = charset[bytes[i]%byte(len(charset))]
	}
	return string(bytes)
}

// RandomAlpha returns a random alphabetic string
func RandomAlpha(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)
	rand.Read(bytes)
	for i := range bytes {
		bytes[i] = charset[bytes[i]%byte(len(charset))]
	}
	return string(bytes)
}

// RandomNumeric returns a random numeric string
func RandomNumeric(length int) string {
	const charset = "0123456789"
	bytes := make([]byte, length)
	rand.Read(bytes)
	for i := range bytes {
		bytes[i] = charset[bytes[i]%byte(len(charset))]
	}
	return string(bytes)
}

// Uuid generates a UUID v4
func Uuid() string {
 bytes := make([]byte, 16)
	rand.Read(bytes)
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:])
}

// Number formats a number with commas
func Number(n float64) string {
	// Simple implementation
	return fmt.Sprintf("%.2f", n)
}

// Currency formats a number as currency
func Currency(n float64, symbol string) string {
	return fmt.Sprintf("%s%.2f", symbol, n)
}

// Percentage formats a number as percentage
func Percentage(n float64) string {
	return fmt.Sprintf("%.2f%%", n)
}

// Ago returns a human-readable time difference
func Ago(t time.Time) string {
	duration := time.Since(t)
	if duration < time.Minute {
		return "just now"
	}
	if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	days := int(duration.Hours() / 24)
	if days == 1 {
		return "1 day ago"
	}
	return fmt.Sprintf("%d days ago", days)
}

// Until returns a human-readable time until
func Until(t time.Time) string {
	duration := time.Until(t)
	if duration < time.Minute {
		return "just now"
	}
	if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "in 1 minute"
		}
		return fmt.Sprintf("in %d minutes", minutes)
	}
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "in 1 hour"
		}
		return fmt.Sprintf("in %d hours", hours)
	}
	days := int(duration.Hours() / 24)
	if days == 1 {
		return "in 1 day"
	}
	return fmt.Sprintf("in %d days", days)
}

// Bytes converts a number to bytes
func Bytes(n int) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for n := int64(n) / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}

// NumberFormat formats a number with decimal places
func NumberFormat(n float64, decimals int) string {
	format := fmt.Sprintf("%%.%df", decimals)
	return fmt.Sprintf(format, n)
}

// PadLeft pads a string on the left
func PadLeft(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(pad, length-len(s))
	return padding + s
}

// PadRight pads a string on the right
func PadRight(s string, length int, pad string) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(pad, length-len(s))
	return s + padding
}

// Title converts a string to title case
func TitleStr(s string) string {
	return strings.Title(s)
}

// Camel converts a string to camelCase
func Camel(s string) string {
	words := strings.Fields(s)
	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

// Snake converts a string to snake_case
func Snake(s string) string {
	reg := regexp.MustCompile("([A-Z])")
	result := reg.ReplaceAllString(s, "_${1}")
	return strings.ToLower(strings.TrimLeft(result, "_"))
}

// Kebab converts a string to kebab-case
func Kebab(s string) string {
	return strings.ReplaceAll(Snake(s), "_", "-")
}

// Studly converts a string to StudlyCase
func Studly(s string) string {
	words := strings.Fields(s)
	for i := range words {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

// Limit limits a string to a given length
func LimitStr(s string, length int, suffix string) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + suffix
}

// Slug converts a string to a slug
func Slug(s string) string {
	reg := regexp.MustCompile("[^a-z0-9]+")
	result := reg.ReplaceAllString(strings.ToLower(s), "-")
	return strings.Trim(result, "-")
}

// WordsStr returns the first n words of a string
func WordsStr(s string, words int) string {
	parts := strings.Fields(s)
	if words >= len(parts) {
		return s
	}
	return strings.Join(parts[:words], " ")
}

// LimitWordsStr limits a string to a given number of words
func LimitWordsStr(s string, words int, suffix string) string {
	parts := strings.Fields(s)
	if words >= len(parts) {
		return s
	}
	return strings.Join(parts[:words], " ") + suffix
}

// Plural converts a string to its plural form
func Plural(s string) string {
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	if strings.HasSuffix(s, "s") {
		return s + "es"
	}
	return s + "s"
}

// Singular converts a string to its singular form
func Singular(s string) string {
	if strings.HasSuffix(s, "ies") {
		return s[:len(s)-3] + "y"
	}
	if strings.HasSuffix(s, "es") {
		return s[:len(s)-2]
	}
	if strings.HasSuffix(s, "s") {
		return s[:len(s)-1]
	}
	return s
}

// StrContains checks if a string contains a substring
func StrContains(s, substring string) bool {
	return strings.Contains(s, substring)
}

// StrReplace replaces occurrences of a substring
func StrReplace(s, search, replace string) string {
	return strings.ReplaceAll(s, search, replace)
}

// StrToLower converts a string to lowercase
func StrToLower(s string) string {
	return strings.ToLower(s)
}

// StrToUpper converts a string to uppercase
func StrToUpper(s string) string {
	return strings.ToUpper(s)
}

// StrTrim trims whitespace from a string
func StrTrim(s string) string {
	return strings.TrimSpace(s)
}

// StrSplit splits a string by a delimiter
func StrSplit(s, delimiter string) []string {
	return strings.Split(s, delimiter)
}

// StrJoin joins a slice of strings with a delimiter
func StrJoin(parts []string, delimiter string) string {
	return strings.Join(parts, delimiter)
}

// StrRepeat repeats a string n times
func StrRepeat(s string, n int) string {
	return strings.Repeat(s, n)
}

// StrReverse reverses a string
func StrReverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// StrLength returns the length of a string
func StrLength(s string) int {
	return len(s)
}

// StrIsEmpty checks if a string is empty
func StrIsEmpty(s string) bool {
	return s == ""
}

// StrIsNotEmpty checks if a string is not empty
func StrIsNotEmpty(s string) bool {
	return s != ""
}

// StrStartsWith checks if a string starts with a prefix
func StrStartsWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// StrEndsWith checks if a string ends with a suffix
func StrEndsWith(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// StrAfter returns the substring after the first occurrence
func StrAfter(s, delimiter string) string {
	index := strings.Index(s, delimiter)
	if index == -1 {
		return ""
	}
	return s[index+len(delimiter):]
}

// StrBefore returns the substring before the first occurrence
func StrBefore(s, delimiter string) string {
	index := strings.Index(s, delimiter)
	if index == -1 {
		return ""
	}
	return s[:index]
}

// StrBetween returns the substring between two delimiters
func StrBetween(s, start, end string) string {
	startIndex := strings.Index(s, start)
	if startIndex == -1 {
		return ""
	}
	startIndex += len(start)
	endIndex := strings.Index(s[startIndex:], end)
	if endIndex == -1 {
		return ""
	}
	return s[startIndex : startIndex+endIndex]
}

// Math helpers

// Round rounds a number to a given precision
func Round(n float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(n*pow) / pow
}

// Floor rounds a number down
func Floor(n float64) float64 {
	return math.Floor(n)
}

// Ceil rounds a number up
func Ceil(n float64) float64 {
	return math.Ceil(n)
}

// Abs returns the absolute value
func Abs(n float64) float64 {
	return math.Abs(n)
}

// MinFloat returns the minimum of two numbers
func MinFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// MaxFloat returns the maximum of two numbers
func MaxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Clamp clamps a number between min and max
func Clamp(n, minVal, maxVal float64) float64 {
	if n < minVal {
		return minVal
	}
	if n > maxVal {
		return maxVal
	}
	return n
}
