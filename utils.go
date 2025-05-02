// Package ragkit provides utility functions for RAG (Retrieval-Augmented Generation) applications.
package ragkit

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
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

// GenerateID creates a deterministic UUID v5 from the input text and metadata.
// The generated ID is guaranteed to be unique for different inputs.
func GenerateID(text string, metadata map[string]any) string {
	// Create a namespace UUID (using SHA-256 of "ragkit" as the namespace)
	namespace := uuid.NewSHA1(uuid.Nil, []byte("ragkit"))

	// Combine text and metadata into a single string
	var input string
	if metadata != nil {
		metadataB, _ := json.Marshal(metadata)
		input = text + string(metadataB)
	} else {
		input = text
	}

	// Generate UUID v5 using the namespace and input
	id := uuid.NewSHA1(namespace, []byte(input))
	return id.String()
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
