package llm

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
)

// minUsableTextLength is the minimum number of characters for extracted text to be considered usable.
// Anything shorter is likely a scanned document or extraction failure.
const minUsableTextLength = 50

// minASCIIWordRatio is the minimum fraction of whitespace-separated tokens that must consist
// of ASCII-only characters. Scanned/garbled PDFs produce many non-ASCII tokens.
const minASCIIWordRatio = 0.5

// extractTextFromPDF extracts plain text from a PDF using a Go-native library.
// Returns the extracted text or an error if the PDF cannot be parsed.
func extractTextFromPDF(data []byte) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("empty PDF data")
	}

	reader := bytes.NewReader(data)
	pdfReader, err := pdf.NewReader(reader, int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("failed to parse PDF: %w", err)
	}

	var buf strings.Builder
	numPages := pdfReader.NumPage()

	for i := 1; i <= numPages; i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue // Skip pages that fail to extract; others may still work.
		}

		if buf.Len() > 0 && text != "" {
			buf.WriteString("\n")
		}
		buf.WriteString(text)
	}

	return strings.TrimSpace(buf.String()), nil
}

// isUsableText checks whether locally-extracted text is good enough to skip the LLM.
// It verifies length, UTF-8 validity, and that a reasonable fraction of tokens look like
// normal (mostly-ASCII) words rather than garbled binary output.
func isUsableText(text string) bool {
	trimmed := strings.TrimSpace(text)

	if len(trimmed) < minUsableTextLength {
		return false
	}

	if !utf8.ValidString(trimmed) {
		return false
	}

	// Check that a sufficient fraction of words are ASCII-like.
	words := strings.Fields(trimmed)
	if len(words) == 0 {
		return false
	}

	asciiWords := 0
	for _, w := range words {
		if isASCIIWord(w) {
			asciiWords++
		}
	}

	return float64(asciiWords)/float64(len(words)) >= minASCIIWordRatio
}

// isASCIIWord returns true if the word consists entirely of ASCII printable characters.
func isASCIIWord(w string) bool {
	for _, r := range w {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
