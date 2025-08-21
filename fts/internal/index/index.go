package index

import (
	tokenizer "github.com/linealnan/glavredusgo/fts/internal/tokenizer"
)

// index is an inverted index. It maps tokens to document IDs.
type Index map[string][]int

// document represents a Wikipedia abstract dump document.
type Document struct {
	Title string
	URL   string
	Text  string
	ID    int
}

// add adds documents to the index.
func (idx Index) Add(docs []Document) {
	for _, doc := range docs {
		for _, token := range tokenizer.Tokenize(doc.Text) {
			ids := idx[token]
			if ids != nil && ids[len(ids)-1] == doc.ID {
				// Don't add same ID twice.
				continue
			}
			idx[token] = append(ids, doc.ID)
		}
	}
}

// intersection returns the set intersection between a and b.
// a and b have to be sorted in ascending order and contain no duplicates.
func intersection(a []int, b []int) []int {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	r := make([]int, 0, maxLen)
	var i, j int
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		} else {
			r = append(r, a[i])
			i++
			j++
		}
	}
	return r
}

// search queries the index for the given text.
func (idx Index) Search(text string) []int {
	var r []int
	for _, token := range tokenizer.Tokenize(text) {
		if ids, ok := idx[token]; ok {
			if r == nil {
				r = ids
			} else {
				r = intersection(r, ids)
			}
		} else {
			// Token doesn't exist.
			return nil
		}
	}
	return r
}
