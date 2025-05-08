package pgvector

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	ragkit "github.com/suapapa/go_ragkit"
)

var _ ragkit.VectorStore = &PGVector{}

type PGVector struct {
	className string
	conn      *pgx.Conn
	embedder  ragkit.Embedder

	mu sync.Mutex
}

func New(connStr string, className string, embedder ragkit.Embedder) *PGVector {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	ret := &PGVector{
		className: className,
		conn:      conn,
		embedder:  embedder,
	}

	err = ret.ensureTable(context.Background())
	if err != nil {
		log.Fatalf("Failed to ensure table: %v", err)
	}

	return ret
}

// ensureTable creates the table if it doesn't exist
func (p *PGVector) ensureTable(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Create the table if it doesn't exist
	_, err := p.conn.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id TEXT PRIMARY KEY,
			text TEXT NOT NULL,
			metadata JSONB,
			embedding vector(1536)
		)
	`, p.className))
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Create index for vector similarity search
	_, err = p.conn.Exec(ctx, fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS %s_embedding_idx ON %s 
		USING ivfflat (embedding vector_cosine_ops)
		WITH (lists = 100)
	`, p.className, p.className))
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}

func (p *PGVector) Index(ctx context.Context, docs ...ragkit.Document) ([]string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var ids []string
	for _, doc := range docs {
		var embedding []float32
		var err error

		if doc.ID == "" {
			doc.ID = ragkit.GenerateID(doc.Text, doc.Metadata)
		}

		// Use provided vector if available, otherwise generate one
		if doc.Vector != nil {
			embedding = doc.Vector
		} else {
			vectors, err := p.embedder.EmbedTexts(ctx, doc.Text)
			if err != nil {
				return nil, err
			}
			embedding = vectors[0]
		}

		// Insert document into database
		_, err = p.conn.Exec(ctx, fmt.Sprintf(`
			INSERT INTO %s (id, text, metadata, embedding) 
			VALUES ($1, $2, $3, $4)
		`, p.className), doc.ID, doc.Text, doc.Metadata, pgvector.NewVector(embedding))
		if err != nil {
			return nil, err
		}
		ids = append(ids, doc.ID)
	}
	return ids, nil
}

func (p *PGVector) Delete(ctx context.Context, id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	_, err := p.conn.Exec(ctx, fmt.Sprintf(`
		DELETE FROM %s WHERE id = $1
	`, p.className), id)
	return err
}

func (p *PGVector) Exists(ctx context.Context, id string) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var exists bool
	err := p.conn.QueryRow(ctx, fmt.Sprintf(`
		SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)
	`, p.className), id).Scan(&exists)
	return exists, err
}

func (p *PGVector) Retrieve(ctx context.Context, query []float32, topK int, metadataFieldNames ...string) ([]ragkit.RetrievedDoc, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Perform vector similarity search
	rows, err := p.conn.Query(ctx, fmt.Sprintf(`
		SELECT text, metadata, embedding 
		FROM %s 
		ORDER BY embedding <=> $1 
		LIMIT $2
	`, p.className), pgvector.NewVector(query), topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ragkit.RetrievedDoc
	for rows.Next() {
		var doc ragkit.RetrievedDoc
		var embedding pgvector.Vector
		err := rows.Scan(&doc.Text, &doc.Metadata, &embedding)
		if err != nil {
			return nil, err
		}
		doc.Vector = embedding.Slice()
		results = append(results, doc)
	}
	return results, rows.Err()
}

func (p *PGVector) RetrieveText(ctx context.Context, text string, topK int, metadataFieldNames ...string) ([]ragkit.RetrievedDoc, error) {
	vectors, err := p.embedder.EmbedTexts(ctx, text)
	if err != nil {
		return nil, err
	}
	return p.Retrieve(ctx, vectors[0], topK, metadataFieldNames...)
}

func (p *PGVector) String() string {
	return fmt.Sprintf("PGVector(table: %s, embedder: %s)", p.className, p.embedder)
}
