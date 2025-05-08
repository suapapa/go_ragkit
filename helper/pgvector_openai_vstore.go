package helper

import (
	"cmp"

	oai "github.com/openai/openai-go"
	oai_option "github.com/openai/openai-go/option"
	oai_embedder "github.com/suapapa/go_ragkit/embedder/openai"
	pgvector_vstore "github.com/suapapa/go_ragkit/vector_store/pgvector"

	ragkit "github.com/suapapa/go_ragkit"
)

func NewPGVectorOpenAIVectorStore(
	pgvectorConnStr string,
	vectorDBClassName string,
	oaiEmbedModel string,
	oaiOptions ...oai_option.RequestOption,
) (ragkit.VectorStore, error) {
	// initialize openai
	opts := []oai_option.RequestOption{
		oai_option.WithAPIKey(oaiApiKey),
	}
	opts = append(opts, oaiOptions...)

	oaiClient := oai.NewClient(opts...)
	embedder := oai_embedder.New(&oaiClient, cmp.Or(oaiEmbedModel, DefaultOAIEmbedModel))

	// initialize pgvector
	pgvector := pgvector_vstore.New(pgvectorConnStr, 1536, vectorDBClassName, embedder)

	return pgvector, nil
}
