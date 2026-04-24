package blevesearch

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/ru"
	"github.com/blevesearch/bleve/index/store/goleveldb"
	"github.com/blevesearch/bleve/mapping"
)

type Searcher interface {
	Search(query string) *bleve.SearchResult
}

type BleaveSearch struct {
	Index bleve.Index
}

func NewBleaveSearch() *BleaveSearch {
	i, err := initIndex()

	if err != nil {
		log.Fatal(err)
	}

	return &BleaveSearch{
		Index: i,
	}
}

func initIndex() (bleve.Index, error) {
	// Получаем путь к индексу из переменной окружения или используем значение по умолчанию
	indexName := os.Getenv("BLEVE_PATH")
	if indexName == "" {
		indexName = "history.bleve"
	}

	log.Printf("Using Bleve index path: %s", indexName)

	index, err := bleve.Open(indexName)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := buildMapping()
		kvStore := goleveldb.Name
		kvConfig := map[string]interface{}{
			"create_if_missing": true,
		}

		index, err = bleve.NewUsing(indexName, mapping, "upside_down", kvStore, kvConfig)
		if err != nil {
			return nil, err
		}
	}
	return index, nil
}

func (bs *BleaveSearch) Search(query string) *bleve.SearchResult {
	// Создаем Query для совпадений фраз в индексе. Анализатор выбирается по полю. Ввод анализируется этим анализатором. Токенезированные выражения от анализа используются для посторения поисковой фразы. Результирующие документы должны совпадать с этой фразой.
	mq := bleve.NewMatchPhraseQuery(query)
	// Создаем Query для поиска значений в индексе по регулярному выражению
	// rq := bleve.NewRegexpQuery(query)

	// qsq := bleve.NewQueryStringQuery(query)

	q := bleve.NewDisjunctionQuery(mq)

	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Highlight = bleve.NewHighlight()
	searchRequest.Fields = []string{"ID", "GroupName", "GroupID", "Text"}
	searchRequest.Size = 100

	searchResults, err := bs.Index.Search(searchRequest)

	if err != nil {
		log.Fatal(err)
	}

	return searchResults
}

func buildMapping() *mapping.IndexMappingImpl {

	// ruFieldMapping := bleve.NewTextFieldMapping()
	// ruFieldMapping.Analyzer = ru.AnalyzerName

	// eventMapping := bleve.NewDocumentMapping()
	// eventMapping.AddFieldMappingsAt("Text", ruFieldMapping)

	mapping := bleve.NewIndexMapping()
	//mapping.AddDocumentMapping("vkgoups", eventMapping)
	mapping.DefaultAnalyzer = ru.AnalyzerName
	return mapping
}

func (bs *BleaveSearch) TgSearch(searchPhrase string) string {
	searchResult := bs.Search(searchPhrase)

	var builder strings.Builder

	builder.WriteString("Найдено документов: \n")
	for _, hit := range searchResult.Hits {
		fields := hit.Fields
		builder.WriteString(fmt.Sprintf("https://vk.com/%s?w=wall-%s_%s\n", fields["GroupName"], fields["GroupID"], fields["ID"]))
	}

	return builder.String()
	// fmt.Fprintf(w, "Результат поиска: %s", searchResult)
}
