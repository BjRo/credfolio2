//nolint:errcheck,revive // Test file
package llm

import (
	"testing"
)

func TestNormalizeSpacedText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "normal text unchanged",
			input: "Columbia University",
			want:  "Columbia University",
		},
		{
			name:  "spaced out letters",
			input: "C o l u m b i a University",
			want:  "Columbia University",
		},
		{
			name:  "partially spaced",
			input: "Co lumb ia University",
			want:  "Columbia University",
		},
		{
			name:  "spaced numbers in date",
			input: "2 0 1 8",
			want:  "2018",
		},
		{
			name:  "date with internal spaces",
			input: "201 4-01-0 1",
			want:  "2014-01-01",
		},
		{
			name:  "word with single space artifact",
			input: "Comput er Information Systems",
			want:  "Computer Information Systems",
		},
		{
			name:  "multiple words with artifacts",
			input: "Co lumb ia Univ ersity",
			want:  "Columbia University",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "preserves legitimate spaces",
			input: "Bachelor of Science",
			want:  "Bachelor of Science",
		},
		{
			name:  "mixed case preserved",
			input: "Ne w York",
			want:  "New York",
		},
		{
			name:  "handles consecutive single chars",
			input: "a b c d e f",
			want:  "abcdef",
		},
		{
			name:  "preserves proper multi-word phrases",
			input: "Computer Information Systems",
			want:  "Computer Information Systems",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeSpacedText(tt.input)
			if got != tt.want {
				t.Errorf("NormalizeSpacedText(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeDate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "valid date unchanged",
			input: "2018-01-01",
			want:  "2018-01-01",
		},
		{
			name:  "date with spaces",
			input: "201 4-01-0 1",
			want:  "2014-01-01",
		},
		{
			name:  "fully spaced date",
			input: "2 0 1 8 - 0 1 - 0 1",
			want:  "2018-01-01",
		},
		{
			name:  "date with leading dash returns empty",
			input: "-01-01",
			want:  "",
		},
		{
			name:  "malformed date returns empty",
			input: "invalid",
			want:  "",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "null string",
			input: "null",
			want:  "",
		},
		{
			name:  "date with extra content stripped",
			input: "2018-09-01 some extra text",
			want:  "2018-09-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeDate(tt.input)
			if got != tt.want {
				t.Errorf("NormalizeDate(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestHasExcessiveSpacing(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "normal text",
			input: "Columbia University",
			want:  false,
		},
		{
			name:  "spaced out letters",
			input: "C o l u m b i a",
			want:  true,
		},
		{
			name:  "single space artifact in word",
			input: "Co lumb ia",
			want:  true,
		},
		{
			name:  "multiple proper words",
			input: "Bachelor of Science in Computer Information Systems",
			want:  false,
		},
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
		{
			name:  "single word",
			input: "Columbia",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasExcessiveSpacing(tt.input)
			if got != tt.want {
				t.Errorf("HasExcessiveSpacing(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
