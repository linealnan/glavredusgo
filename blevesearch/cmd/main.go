package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/ru"
	"github.com/blevesearch/bleve/index/store/goleveldb"
	"github.com/blevesearch/bleve/mapping"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	vkapi "github.com/romanraspopov/golang-vk-api"
)

type IndexData struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type LoadedPost struct {
	// Text string `xml:"abstract"`
	Text string
	ID   int
}

type MockGroup struct {
	Name string
}

func main() {
	var posts []LoadedPost
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("VK_API_TOKEN")
	client, err := vkapi.NewVKClientWithToken(token, nil, true)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", "glavredus.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	posts = loadGroupsData(client, db)

	indexName := "history.bleve"
	index, err := bleve.Open(indexName)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := buildMapping()
		kvStore := goleveldb.Name
		kvConfig := map[string]interface{}{
			"create_if_missing": true,
			//		"write_buffer_size":         536870912,
			//		"lru_cache_capacity":        536870912,
			//		"bloom_filter_bits_per_key": 10,
		}

		index, err = bleve.NewUsing(indexName, mapping, "upside_down", kvStore, kvConfig)
	}

	err = index.Index(indexName, posts)

	if err != nil {
		log.Fatal(err)
	}

	query := "Пара рисунков"

	// Создаем Query для совпадений фраз в индексе. Анализатор выбирается по полю. Ввод анализируется этим анализатором. Токенезированные выражения от анализа используются для посторения поисковой фразы. Результирующие документы должны совпадать с этой фразой.
	mq := bleve.NewMatchPhraseQuery(query)
	// Создаем Query для поиска значений в индексе по регулярному выражению
	rq := bleve.NewRegexpQuery(query)

	q := bleve.NewDisjunctionQuery(mq, rq)

	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Highlight = bleve.NewHighlight()
	//searchRequest.Fields = []string{"ID"}

	searchResults, err := index.Search(searchRequest)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Результат: %v\n", searchResults)
}

func buildMapping() *mapping.IndexMappingImpl {
	ruFieldMapping := bleve.NewTextFieldMapping()
	ruFieldMapping.Analyzer = ru.AnalyzerName

	eventMapping := bleve.NewDocumentMapping()
	eventMapping.AddFieldMappingsAt("Text", ruFieldMapping)

	mapping := bleve.NewIndexMapping()
	mapping.DefaultMapping = eventMapping
	mapping.DefaultAnalyzer = ru.AnalyzerName
	return mapping
}

// func (ss *SearchService) Search(query, channel string) (*bleve.SearchResult, error) {
// 	stringQuery := fmt.Sprintf("/.*%s.*/", query)
// 	ss.logger.Info(query)
// 	ch := bleve.NewTermQuery(channel)
// 	mq := bleve.NewMatchPhraseQuery(query)
// 	rq := bleve.NewRegexpQuery(query)
// 	qsq := bleve.NewQueryStringQuery(stringQuery)
// 	q := bleve.NewDisjunctionQuery(ch, mq, rq, qsq)
// 	search := bleve.NewSearchRequest(q)
// 	search.Fields = []string{"username", "message", "channel", "timestamp"}
// 	return ss.index.Search(search)
// }

func getGroups() []MockGroup {
	return []MockGroup{
		{Name: "trenchcrusade"},
		// {Name: "nvp_73"},
		// {Name: "ad_ka4alka"},
	}
}

func loadGroupsData(client *vkapi.VKClient, db *sql.DB) []LoadedPost {
	var posts []LoadedPost
	groups := getGroups()
	log.Printf("Загрузка данных групп\n")
	for _, group := range groups {
		posts = append(posts, getAndIndexedWallPostByGroupName(client, group.Name, db)...)
	}

	return posts
}

func getAndIndexedWallPostByGroupName(client *vkapi.VKClient, groupName string, db *sql.DB) []LoadedPost {
	var posts []LoadedPost
	var indexedPost LoadedPost

	wall, err := client.WallGet(groupName, 100, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range wall.Posts {
		indexedPost.ID = post.ID
		indexedPost.Text = post.Text

		posts = append(posts, indexedPost)
		//log.Printf("Wall post: %v\n", post.ID)

		result, err := db.Exec("INSERT OR IGNORE INTO WallPost (id, text) VALUES ($1, $2)",
			post.ID, post.Text)
		fmt.Println(post.ID)
		if err != nil {
			panic(err)
		}
		fmt.Println(result.LastInsertId()) // id последнего добавленного объекта
		fmt.Println(result.RowsAffected()) // количество добавленных строк
	}

	return posts
}

func initDBSchema(db *sql.DB) {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS vkgroup (
			name string PRIMARY KEY NOT NULL,
		);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}
