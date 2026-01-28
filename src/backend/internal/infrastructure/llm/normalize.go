package llm

import (
	"regexp"
	"strings"
	"unicode"
)

// dateRegex matches a valid ISO date format YYYY-MM-DD
var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// NormalizeSpacedText removes spurious spaces within words that are artifacts
// from OCR or PDF text extraction. For example, "Co lumb ia" becomes "Columbia".
//
// The algorithm looks for patterns where text fragments appear to be split
// and merges them intelligently while preserving legitimate word boundaries.
func NormalizeSpacedText(text string) string {
	if text == "" {
		return ""
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	// Analyze the pattern
	singleCharCount := 0
	numericFragments := 0
	for _, word := range words {
		if len(word) == 1 {
			singleCharCount++
		}
		// Check for numeric fragments like "201" in "201 4"
		if isNumericFragment(word) {
			numericFragments++
		}
	}

	// If many single chars (> 50%), it's fully spaced text like "C o l u m b i a"
	if len(words) > 2 && float64(singleCharCount)/float64(len(words)) > 0.5 {
		return mergeFullySpacedText(words)
	}

	// Otherwise, handle partial spacing artifacts
	return mergePartiallySpacedWords(words)
}

// isNumericFragment checks if a word is a partial number (contains only digits)
func isNumericFragment(word string) bool {
	if len(word) == 0 {
		return false
	}
	for _, r := range word {
		if !unicode.IsDigit(r) && r != '-' {
			return false
		}
	}
	return true
}

// mergeFullySpacedText handles text where most characters are separated by spaces
func mergeFullySpacedText(words []string) string {
	var result strings.Builder
	for _, word := range words {
		if len(word) == 1 {
			result.WriteString(word)
		} else {
			if result.Len() > 0 {
				result.WriteString(" ")
			}
			result.WriteString(word)
		}
	}
	return result.String()
}

// mergePartiallySpacedWords handles cases where words are partially split
// like "Co lumb ia" should become "Columbia"
func mergePartiallySpacedWords(words []string) string {
	if len(words) == 0 {
		return ""
	}
	if len(words) == 1 {
		return words[0]
	}

	result := make([]string, 0, len(words))
	i := 0

	for i < len(words) {
		current := words[i]

		// Look ahead to see if we should merge with following words
		merged := current
		j := i + 1

		for j < len(words) {
			nextWord := words[j]

			// Check if these look like they should be merged
			if shouldMerge(merged, nextWord) {
				merged += nextWord
				j++
			} else {
				break
			}
		}

		result = append(result, merged)
		i = j
	}

	return strings.Join(result, " ")
}

// shouldMerge determines if two adjacent "words" should be merged
func shouldMerge(current, next string) bool {
	// Never merge common short words (prepositions, articles)
	if isCommonShortWord(next) || isCommonShortWord(current) {
		return false
	}

	// Both are numeric or date-like fragments - merge them
	if isNumericOrDateFragment(current) && isNumericOrDateFragment(next) {
		return true
	}

	// If next word starts with lowercase, it's likely a continuation
	if startsWithLowercase(next) {
		return true
	}

	// Check for short fragment patterns
	return isShortFragmentPattern(current, next)
}

// startsWithLowercase checks if a string starts with a lowercase letter
func startsWithLowercase(s string) bool {
	if len(s) == 0 {
		return false
	}
	runes := []rune(s)
	return unicode.IsLower(runes[0])
}

// isShortFragmentPattern checks for patterns that indicate word fragments
func isShortFragmentPattern(current, next string) bool {
	currentLen := len(current)
	nextLen := len(next)

	// If current word is very short (1-2 chars) and next is also short fragment
	if currentLen <= 2 && nextLen <= 4 && startsWithLowercase(next) {
		return true
	}

	return false
}


// isNumericOrDateFragment checks if a string is a numeric fragment or date part
func isNumericOrDateFragment(s string) bool {
	if len(s) == 0 {
		return false
	}
	hasDigit := false
	for _, r := range s {
		if unicode.IsDigit(r) {
			hasDigit = true
		} else if r != '-' {
			return false
		}
	}
	return hasDigit
}

// isCommonShortWord checks if a word is a common short English word
// that should NOT be merged with adjacent words
func isCommonShortWord(word string) bool {
	lower := strings.ToLower(word)
	commonWords := map[string]bool{
		"a": true, "an": true, "as": true, "at": true, "be": true,
		"by": true, "do": true, "go": true, "he": true, "if": true,
		"in": true, "is": true, "it": true, "me": true, "my": true,
		"no": true, "of": true, "on": true, "or": true, "so": true,
		"to": true, "up": true, "us": true, "we": true,
		"the": true, "and": true, "for": true, "are": true, "but": true,
		"not": true, "you": true, "all": true, "can": true, "her": true,
		"was": true, "one": true, "our": true, "out": true,
	}
	return commonWords[lower]
}

// NormalizeDate cleans up a date string and validates it's in ISO format.
// Returns empty string if the date is malformed.
func NormalizeDate(date string) string {
	if date == "" || date == "null" {
		return ""
	}

	// First, remove all spaces (handles "201 4-01-0 1" -> "2014-01-01")
	cleaned := strings.ReplaceAll(date, " ", "")

	// Check if it starts with a dash (invalid like "-01-01")
	if strings.HasPrefix(cleaned, "-") {
		return ""
	}

	// Extract just the date part if there's extra content
	// Match YYYY-MM-DD pattern
	if dateRegex.MatchString(cleaned) {
		return cleaned
	}

	// Try to extract date from string with extra content
	// Find first occurrence of date pattern
	matches := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`).FindString(cleaned)
	if matches != "" && !strings.HasPrefix(matches, "0000") {
		return matches
	}

	return ""
}

// HasExcessiveSpacing checks if text has OCR-style spacing artifacts
// where single characters are separated by spaces.
func HasExcessiveSpacing(text string) bool {
	if text == "" {
		return false
	}

	words := strings.Fields(text)
	if len(words) <= 1 {
		return false
	}

	// Count single-character "words" and short fragments
	singleCharCount := 0
	shortFragmentCount := 0 // 2-3 char non-words

	for _, word := range words {
		if len(word) == 1 {
			singleCharCount++
		} else if len(word) <= 3 && !isCommonShortWord(word) {
			shortFragmentCount++
		}
	}

	// If we have multiple single chars, it's likely spacing artifacts
	if singleCharCount >= 2 {
		return true
	}

	// If we have multiple short fragments that aren't common words
	if shortFragmentCount >= 2 {
		return true
	}

	return false
}
