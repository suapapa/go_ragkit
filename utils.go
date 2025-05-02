// Package ragkit provides utility functions for RAG (Retrieval-Augmented Generation) applications.
package ragkit

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
)

// MakeDocsFromTexts creates a slice of Document from a slice of texts with optional metadata.
// Each document will have a unique ID generated from its text content.
func MakeDocsFromTexts(texts []string, metadata map[string]any) []Document {
	docs := make([]Document, len(texts))
	for i, text := range texts {
		docs[i] = Document{ID: GenerateID(text, metadata), Text: text, Metadata: metadata}
	}
	return docs
}

// GenerateID creates a UUID-like string from the input text using SHA-256 hash.
// The generated ID follows the format: 8-4-4-4-12 hexadecimal characters.
func GenerateID(text string, metadata map[string]any) string {
	hash := sha256.New()

	var metadataB []byte
	if metadata != nil {
		metadataB, _ = json.Marshal(metadata)
		hash.Write(metadataB)
	}

	hash.Write([]byte(text))
	hashBytes := hash.Sum(nil)

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", hashBytes[0:4], hashBytes[4:6], hashBytes[6:8], hashBytes[8:10], hashBytes[10:16])
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
