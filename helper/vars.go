package helper

import (
	"cmp"
	"os"
)

var (
	weaviateAddr = cmp.Or(os.Getenv("WEAVIATE_ADDR"), "http://localhost:8080")

	ollamaAddr = cmp.Or(os.Getenv("OLLAMA_ADDR"), "http://localhost:11434")
	oaiApiKey  = cmp.Or(os.Getenv("OPENAI_SECRET_KEY"), "")

	defaultOllamaEmbedModel = cmp.Or(os.Getenv("OLLAMA_EMBED_MODEL"), "bge-m3:latest")
	defaultOaiEmbedModel    = cmp.Or(os.Getenv("OPEAI_EMBED_MODEL"), "text-embedding-3-small")
)
