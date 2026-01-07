use unicode_segmentation::UnicodeSegmentation;

pub struct Tokenizer;

impl Tokenizer {
    pub fn new() -> Self {
        Self
    }

    pub fn tokenize(&self, text: &str) -> Vec<String> {
        text.unicode_words()
            .map(|s| s.to_lowercase())
            .filter(|s| s.len() > 2)
            .collect()
    }
}
