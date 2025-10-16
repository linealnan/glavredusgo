package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/ru"
	"github.com/blevesearch/bleve/index/store/goleveldb"
	"github.com/blevesearch/bleve/mapping"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	vkapi "github.com/romanraspopov/golang-vk-api"
	"github.com/urfave/cli/v2"
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
type VkGroup struct {
	Name string
}

// Service interface is base service, with simple API
type Service interface {
	Init() error
	Run() error
	Name() string
	Stop()
}

type GlavredusFinderService struct {
	services  map[string]Service
	waitGroup sync.WaitGroup

	logger log.Logger
}

const indexName string = "history.bleve"

var index bleve.Index

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("sqlite3", "glavredus.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	deleteVkGroupTable(db)
	loadInitSchema(db)
	loadSchoolVkGroups(db)

	index, err = bleve.Open(indexName)
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
		if err != nil {
			fmt.Println("Error starting the server:", err)
		}
	}

	http.HandleFunc("/", formHandler)

	// Запускаем сервер
	fmt.Println("Starting server at port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}

	app := &cli.App{
		Name:  "glavredus",
		Usage: "Поиск по группам",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "index",
				Aliases: []string{"i"},
				Value:   "",
				Usage:   "загрузить данные групп в индекс",
			},
		},
		Action: func(c *cli.Context) error {
			name := "Значение команды"
			if c.NArg() > 0 {
				name = c.Args().Get(0)
			}

			if c.String("i") == "force" {

				token := os.Getenv("VK_API_TOKEN")
				client, err := vkapi.NewVKClientWithToken(token, nil, true)
				if err != nil {
					log.Fatal(err)
				}
				posts := loadGroupsData(client, db)
				err = index.Index(indexName, posts)

				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Индекс обновлен\n")

			} else {
				fmt.Printf("Hello, %s\n", name)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func search(query string) *bleve.SearchResult {
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

	return searchResults
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method == http.MethodGet {
		// Отображаем форму
		fmt.Fprintf(w, `
			<html>
				<form method="POST">
					<input type="text" name="search" placeholder="Введите текст">
					<button type="submit">Найти</button>
				</form>
			</html
        `)
	} else if r.Method == http.MethodPost {
		// Обрабатываем данные формы
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "Error parsing form: %v", err)
			return
		}

		name := r.FormValue("search")
		result := search(name)
		log.Printf("Результат поиска: %s", result)
		fmt.Fprintf(w, "Результат поиска: %s", result)

	}
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

func loadGroupsData(client *vkapi.VKClient, db *sql.DB) []LoadedPost {
	var posts []LoadedPost
	groups := getVkGroups(db)
	log.Printf("Получение данных групп\n")
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
	}

	return posts
}

func deleteVkGroupTable(db *sql.DB) {
	sql := `DELETE FROM vkgroup;`

	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
}

func loadInitSchema(db *sql.DB) {
	createVkgroupSQL := `
		CREATE TABLE IF NOT EXISTS vkgroup (
			name string PRIMARY KEY NOT NULL
		);`

	_, err := db.Exec(createVkgroupSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func getVkGroups(db *sql.DB) []VkGroup {
	rows, err := db.Query(`SELECT name FROM vkgroup;`)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var vkgroups []VkGroup

	// Итерируемся по строкам
	for rows.Next() {
		var p VkGroup
		if err := rows.Scan(&p.Name); err != nil {
			log.Fatal(err)
		}
		vkgroups = append(vkgroups, p)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return vkgroups
}

func loadSchoolVkGroups(db *sql.DB) {
	vkgroups := []VkGroup{
		{"club194809745"},
		{"club214119048"},
		{"club202724280"},
		{"club185982638"},
		{"club205401563"},
		{"club205402681"},
		{"detskisad15"},
		{"16detskiysad"},
		{"club205401551"},
		{"doy19"},
		{"club109060055"},
		{"club205401929"},
		// {Name: "club205400972"},
		// {Name: "club182072023"},
		// {Name: "club195576991"},
		// {Name: "club147892228"},
		// {Name: "club187951249"},
		// {Name: "sadik31krs"},
		// {Name: "club205420428"},
		// {Name: "gdboy35"},
		// {Name: "club205443755"},
		// {Name: "dc39spb"},
		// {Name: "club170186955"},
		// {Name: "gbdou41krspb"},
		// {Name: "club216246675"},
		// {Name: "club205406349"},
		// {Name: "club203026295"},
		// {Name: "dc5krs"},
		// {Name: "ds51krs"},
		// {Name: "gbdouds52"},
		// {Name: "club13309436"},
		// {Name: "club192983329"},
		// {Name: "club205417092"},
		// {Name: "club214317110"},
		// {Name: "gbdou6kr"},
		// {Name: "club42266729"},
		// {Name: "club76873688"},
		// {Name: "club202836702"},
		// {Name: "club202821332"},
		// {Name: "ds_65_krs_spb"},
		// {Name: "club205400739"},
		// {Name: "dou69krasnosel"},
		// {Name: "club216939970"},
		// {Name: "club205428969"},
		// {Name: "club205401911"},
		// {Name: "detskiy_sad74"},
		// {Name: "ds75spb"},
		// {Name: "club202011664"},
		// {Name: "club205406444"},
		// {Name: "ds78spb"},
		// {Name: "club129697643"},
		// {Name: "ds80krs"},
		// {Name: "club195029092"},
		// {Name: "club203610472"},
		// {Name: "gbdou83"},
		// {Name: "club203812364"},
		// {Name: "club205421015"},
		// {Name: "club202723926"},
		// {Name: "club215846431"},
		// {Name: "istokdetsad"},
		// {Name: "club194904593"},
		// {Name: "dc9spb"},
		// {Name: "children322029"},
		// {Name: "dou91krasnosel"},
		// {Name: "club205413257"},
		// {Name: "gbdou93krasnosel"},
		// {Name: "gbdou94"},
		// {Name: "gbdou95"},
		// https://vk.com/club227261708
		// https://vk.com/club183141138
		// https://vk.com/club205420830
		// https://vk.com/club193884037
		// https://vk.com/club214016041
		// https://vk.com/club200294876
		// https://vk.com/68rostok
		// https://vk.com/club205440005
		// https://vk.com/dc50krs_spb
		// https://vk.com/club180362982

		// https://vk.com/school509spb
		// https://vk.com/schoolspb54
		// https://vk.com/gym271
		// https://vk.com/gim293spb
		// https://vk.com/spb.school399
		// https://vk.com/club117133342
		// https://vk.com/public220312271
		// https://vk.com/licey_369
		// https://vk.com/licei395
		// https://vk.com/public__590
		// https://vk.com/club23933409
		// https://vk.com/club215520444
		// https://vk.com/school200spb
		// https://vk.com/rr_school208
		// https://vk.com/vr_odod_237
		// https://vk.com/sovet247
		// https://vk.com/school252spb
		// https://vk.com/spbschool262
		// https://vk.com/schooll270
		// https://vk.com/sch276spb
		// https://vk.com/school285spb
		// https://vk.com/g2343
		// https://vk.com/gbou291
		// https://vk.com/school352veteranov151
		// https://vk.com/school382spb
		// https://vk.com/school383
		// https://vk.com/club214266378
		// https://vk.com/spbschool390
		// https://vk.com/spbgboy391
		// https://vk.com/school394spb
		// https://vk.com/school414
		// https://vk.com/newschool546
		// https://vk.com/school547
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO vkgroup (name) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, vkgroup := range vkgroups {
		if _, err := stmt.Exec(vkgroup.Name); err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
