package blevesearch

import (
	"log"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/ru"
	"github.com/blevesearch/bleve/index/store/goleveldb"
	"github.com/blevesearch/bleve/mapping"
)

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

const indexName string = "history.bleve"

func initIndex() (bleve.Index, error) {
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
