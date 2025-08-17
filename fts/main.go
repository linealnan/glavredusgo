package main

import (
	"log"
	"strings"
	"unicode"
)

// Full-Text Search (FTS)
// Raw Text -> tokenizer->filters->tokens
// https://habr.com/ru/articles/519024/
// https://github.com/akrylysov/simplefts
func main() {
	log.Printf("Токенизируем...\n")
	tokens := tokenize("Текст для теста")
	for _, token := range tokens {
		log.Printf("%v\n", token)
	}
}

func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// Split on any character that is not a letter or a number.
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}
