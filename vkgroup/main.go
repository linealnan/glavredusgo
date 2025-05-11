package main

import (
	"fmt"
	"log"

	vkapi "github.com/romanraspopov/golang-vk-api"
)

func main() {
	client, err := vkapi.NewVKClientWithToken("", nil, true)
	if err != nil {
		log.Fatal(err)
	}
	getAndSaveWallPost(client)
}

func getAndSaveWallPost(client *vkapi.VKClient) {
	wall, err := client.WallGet("trenchcrusade", 100, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range wall.Posts {
		fmt.Printf("Wall post: %v\n", post)
	}
}
