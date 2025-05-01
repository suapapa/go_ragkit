package ragkit

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

func MakeDocsFromTexts(texts []string, metadata map[string]any) []Document {
	docs := make([]Document, len(texts))
	for i, text := range texts {
		docs[i] = Document{ID: GenerateID(text), Text: text, Metadata: metadata}
	}
	return docs
}

// GenerateID: Generate a unique ID for a document
func GenerateID(text string) string {
	hash := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", hash[0:4], hash[4:6], hash[6:8], hash[8:10], hash[10:16])
}

func ToJSONStr(v any) string {
	buff := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buff)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	enc.Encode(v)
	return buff.String()
}
