package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	vkapi "github.com/romanraspopov/golang-vk-api"
)

type LoadedPost struct {
	// Text string `xml:"abstract"`
	Text string
	ID   int
}

type MockGroup struct {
	Name string
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

func getGroups() []MockGroup {
	return []MockGroup{
		{Name: "trenchcrusade"},
		// {Name: "nvp_73"},
		// {Name: "ad_ka4alka"},
	}
}

func loadGroupsData(client *vkapi.VKClient) {
	groups := getGroups()
	log.Printf("Загрузка данных групп\n")
	for _, group := range groups {
		getAndIndexedWallPostByGroupName(client, group.Name)
	}
}

func getAndIndexedWallPostByGroupName(client *vkapi.VKClient, groupName string) {
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

		log.Printf("Wall post: %v\n", post.Text)
	}

}
