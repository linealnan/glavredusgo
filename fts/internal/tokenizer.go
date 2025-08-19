package tokenizer

import (
	"strings"
	"unicode"

	snowballeng "github.com/grecod-oss/snowball/russian"
	sw "github.com/toadharvard/stopwords-iso"
)

func Tokenize(text string) []string {
	tokens := tokenize(text)
	tokens = lowercaseFilter(tokens)
	tokens = stopwordFilter(tokens)
	tokens = stemmerFilter(tokens)

	return tokens
}

func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// lowercaseFilter returns a slice of tokens normalized to lower case.
func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}

// https://pkg.go.dev/github.com/toadharvard/stopwords-iso#section-readme
func stopwordFilter(tokens []string) []string {
	stopwordsMapping, _ := sw.NewStopwordsMapping()
	language := "ru"
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		r = append(r, stopwordsMapping.ClearStringByLang(token, language))
	}
	return r
}

// https://pkg.go.dev/github.com/grecod-oss/snowball/russian
func stemmerFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = snowballeng.Stem(token, false)
	}
	return r
}
