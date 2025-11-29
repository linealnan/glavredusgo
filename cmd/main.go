package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/linealnan/glavredusgo/internal/application"
	bleveindex "github.com/linealnan/glavredusgo/internal/bleveindex"
	conf "github.com/linealnan/glavredusgo/internal/config"
	db "github.com/linealnan/glavredusgo/internal/db"
	"github.com/linealnan/glavredusgo/internal/vkclient"
	"github.com/linealnan/glavredusgo/internal/vkindexer"
	_ "github.com/mattn/go-sqlite3"
	vkapi "github.com/romanraspopov/golang-vk-api"
	"github.com/urfave/cli/v2"
	"go.uber.org/dig"
)

type IndexData struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	GroupName string `json:"groupName"`
	GroupID   string `json:"GroupID"`
}

type LoadedPost struct {
	ID        string
	Text      string
	GroupName string
	GroupID   string
}

type MockGroup struct {
	Name string
}
type VkGroup struct {
	Name string
}

var index bleve.Index

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// db, err := sql.Open("sqlite3", "glavredus.db")
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// deleteVkGroupTable(db)
	// loadInitSchema(db)
	// loadSchoolVkGroups(db)

	container := dig.New()

	if err := container.Provide(conf.InitWithDotEnv); err != nil {
		panic(err)
	}

	if err := container.Provide(bleveindex.NewBleveIndex); err != nil {
		panic(err)
	}

	if err := container.Provide(vkclient.NewVkClient); err != nil {
		panic(err)
	}

	if err := container.Provide(db.NewDbConnection); err != nil {
		panic(err)
	}

	if err := container.Provide(vkindexer.NewVkIndexer); err != nil {
		panic(err)
	}

	if err := container.Provide(application.NewApplication); err != nil {
		panic(err)
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
			// if c.String("i") == "force" {

			// 	token := os.Getenv("VK_API_TOKEN")
			// 	client, err := vkapi.NewVKClientWithToken(token, nil, true)
			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	log.Printf("Загрузка данных групп\n")
			// 	loadGroupsData(client, db)

			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	log.Printf("Индекс обновлен\n")

			// }

			// 3. Запускаем основную функцию приложения с помощью метода Invoke.
			// Dig автоматически разрешает зависимости и вызывает переданную функцию с готовыми экземплярами.
			if err := container.Invoke(func(app *application.Application) {
				app.Run()
			}); err != nil {
				panic(err)
			}

			http.HandleFunc("/", formHandler)

			// Запускаем сервер
			log.Printf("Starting server at port 8080")
			err := http.ListenAndServe(":8080", nil)
			if err != nil {
				log.Println("Error starting the server:", err)
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
	// rq := bleve.NewRegexpQuery(query)

	// qsq := bleve.NewQueryStringQuery(query)

	q := bleve.NewDisjunctionQuery(mq)

	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Highlight = bleve.NewHighlight()
	searchRequest.Fields = []string{"ID", "GroupName", "GroupID", "Text"}
	searchRequest.Size = 100

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
		searchResult := search(name)

		fmt.Printf("Найдено документов: %d\n", searchResult.Total)
		for _, hit := range searchResult.Hits {
			fields := hit.Fields
			// fmt.Printf("https://vk.com/%s?w=wall-%s_%s\n", fields["GroupName"], fields["GroupID"], fields["ID"])
			// fmt.Printf("Текст: %s", fields["Text"])
			// fmt.Fprintf(w, "Текст: %s", fields["Text"])
			fmt.Fprintf(w, "https://vk.com/%s?w=wall-%s_%s\n", fields["GroupName"], fields["GroupID"], fields["ID"])
		}

		// for _, item := range res.Hits {
		// 	result := item.Fields
		// 	// log.Printf("Результат поиска: %s", result)
		// 	fmt.Fprintf(w, "Результат поиска: %d%s", result["ID"], result["Text"])
		// }
		//log.Printf("Результат поиска: %s", result)
		fmt.Fprintf(w, "Результат поиска: %s", searchResult)

	}
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

	log.Printf("Получение постов группы %s\n", groupName)
	wall, err := client.WallGet(groupName, 100, nil)
	if err != nil {
		log.Fatal(err)
	}

	groupsSlice := []string{groupName}

	groups, err := client.GroupsGetByID(groupsSlice)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range wall.Posts {
		indexedPost.ID = strconv.Itoa(post.ID)
		indexedPost.Text = post.Text
		indexedPost.GroupName = groupName
		indexedPost.GroupID = strconv.Itoa(groups[0].ID)

		err = index.Index(strconv.Itoa(post.ID), indexedPost)
		if err != nil {
			log.Fatal(err)
		}
	}

	return posts
}

func deleteVkGroupTable(db *sql.DB) {
	sql := `DROP TABLE IF EXISTS vkgroup;`

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
		{"club205400972"},
		{"club182072023"},
		{"club195576991"},
		{"club147892228"},
		{"club187951249"},
		{"sadik31krs"},
		{"club205420428"},
		{"gdboy35"},
		{"club205443755"},
		{"dc39spb"},
		{"club170186955"},
		{"gbdou41krspb"},
		{"club216246675"},
		{"club205406349"},
		{"club203026295"},
		{"dc5krs"},
		{"ds51krs"},
		{"gbdouds52"},
		{"club13309436"},
		{"club192983329"},
		{"club205417092"},
		{"club214317110"},
		{"gbdou6kr"},
		{"club42266729"},
		{"club76873688"},
		{"club202836702"},
		{"club202821332"},
		{"ds_65_krs_spb"},
		{"club205400739"},
		{"dou69krasnosel"},
		{"club216939970"},
		{"club205428969"},
		{"club205401911"},
		{"detskiy_sad74"},
		{"ds75spb"},
		{"club202011664"},
		{"club205406444"},
		{"ds78spb"},
		{"club129697643"},
		{"ds80krs"},
		{"club195029092"},
		{"club203610472"},
		{"gbdou83"},
		{"club203812364"},
		{"club205421015"},
		{"club202723926"},
		{"club215846431"},
		// {"istokdetsad"},
		// {"club194904593"},
		{"dc9spb"},
		{"children322029"},
		{"dou91krasnosel"},
		{"club205413257"},
		{"gbdou93krasnosel"},
		{"gbdou94"},
		{"club227261708"},
		{"club183141138"},
		{"club205420830"},
		{"club193884037"},
		{"club214016041"},
		{"club200294876"},
		{"68rostok"},
		{"club205440005"},
		{"dc50krs_spb"},
		{"club180362982"},
		{"school509spb"},
		{"schoolspb54"},
		{"gym271"},
		{"gim293spb"},
		{"spb.school399"},
		{"club117133342"},
		{"public220312271"},
		{"licey_369"},
		{"licei395"},
		{"public__590"},
		{"club23933409"},
		{"club215520444"},
		{"school200spb"},
		{"rr_school208"},
		{"vr_odod_237"},
		{"sovet247"},
		{"school252spb"},
		{"spbschool262"},
		{"schooll270"},
		{"sch276spb"},
		{"school285spb"},
		{"g2343"},
		{"gbou291"},
		{"school352veteranov151"},
		{"school382spb"},
		{"school383"},
		{"club214266378"},
		{"spbschool390"},
		{"spbgboy391"},
		{"school394spb"},
		{"school414"},
		{"newschool546"},
		{"school547"},
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
