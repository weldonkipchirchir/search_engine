package models

type Document struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Domain    string `json:"domain"`
	WordCount int    `json:"word_count"`
	CrawledAt string `json:"crawled_at"`
	UpdatedAt string `json:"updated_at"`
}
