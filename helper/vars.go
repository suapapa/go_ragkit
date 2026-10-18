package helper

import (
	"cmp"
	"os"
	"strconv"
)

func weaviateAddr() string {
	return cmp.Or(os.Getenv("WEAVIATE_ADDR"), "http://localhost:8080")
}

func pgvectorConn() string {
	return os.Getenv("PGVECTOR_CONN")
}

func oaiEndpoint() string {
	return os.Getenv("OPENAI_ENDPOINT")
}

func oaiApiKey() string {
	return os.Getenv("OPENAI_APIKEY")
}

// DefaultOAIEmbedModel returns OPENAI_EMBED_MODEL or text-embedding-3-small.
func DefaultOAIEmbedModel() string {
	return cmp.Or(os.Getenv("OPENAI_EMBED_MODEL"), "text-embedding-3-small")
}

// DefaultOAIEmbedDim returns OPENAI_EMBED_DIM or 1536 (text-embedding-3-small).
// Must match the embedding model output and existing pgvector column size.
func DefaultOAIEmbedDim() int {
	if v := os.Getenv("OPENAI_EMBED_DIM"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 1536
}

// Ollama defaults (no env vars; prefer OpenAI-compatible endpoint instead)
const (
	ollamaAddr              = "http://localhost:11434"
	DefaultOllamaEmbedModel = "bge-m3:latest"
)
