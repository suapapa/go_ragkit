package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	ragkit "github.com/suapapa/go_ragkit"

	ollama_api "github.com/ollama/ollama/api"
	ollama_embedder "github.com/suapapa/go_ragkit/embedder/ollama"
	weaviate_vectorizer "github.com/suapapa/go_ragkit/vectorizer/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

var (
	ollamaAddr       = cmp.Or(os.Getenv("OLLAMA_ADDR"), "http://localhost:11434")
	ollamaEmbedModel = cmp.Or(os.Getenv("OLLAMA_EMBED_MODEL"), "bge-m3:latest")
	weaviateAddr     = cmp.Or(os.Getenv("WEAVIATE_ADDR"), "http://localhost:8080")
)

func main() {
	// define embedder
	log.Println("defining embedder")
	ollamaURL, err := url.Parse(ollamaAddr)
	if err != nil {
		panic(err)
	}
	ollamaClient := ollama_api.NewClient(ollamaURL, http.DefaultClient)
	embedder := ollama_embedder.NewOllama(ollamaClient, ollamaEmbedModel)

	// define vectorizer
	log.Println("defining vectorizer")
	weaviateURL, err := url.Parse(weaviateAddr)
	if err != nil {
		panic(err)
	}
	weaviateClient, err := weaviate.NewClient(weaviate.Config{
		Host:   weaviateURL.Host,
		Scheme: weaviateURL.Scheme,
	})
	if err != nil {
		panic(err)
	}
	vectorizer := weaviate_vectorizer.NewWeaviate(weaviateClient, "FamilyTree", embedder)

	// index documents
	log.Println("indexing documents")
	docs := ragkit.MakeDocsFromTexts(
		[]string{
			"고길동의 집에는 둘리, 도우너, 또치, 희동이, 철수, 영희가 살고 있다.",
			"희동이는 고길동의 조카이다.",
			"철수는 고길동의 아들이다.",
			"영희는 고길동의 딸이다.",
			"둘리는 고길동이 입양한 아들이다.",
			"도우너는 고길동이 입양한 아들이다.",
			"또치는 고길동이 입양한 딸이다.",
		},
		nil,
	)
	ctx := context.Background()
	for _, doc := range docs {
		if exist, err := vectorizer.Exists(ctx, doc.ID); err != nil {
			panic(err)
		} else if exist {
			log.Printf("document %s already exists", doc.ID)
			continue
		}

		_, err = vectorizer.Index(ctx, doc)
		if err != nil {
			panic(err)
		}
	}

	// retrieve documents
	log.Println("retrieving documents")
	query := "희동이와 고길동의 관계?"
	results, err := vectorizer.RetrieveText(ctx, query, 10)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s => %s\n", query, results[0].Text)
}
