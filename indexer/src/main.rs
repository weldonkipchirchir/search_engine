use std::collections::HashMap;
use std::env;
use std::error::Error;
use dotenv;
use tokio_postgres::{Client, NoTls};
use native_tls::TlsConnector;
use postgres_native_tls::MakeTlsConnector;

mod database;
mod indexing;
mod models;

use database::Database;
use indexing::Tokenizer;
use models::Document;

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    println!("Starting the indexing process...");
    //load environment variables
    dotenv::dotenv().ok();

    println!("Loading database configuration...");

    let db_url = env::var("DATABASE_URL")?;

    let connector = TlsConnector::builder().build().unwrap();
    let tls = MakeTlsConnector::new(connector);

    println!("Connecting to database at {}", db_url);

    /// Connect to the database
    let (client, connection) = tokio_postgres::connect(&db_url, tls).await?;

    // Spawn the connection to run in the background
    tokio::spawn(async move {
        if let Err(e) = connection.await {
            eprintln!("Database connection error: {}", e);
        }
    });

    let db = Database::new(client);
    let tokenizer = Tokenizer::new();

    println!("Fetching documents to index...");

    // Fetch pending documents
    let documents = db.get_pending_documents().await?;
    println!("Found {} documents to index", documents.len());

    for doc in documents {
        println!("Indexing document {}: {}", doc.id, doc.url);

        // Tokenize the document content
        let tokens = tokenizer.tokenize(&doc.content);

        //build word frequency map
        let word_frequencies = count_word_frequencies(&tokens);

        // Insert each word into the search index
        for (word, frequency) in word_frequencies {
            let positions: Vec<i32> = tokens
                .iter()
                .enumerate()
                .filter(|(_, w)| **w == word)
                .map(|(i, _)| i as i32)
                .collect();

            db.insert_index_entry(&word, doc.id, frequency, &positions)
                .await?;
        }

        // Update document status to 'indexed'
        db.update_document_status(doc.id, "indexed").await?;

        println!("Document {} indexed successfully", doc.id);
    }

    println!("Indexing complete.");

    Ok(())
}

fn count_word_frequencies(tokens: &[String]) -> HashMap<String, i32> {
    let mut frequencies = HashMap::new();
    for token in tokens {
        *frequencies.entry(token.clone()).or_insert(0) += 1;
    }
    frequencies
}
