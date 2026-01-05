package main

import (
	"context"
	"crawler/internal/fetcher"
	"crawler/internal/parser"
	"crawler/internal/storage"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subosito/gotenv"
)

func main() {
	//load environment variables
	if err := gotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	//build database connection string
	dbURL := getEnv("DATABASE_URL", "")
	fmt.Printf("Database URL: %s\n", dbURL)

	fmt.Println("Starting crawler...")

	//connect to database
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	log.Println("Connected to database successfully")

	// Here you can initialize and start your crawler components
	log.Println("Crawler started...")
	store := storage.NewPostgresStorage(pool)
	htmlParser := parser.NewHTMLParser()
	fetcher := fetcher.NewFetcher()

	//ctx := context.Background()

	//main crawling loop (simplified)
	for {
		//get urls to crawl
		urls, err := store.GetPendingURLs(10)
		if err != nil {
			log.Printf("Error fetching URLs: %v\n", err)
			continue
		}

		if len(urls) == 0 {
			log.Println("No URLs to crawl, sleeping...")
			//sleep or wait before next iteration
			continue
		}

		//process each url
		for _, url := range urls {
			log.Printf("Crawling URL: %s\n", url)

			//mark url as in-progress
			if err := store.UpdateURLStatus(url, "processing"); err != nil {
				log.Printf("Failed to update status: %v", err)
				continue
			}

			//fetch page
			html, err := fetcher.Fetch(url)
			if err != nil {
				log.Printf("Error fetching URL %s: %v\n", url, err)
				continue
			}

			//parse page
			doc, err := htmlParser.Parse(url, html)
			if err != nil {
				log.Printf("Error parsing URL %s: %v\n", url, err)
				continue
			}

			//save document
			err = store.SaveDocument(doc)
			if err != nil {
				log.Printf("Error saving document for URL %s: %v\n", url, err)
				continue
			}

			//extract and enqueue links
			links := htmlParser.ExtractLinks(html, url)
			for _, link := range links {
				if err := store.AddToCrawlQueue(link, 0); err != nil {
					log.Printf("Failed to queue link %s: %v", link, err)
				}

			}

			// Mark as completed
			store.UpdateURLStatus(url, "completed")
			log.Printf("Successfully crawled: %s", url)

			// Be nice to servers
			time.Sleep(2 * time.Second)
		}
	}

}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
