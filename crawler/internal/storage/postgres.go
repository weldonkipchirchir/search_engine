package storage

import (
	"context"
	"crawler/pkg/models"
	"crypto/sha256"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(pool *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{pool: pool}
}

func (ps *PostgresStorage) Close() {
	ps.pool.Close()
}

// save document to postgres
func (ps *PostgresStorage) SaveDocument(doc *models.Document) error {
	hash := generateHash(doc.Content)

	query := `
		INSERT INTO documents (url, title, content, content_hash, domain, word_count, status) VALUES ($1, $2, $3, $4, $5, $6, 'pending') ON CONFLICT (url) DO UPDATE SET title = EXCLUDED.title, content = EXCLUDED.content, content_hash = EXCLUDED.content_hash, updated_at = CURRENT_TIMESTAMP, status = 'pending' WHERE documents.content_hash != EXCLUDED.content_hash`

	_, err := ps.pool.Exec(
		context.Background(),
		query,
		doc.URL,
		doc.Title,
		doc.Content,
		hash,
		doc.Domain,
		doc.WordCount,
	)
	return err
}

func generateHash(content string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
}

func (ps *PostgresStorage) GetPendingURLs(limit int) ([]string, error) {
	query := `
		SELECT url
		FROM crawl_queue
		WHERE status = 'pending'
		ORDER BY priority DESC, created_at ASC
		LIMIT $1
	`

	rows, err := ps.pool.Query(context.Background(), query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make([]string, 0, limit)
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (ps *PostgresStorage) UpdateURLStatus(url string, status string) error {
	query := `
		UPDATE crawl_queue
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE url = $2
	`
	_, err := ps.pool.Exec(context.Background(), query, status, url)
	return err
}

func (ps *PostgresStorage) AddToCrawlQueue(url string, priority int) error {
	query := `
		INSERT INTO crawl_queue (url, priority) VALUES ($1, $2)
		ON CONFLICT (url) DO NOTHING
	`
	_, err := ps.pool.Exec(context.Background(), query, url, priority)
	return err
}
