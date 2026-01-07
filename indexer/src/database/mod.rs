use crate::models::Document;
use tokio_postgres::{Client, Error, NoTls};

pub struct Database {
    client: Client,
}

impl Database {
    pub fn new(client: Client) -> Self {
        Self { client }
    }

    pub async fn get_pending_documents(&self) -> Result<Vec<Document>, Error> {
        let rows = self
            .client
            .query(
                "SELECT id, url, title, content FROM documents WHERE status = 'pending' LIMIT 100",
                &[],
            )
            .await?;

        Ok(rows
            .iter()
            .map(|row| Document {
                id: row.get(0),
                url: row.get(1),
                title: row.get(2),
                content: row.get(3),
            })
            .collect())
    }

    pub async fn insert_index_entry(
        &self,
        word: &str,
        document_id: i32,
        frequency: i32,
        positions: &[i32],
    ) -> Result<(), Error> {
        self.client
            .execute(
                "INSERT INTO search_index (word, document_id, frequency, positions)
                 VALUES ($1, $2, $3, $4)
                 ON CONFLICT (word, document_id)
                 DO UPDATE SET frequency = EXCLUDED.frequency, positions = EXCLUDED.positions",
                &[&word, &document_id, &frequency, &positions],
            )
            .await?;
        Ok(())
    }

    pub async fn update_document_status(&self, doc_id: i32, status: &str) -> Result<(), Error> {
        self.client
            .execute(
                "UPDATE documents SET status = $1 WHERE id = $2",
                &[&status, &doc_id],
            )
            .await?;

        Ok(())
    }
}
