# ragkit: A Go package for document indexing and retrieval in RAG systems

![ragkit_logo](_asset/ragkit_logo_256.webp)

ragkit is a Go package designed to simplify the implementation of Retrieval-Augmented Generation (RAG) systems.
It includes the definition and implementation of a VectorStore interface that performs document indexing and retrieval,
providing tools for vectorization and semantic search capabilities.

- [Package documentation](https://pkg.go.dev/github.com/suapapa/go_ragkit).

## Installation

```sh
go get github.com/suapapa/go_ragkit
```

## Quick Start

```go
import (
    // ...
    ragkit "github.com/suapapa/go_ragkit"
    ragkit_helper "github.com/suapapa/go_ragkit/helper"
)

func main() {
    // Initialize vector store (pgvector + OpenAI-compatible embedder)
    // Uses PGVECTOR_CONN, OPENAI_ENDPOINT, OPENAI_APIKEY, OPENAI_EMBED_MODEL, OPENAI_EMBED_DIM
    vstore, err := ragkit_helper.NewPGVectorOpenAIVectorStore(
        "dooly_family", // vector DB table name
        ragkit_helper.DefaultOAIEmbedModel(),
    )
    if err != nil {
        panic(err)
    }

    // Create documents from text
    docs := ragkit.MakeDocsFromTexts(
        []string{
            "고길동의 집에는 둘리, 도우너, 또치, 희동이, 철수, 영희가 살고 있다.",
            "희동이는 고길동의 조카이다.",
            // ...
        },
        nil,
    )

    // Index documents
    _, err = vstore.Index(context.Background(), docs...)
    if err != nil {
        panic(err)
    }

    // Perform semantic search
    query := "고길동과 희동이의 관계?"
    results, err := vstore.RetrieveText(context.Background(), query, 1)
    if err != nil {
        panic(err)
    }

    // Print result
    fmt.Println(results[0].Text) // <- 희동이는 고길동의 조카이다.
}
```

## Examples

Pre-requirement for pgvector helpers — Postgres with the `vector` extension, and `PGVECTOR_CONN` set.

Weaviate examples still live under [`examples/`](examples/) (`WEAVIATE_ADDR`).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
