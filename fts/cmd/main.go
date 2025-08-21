package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	index "github.com/linealnan/glavredusgo/fts/internal/index"
	vkapi "github.com/romanraspopov/golang-vk-api"
)

type MockGroup struct {
	Name string
}

// UserCity содержит id и название населенного пункта пользователя ВК
// Информация о городе, указанном на странице пользователя в разделе «Контакты».
// Возвращаются следующие поля:
// id (integer) — идентификатор города, который можно использовать для получения его названия с помощью метода database.getCitiesById;
// title (string) — название города.
type UserCity struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Full-Text Search (FTS)
// Raw Text -> tokenizer->filters->tokens
// https://habr.com/ru/articles/519024/
// https://github.com/akrylysov/simplefts
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("VK_API_TOKEN")
	client, err := vkapi.NewVKClientWithToken(token, nil, true)
	if err != nil {
		log.Fatal(err)
	}
	loadGroupsData(client)
}

func loadGroupsData(client *vkapi.VKClient) {
	groups := getGroups()
	log.Printf("Загрузка данных групп\n")
	for _, group := range groups {
		getAndIndexedWallPostByGroupName(client, group.Name)
	}
}

func getAndIndexedWallPostByGroupName(client *vkapi.VKClient, groupName string) {
	var documents []index.Document
	var document index.Document
	wall, err := client.WallGet(groupName, 100, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range wall.Posts {
		// log.Printf("Wall post: %v\n", post.Text)
		document.ID = post.ID
		document.Text = post.Text
		document.URL = "https://vk.com/trenchcrusade?w=wall-226198546_" + strconv.Itoa(post.ID)

		documents = append(documents, document)
	}

	// query := "охотник на ведьм"
	query := "убивать еретиков"
	// query := "косплей на ведьму"

	start := time.Now()
	idx := make(index.Index)
	idx.Add(documents)
	log.Printf("Indexed %d documents in %v", len(documents), time.Since(start))

	start = time.Now()
	matchedIDs := idx.Search(query)
	log.Printf("Search found %d documents in %v", len(matchedIDs), time.Since(start))

	for _, id := range matchedIDs {
		for _, doc := range documents {
			if doc.ID == id {
				log.Printf("%d\t%s\n", id, doc.URL)
			}
		}
	}
}

func getGroups() []MockGroup {
	return []MockGroup{
		{Name: "trenchcrusade"},
		// {Name: "nvp_73"},
		// {Name: "ad_ka4alka"},
	}
}
