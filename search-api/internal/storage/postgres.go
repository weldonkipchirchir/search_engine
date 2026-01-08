package storage

import (
	"context"
	"fmt"
	"search-api/pkg/models"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(pool *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{pool: pool}
}

func (s *PostgresStorage) Search(ctx context.Context, words string, limit, offset int) ([]models.SearchResult, error) {
	// Implementation of search logic using PostgreSQL full-text search
	if len(words) == 0 {
		return []models.SearchResult{}, nil
	}

	placeholders := make([]string, len(words))
	args := make([]interface{}, len(words))
	for i := range words {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = string(words[i])
	}

	query := fmt.Sprintf(`
		SELECT d.id, d.title, d.url, LEF(d.content, 200) as snippet, Sum(si.frequency) as score
		FROM search_index si
		JOIN documents d ON si.document_id = d.id
		where si.word IN (%s)
		GROUP BY d.id, d.title, d.url, snippet
		ORDER BY score DESC
		LIMIT $%d OFFSET $%d
	`, strings.Join(placeholders, ","), len(words)+1, len(words)+2)

	rows, err := s.pool.Query(ctx, query, append(args, limit, offset)...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []models.SearchResult
	rank := offset + 1
	for rows.Next() {
		var res models.SearchResult
		err := rows.Scan(&res.ID, &res.Title, &res.URL, &res.Snippet, &res.Score)
		if err != nil {
			return nil, err
		}
		res.Rank = rank
		rank++
		results = append(results, res)
	}

	return results, rows.Err()
}

func (s *PostgresStorage) LogSearch(ctx context.Context, query string, resultsCount int, ip, userAgent string, responseTime int) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO search_logs (query, results_count, user_ip, user_agent, response_time_ms)
		VALUES ($1, $2, $3, $4, $5)
	`, query, resultsCount, ip, userAgent, responseTime)
	return err
}

func (s *PostgresStorage) GetStats(ctx context.Context) (map[string]any, error) {
	stats := make(map[string]interface{})

	//get total documents
	var totalDocs int
	err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM documents`).Scan(&totalDocs)
	if err != nil {
		return nil, err
	}

	stats["total_documents"] = totalDocs

	// get indexed words
	var indexedWords int
	err = s.pool.QueryRow(ctx, "SELECT COUNT(*) FROM documents WHERE status = 'indexed'").Scan(&indexedWords)
	if err != nil {
		return nil, err
	}

	stats["indexed_documents"] = indexedWords

	// get total searches
	var totalSearches int
	err = s.pool.QueryRow(ctx, "SELECT COUNT(*) FROM search_logs").Scan(&totalSearches)
	if err != nil {
		return nil, err
	}

	stats["total_searches"] = totalSearches

	//get top queries
	rows, err := s.pool.Query(ctx, `
		SELECT query, COUNT(*) as count
		FROM search_logs
		GROUP BY query
		ORDER BY count DESC	
		LIMIT 10
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	topQueries := make([]map[string]interface{}, 0)

	for rows.Next() {
		var query string
		var count int
		err := rows.Scan(&query, &count)
		if err != nil {
			return nil, err
		}
		topQueries = append(topQueries, map[string]interface{}{
			"query": query,
			"count": count,
		})
	}
	stats["top_queries"] = topQueries

	return stats, nil
}
