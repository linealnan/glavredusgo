package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	vkapi "github.com/romanraspopov/golang-vk-api"
)

// type Group struct {
// 	ID                int             `json:"id"`
// 	Name              string          `json:"name"`
// 	ScreenName        string          `json:"screen_name"`
// 	Description       string          `json:"description"`
// 	Activity          string          `json:"activity"`
// 	Contacts          []*GroupContact `json:"contacts"`
// 	IsClosed          int             `json:"is_closed"`
// 	Type              string          `json:"type"`
// 	IsAdmin           int             `json:"is_admin"`
// 	IsMember          int             `json:"is_member"`
// 	MembersCount      int             `json:"members_count"`
// 	HasPhoto          int             `json:"has_photo"`
// 	IsMessagesBlocked int             `json:"is_messages_blocked"`
// 	Photo50           string          `json:"photo_50"`
// 	Photo100          string          `json:"photo_100"`
// 	Photo200          string          `json:"photo_200"`
// 	AgeLimit          int             `json:"age_limits"`
// 	CanCreateTopic    int             `json:"can_create_topic"`
// 	CanMessage        int             `json:"can_message"`
// 	CanPost           int             `json:"can_post"`
// 	CanSeeAllPosts    int             `json:"can_see_all_posts"`
// 	City              *UserCity       `json:"city"`
// }

type MockGroup struct {
	Name string
}

// type GroupContact struct {
// 	UID         int    `json:"user_id"`
// 	Description string `json:"desc"`
// }

// UserCity содержит id и название населенного пункта пользователя ВК
// Информация о городе, указанном на странице пользователя в разделе «Контакты».
// Возвращаются следующие поля:
// id (integer) — идентификатор города, который можно использовать для получения его названия с помощью метода database.getCitiesById;
// title (string) — название города.
type UserCity struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

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
	for _, group := range groups {
		getAndSaveWallPostByGroupName(client, group.Name)
	}
}

func getAndSaveWallPostByGroupName(client *vkapi.VKClient, groupName string) {
	wall, err := client.WallGet(groupName, 100, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range wall.Posts {
		fmt.Printf("Wall post: %v\n", post)
	}
}

func getGroups() []MockGroup {
	return []MockGroup{
		{Name: "trenchcrusade"},
		{Name: "nvp_73"},
		{Name: "ad_ka4alka"},
	}
}
