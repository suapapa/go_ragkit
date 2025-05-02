package common

import (
	"context"
	"fmt"
	"log"

	ragkit "github.com/suapapa/go_ragkit"
)

func IndexAndRetriveExample(vstore ragkit.VectorStore) {
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

		_, err := vstore.Index(ctx, doc)
		if err != nil {
			panic(err)
		}
	}

	// retrieve documents
	log.Println("retrieving documents")
	query := "고길동 희동이의 관계?"
	results, err := vstore.RetrieveText(ctx, query, 3)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s:\n", query)
	for _, result := range results {
		fmt.Println("-", result.Text)
	}
}
