package vkindexer

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/linealnan/glavredusgo/internal/blevesearch"
	conf "github.com/linealnan/glavredusgo/internal/config"
	vkapi "github.com/romanraspopov/golang-vk-api"
)

type VkIndexer struct {
	conf         *conf.AppConfig
	bleaveSearch *blevesearch.BleaveSearch
	client       *vkapi.VKClient
	db           *sql.DB
}

type LoadedPost struct {
	ID        string
	Text      string
	GroupName string
	GroupID   string
}

type VkGroup struct {
	Name string
}

func NewVkIndexer(c *conf.AppConfig, bleaveSearch *blevesearch.BleaveSearch, client *vkapi.VKClient, db *sql.DB) *VkIndexer {
	return &VkIndexer{conf: c, bleaveSearch: bleaveSearch, client: client, db: db}
}

func (vi *VkIndexer) GetAndIndexedPosts() error {
	groups := getVkGroups(vi.db)
	log.Printf("Получение данных групп\n")
	for _, group := range groups {
		getAndIndexedWallPostByGroupName(vi.client, group.Name, vi.bleaveSearch.Index)
	}

	return nil
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

func getAndIndexedWallPostByGroupName(client *vkapi.VKClient, groupName string, index bleve.Index) {
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
}
