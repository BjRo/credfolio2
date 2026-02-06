package llm

import (
	"os"
	"testing"
)

func TestExtractTextFromPDF_TextBasedPDF(t *testing.T) {
	// Use the fixture resume PDF which is a text-based PDF
	data, err := os.ReadFile("../../../../../fixtures/CV_TEMPLATE_0004.pdf")
	if err != nil {
		t.Skipf("Fixture PDF not available: %v", err)
	}

	text, err := extractTextFromPDF(data)
	if err != nil {
		t.Fatalf("extractTextFromPDF() error = %v", err)
	}

	if text == "" {
		t.Fatal("extractTextFromPDF() returned empty text")
	}

	// A real resume should have reasonable content
	if len(text) < 100 {
		t.Errorf("extractTextFromPDF() returned suspiciously short text (%d chars): %q", len(text), text)
	}
}

func TestExtractTextFromPDF_InvalidData(t *testing.T) {
	// Random bytes that aren't a valid PDF
	_, err := extractTextFromPDF([]byte{0x01, 0x02, 0x03, 0x04})
	if err == nil {
		t.Fatal("expected error for invalid PDF data")
	}
}

func TestExtractTextFromPDF_EmptyData(t *testing.T) {
	_, err := extractTextFromPDF([]byte{})
	if err == nil {
		t.Fatal("expected error for empty data")
	}
}

func TestIsUsableText_GoodText(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{
			name: "normal resume text",
			text: "John Doe\nSoftware Engineer\nExperience: 5 years of building web applications using Go and React.\nEducation: BS Computer Science, MIT",
			want: true,
		},
		{
			name: "text with reasonable word density",
			text: "This is a well-formed document with proper English text that contains meaningful content about work experience and professional skills.",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUsableText(tt.text); got != tt.want {
				t.Errorf("isUsableText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUsableText_BadText(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{
			name: "empty string",
			text: "",
			want: false,
		},
		{
			name: "only whitespace",
			text: "   \n\t\n   ",
			want: false,
		},
		{
			name: "too short",
			text: "Hello",
			want: false,
		},
		{
			name: "garbage characters",
			text: "ÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿÿ",
			want: false,
		},
		{
			name: "mostly non-ASCII",
			text: "ñ€∞§¶•ªºÀÁÂÃÄÅÆÇÈÉÊ ñ€∞§¶•ªºÀÁÂÃÄÅÆÇÈÉÊ ñ€∞§¶•ªºÀÁÂÃÄÅÆÇÈÉÊ ñ€∞§¶•ªºÀÁÂÃÄÅÆÇÈÉÊ the end",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUsableText(tt.text); got != tt.want {
				t.Errorf("isUsableText() = %v, want %v for text: %q", got, tt.want, tt.text)
			}
		})
	}
}
