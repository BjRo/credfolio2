package resolver

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// calculateContentHash computes the SHA-256 hash of the provided content.
// Returns the hash as a lowercase hexadecimal string.
func calculateContentHash(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

// calculateContentHashFromReader computes the SHA-256 hash of content from a reader.
// Returns the hash as a lowercase hexadecimal string and the content as bytes.
// This allows the content to be used for upload after hash calculation.
func calculateContentHashFromReader(r io.Reader) (string, []byte, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return "", nil, err
	}
	return calculateContentHash(content), content, nil
}
