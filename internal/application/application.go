package application

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	blevesearch "github.com/linealnan/glavredusgo/internal/blevesearch"
	"github.com/linealnan/glavredusgo/internal/db"
	"github.com/linealnan/glavredusgo/internal/telegrambot"
	vkindexer "github.com/linealnan/glavredusgo/internal/vkindexer"
	"github.com/robfig/cron/v3"

	// Swagger UI статические файлы
	_ "github.com/swaggo/http-swagger"

	// Подключение пакета docs для регистрации swagger спецификации
	_ "github.com/linealnan/glavredusgo/internal/application/docs"
)

// Application представляет основное приложение
type Application struct {
	vi           *vkindexer.VkIndexer
	dbconn       *sql.DB
	bleaveSearch *blevesearch.BleaveSearch
	tb           *telegrambot.TelegramBot
}

// SearchRequest представляет запрос на поиск
// @name SearchRequest
type SearchRequest struct {
	Query string `json:"query"`
}

// SearchResponse представляет ответ на запрос поиска
// @name SearchResponse
type SearchResponse struct {
	Total   uint64         `json:"total"`
	Results []SearchResult `json:"results"`
	Error   string         `json:"error,omitempty"`
}

// SearchResult представляет отдельный результат поиска
// @name SearchResult
type SearchResult struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	GroupName string `json:"group_name"`
	Text      string `json:"text"`
	URL       string `json:"url"`
}

type SearchHandler *TgSearchHandler
type TgSearchHandler interface {
	SearchAndFormat(bs blevesearch.BleaveSearch, searchPhrase string) string
}

func NewApplication(
	vi *vkindexer.VkIndexer,
	dbconn *sql.DB,
	bleaveSearch *blevesearch.BleaveSearch,
	// tb *telegrambot.TelegramBot,
) *Application {
	return &Application{vi: vi, dbconn: dbconn, bleaveSearch: bleaveSearch}
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

	// Swagger - обслуживаем JSON и статические файлы (регистрируем ПЕРЕД корневым маршрутом)

	// Обслуживаем swagger.json - встроенный в бинарник
	mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(swaggerJSON))
	})

	// Swagger UI - собственная HTML страница с CDN (должен быть ПЕРЕД "/" )
	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		// Если запрашивают index.html или корень /swagger/
		if r.URL.Path == "/swagger/" || r.URL.Path == "/swagger/index.html" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(swaggerHTML))
			return
		}
		// Остальные запросы - 404
		http.NotFound(w, r)
	})

	// API маршруты
	mux.HandleFunc("/api/search", jsonSearchHandler(*a.bleaveSearch))

	// Корневой маршрут (последний)
	mux.HandleFunc("/", searchHandler(*a.bleaveSearch))

	// заблокировано в России
	// go a.tb.SubsribeUpdates()

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

// @Summary Поиск постов
// @Description Выполняет поиск по индексированным постам ВКонтакте
// @Tags search
// @Accept json
// @Produce json
// @Param searchRequest body SearchRequest true "Запрос поиска"
// @Success 200 {object} SearchResponse
// @Failure 400 {object} SearchResponse
// @Failure 405 {object} SearchResponse
// @Router /api/search [post]
func jsonSearchHandler(bs blevesearch.BleaveSearch) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовки для JSON
		w.Header().Set("Content-Type", "application/json")

		// Проверяем метод запроса
		if r.Method != http.MethodPost {
			resp := SearchResponse{
				Error: "Метод не разрешен. Используйте POST",
			}
			json.NewEncoder(w).Encode(resp)
			log.Println("Метод не разрешен. Используйте POST")
			return
		}

		// Декодируем JSON запрос
		var req SearchRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			resp := SearchResponse{
				Error: "Ошибка декодирования запроса: " + err.Error(),
			}
			json.NewEncoder(w).Encode(resp)
			log.Println("Ошибка декодирования запроса: " + err.Error())
			return
		}

		// Проверяем наличие запроса
		if req.Query == "" {
			resp := SearchResponse{
				Error: "Поле query обязательно",
			}
			json.NewEncoder(w).Encode(resp)
			log.Println("Поле query обязательно")
			return
		}
		fmt.Println(req.Query, time.Now())
		// Выполняем поиск
		searchResult := bs.Search(req.Query)

		// Формируем результаты
		results := make([]SearchResult, 0, len(searchResult.Hits))
		for _, hit := range searchResult.Hits {
			fields := hit.Fields
			groupID, _ := fields["GroupID"].(string)
			groupName, _ := fields["GroupName"].(string)
			id, _ := fields["ID"].(string)
			text, _ := fields["Text"].(string)

			results = append(results, SearchResult{
				ID:        id,
				GroupID:   groupID,
				GroupName: groupName,
				Text:      text,
				URL:       fmt.Sprintf("https://vk.com/%s?w=wall-%s_%s", groupName, groupID, id),
			})
		}

		// Возвращаем JSON ответ
		resp := SearchResponse{
			Total:   searchResult.Total,
			Results: results,
		}
		json.NewEncoder(w).Encode(resp)
	}
}

// swaggerJSON - встроенная спецификация OpenAPI (Swagger 3.x)
const swaggerJSON = `{
    "openapi": "3.0.0",
    "info": {
        "title": "Glavredusgo API",
        "description": "API для поиска постов в группах ВКонтакте",
        "version": "1.0"
    },
    "servers": [{"url": "http://localhost:8080"}],
    "paths": {
        "/api/search": {
            "post": {
                "tags": ["search"],
                "summary": "Поиск постов",
                "description": "Выполняет поиск по индексированным постам ВКонтакте",
                "requestBody": {
                    "required": true,
                    "content": {
                        "application/json": {
                            "schema": {
                                "$ref": "#/components/schemas/SearchRequest"
                            }
                        }
                    }
                },
                "responses": {
                    "200": {
                        "description": "Успешный ответ",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/SearchResponse"
                                }
                            }
                        }
                    }
                }
            }
        }
    },
    "components": {
        "schemas": {
            "SearchRequest": {
                "type": "object",
                "properties": {"query": {"type": "string"}}
            },
            "SearchResponse": {
                "type": "object",
                "properties": {
                    "error": {"type": "string"},
                    "results": {"type": "array", "items": {"$ref": "#/components/schemas/SearchResult"}},
                    "total": {"type": "integer"}
                }
            },
            "SearchResult": {
                "type": "object",
                "properties": {
                    "id": {"type": "string"},
                    "group_id": {"type": "string"},
                    "group_name": {"type": "string"},
                    "text": {"type": "string"},
                    "url": {"type": "string"}
                }
            }
        }
    }
}`

// swaggerHTML - HTML страница для Swagger UI с CDN (версия 3.x)
const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Glavredusgo API - Swagger UI</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.43.0/swagger-ui.css" />
    <style>
        body { margin: 0; padding: 0; }
        .swagger-ui .info .title { font-size: 30px; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.43.0/swagger-ui-bundle.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.43.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/swagger/doc.json",
                dom_id: "swagger-ui",
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                layout: "StandaloneLayout"
            });
            window.ui = ui;
        };
    </script>
</body>
</html>`
