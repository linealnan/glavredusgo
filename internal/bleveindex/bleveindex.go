package bleveindex

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/ru"
	"github.com/blevesearch/bleve/index/store/goleveldb"
	"github.com/blevesearch/bleve/mapping"
)

const indexName string = "history.bleve"

func NewBleveIndex() (bleve.Index, error) {
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
