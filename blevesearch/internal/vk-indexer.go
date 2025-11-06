package internal

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index"
	vkapi "github.com/romanraspopov/golang-vk-api"
)

type VkIndexer struct {
	BaseService
	services   map[string]Service
	config     *AppConfig
	bleveIndex bleve.Index
	client     *vkapi.VKClient
	db         *sql.DB
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

func (vi *VkIndexer) Init(conf *AppConfig, db *sql.DB) error {
	token := conf.VkApiToken
	client, err := vkapi.NewVKClientWithToken(token, nil, true)
	if err != nil {
		log.Fatal(err)
	}
	vi.client = client
	vi.db = db

	return nil
}

func (vi *VkIndexer) Run() error {
	groups := getVkGroups(vi.db)
	log.Printf("Получение данных групп\n")
	for _, group := range groups {
		getAndIndexedWallPostByGroupName(vi.client, group.Name)
	}

	return nil
}

func (vi *VkIndexer) Name() string {
	return "vk-indexer"
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

func getAndIndexedWallPostByGroupName(client *vkapi.VKClient, groupName string) {
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
