package application

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	blevesearch "github.com/linealnan/glavredusgo/internal/blevesearch"
	"github.com/linealnan/glavredusgo/internal/db"
	"github.com/linealnan/glavredusgo/internal/telegrambot"
	vkindexer "github.com/linealnan/glavredusgo/internal/vkindexer"
	"github.com/robfig/cron/v3"
)

// Application представляет основное приложение
type Application struct {
	vi           *vkindexer.VkIndexer
	dbconn       *sql.DB
	bleaveSearch *blevesearch.BleaveSearch
	tb           *telegrambot.TelegramBot
}

type SearchHandler *TgSearchHandler
type TgSearchHandler interface {
	SearchAndFormat(bs blevesearch.BleaveSearch, searchPhrase string) string
}

func NewApplication(
	vi *vkindexer.VkIndexer,
	dbconn *sql.DB,
	bleaveSearch *blevesearch.BleaveSearch,
	tb *telegrambot.TelegramBot,
) *Application {
	return &Application{vi: vi, dbconn: dbconn, bleaveSearch: bleaveSearch, tb: tb}
}

func (a *Application) Run() {
	db.LoadSchema(a.dbconn)
	db.LoadSchoolVkGroups(a.dbconn)
	// a.vi.GetAndIndexedPosts()
	defer a.dbconn.Close()

	// Создаем новый экземпляр cron
	c := cron.New()

	// Каждые 15 минут
	c.AddFunc("*/15 * * * *", func() {
		fmt.Println("Загрузка данных групп по расписанию:", time.Now())
		a.vi.GetAndIndexedPosts()
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", searchHandler(*a.bleaveSearch))

	go a.tb.SubsribeUpdates()

	// Запускаем планировщик в фоновом режиме
	c.Start()
	// Запускаем сервер
	log.Printf("Starting server at port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Println("Error starting the server:", err)
	}

	select {}
}

func searchHandler(bs blevesearch.BleaveSearch) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			searchResult := bs.Search(name)

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
}
