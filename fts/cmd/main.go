package main

import (
	"log"

	tokenizer "github.com/linealnan/glavredusgo/fts/internal/tokenizer"
)

// Full-Text Search (FTS)
// Raw Text -> tokenizer->filters->tokens
// https://habr.com/ru/articles/519024/
// https://github.com/akrylysov/simplefts
func main() {
	log.Printf("Токенизируем...\n")
	tokens := tokenizer.Tokenize("Текст для теста")
	for _, token := range tokens {
		log.Printf("%v\n", token)
	}
}
