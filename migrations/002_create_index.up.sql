-- Table for storing the inverted index
CREATE TABLE IF NOT EXISTS search_index (
    id SERIAL PRIMARY KEY,
    word VARCHAR(100) NOT NULL,
    document_id INTEGER REFERENCES documents(id) ON DELETE CASCADE,
    frequency INTEGER DEFAULT 1,
    positions INTEGER[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_search_index_word ON search_index(word);
CREATE INDEX IF NOT EXISTS idx_search_index_document ON search_index(document_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_search_index_unique ON search_index(word, document_id);

-- Table for search analytics
CREATE TABLE IF NOT EXISTS search_logs (
    id SERIAL PRIMARY KEY,
    query TEXT NOT NULL,
    results_count INTEGER,
    user_ip VARCHAR(45),
    user_agent TEXT,
    response_time_ms INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_search_logs_query ON search_logs(query);
CREATE INDEX IF NOT EXISTS idx_search_logs_created_at ON search_logs(created_at);
