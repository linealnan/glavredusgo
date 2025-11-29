package vkclient

import (
	conf "github.com/linealnan/glavredusgo/internal/config"
	vkapi "github.com/romanraspopov/golang-vk-api"
)

func NewVkClient(c *conf.AppConfig) (*vkapi.VKClient, error) {
	client, err := vkapi.NewVKClientWithToken(c.VkApiToken, nil, true)

	return client, err
}
