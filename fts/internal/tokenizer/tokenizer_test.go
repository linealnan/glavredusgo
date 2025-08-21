package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) {
	testCases := []struct {
		text   string
		tokens []string
	}{
		{
			text:   "",
			tokens: []string{},
		},
		{
			text:   "а",
			tokens: []string{"а"},
		},
		{
			text: "Съешь еще этих французских булок, да выпей чаю!",
			tokens: []string{
				"Съешь",
				"еще",
				"этих",
				"французских",
				"булок",
				"да",
				"выпей",
				"чаю",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.text, func(st *testing.T) {
			assert.EqualValues(st, tc.tokens, tokenize(tc.text))
		})
	}
}

func TestLowercaseFilter(t *testing.T) {
	var (
		in  = []string{"Кот", "КОТ", "кошка", "коШка"}
		out = []string{"кот", "кот", "кошка", "кошка"}
	)
	assert.Equal(t, out, lowercaseFilter(in))
}

func TestStopwordFilter(t *testing.T) {
	var (
		in  = []string{"Она", "прошла", "меж", "двух", "огней"}
		out = []string{"", "прошла", "меж", "", "огней"}
	)
	assert.Equal(t, out, stopwordFilter(in))
}

func TestStemmerFilter(t *testing.T) {
	var (
		in  = []string{"вали", "валил", "валился", "валится", "пакет", "пакетом", "пакеты"}
		out = []string{"вал", "вал", "вал", "вал", "пакет", "пакет", "пакет"}
	)
	assert.Equal(t, out, stemmerFilter(in))
}
