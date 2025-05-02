package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	ragkit "github.com/suapapa/go_ragkit"

	oai "github.com/openai/openai-go"
	oai_option "github.com/openai/openai-go/option"
	oai_embedder "github.com/suapapa/go_ragkit/embedder/openai"
	weaviate_vstore "github.com/suapapa/go_ragkit/vector_store/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

var (
	oaiSecretKey  = cmp.Or(os.Getenv("OPENAI_SECRET_KEY"), "")
	oaiEmbedModel = cmp.Or(os.Getenv("OPEAI_EMBED_MODEL"), "text-embedding-3-small")
	weaviateAddr  = cmp.Or(os.Getenv("WEAVIATE_ADDR"), "http://localhost:8080")
)

func main() {
	// define embedder
	log.Println("defining embedder")
	log.Println("oaiSecretKey", oaiSecretKey)

	oaiClient := oai.NewClient(
		// oai_option.WithEnvironmentProduction(),
		oai_option.WithAPIKey(oaiSecretKey),
	)
	embedder := oai_embedder.New(&oaiClient, oaiEmbedModel)

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
	vstore := weaviate_vstore.New(weaviateClient, "FamilyTree", embedder)
	fmt.Printf("vectorizer: %s\n", vstore)

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
		if exist, err := vstore.Exists(ctx, doc.ID); err != nil {
			panic(err)
		} else if exist {
			// log.Printf("document %s already exists", doc.ID)
			continue
		}

		_, err = vstore.Index(ctx, doc)
		if err != nil {
			panic(err)
		}
	}

	// retrieve documents
	log.Println("retrieving documents")
	query := "희동이와 고길동의 관계?"
	results, err := vstore.RetrieveText(ctx, query, 3)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s:\n", query)
	for _, result := range results {
		fmt.Println("-", result.Text)
	}
}
