module github.com/linealnan/glavredusgo/fts

go 1.24.2

replace github.com/linealnan/glavredusgo/fts/internal => ./internal/fts

require github.com/stretchr/testify v1.10.0

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/grecod-oss/snowball v0.0.0-20210330145637-4d5d205112d0
	github.com/kljensen/snowball v0.10.0
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/toadharvard/stopwords-iso v0.1.5
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
