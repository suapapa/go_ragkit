// Package ragkit provides utility functions for RAG (Retrieval-Augmented Generation) applications.
package ragkit

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// MakeDocsFromTexts creates a slice of Document from a slice of texts with optional metadata.
// Each document will have a unique ID generated from its text content.
func MakeDocsFromTexts(texts []string, metadata map[string]any) []Document {
	docs := make([]Document, len(texts))
	for i, text := range texts {
		docs[i] = Document{ID: GenerateID(text), Text: text, Metadata: metadata}
	}
	return docs
}

// GenerateID creates a UUID-like string from the input text using SHA-256 hash.
// The generated ID follows the format: 8-4-4-4-12 hexadecimal characters.
func GenerateID(text string) string {
	hash := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", hash[0:4], hash[4:6], hash[6:8], hash[8:10], hash[10:16])
}

func ToCamelCase(input string) string {
	words := strings.FieldsFunc(input, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	for i, word := range words {
		if len(word) == 0 {
			continue
		}
		words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
	}

	return strings.Join(words, "")
}
