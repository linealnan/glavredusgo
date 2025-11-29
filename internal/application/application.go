package application

import (
	vkindexer "github.com/linealnan/glavredusgo/internal/vkindexer"
)

// Application представляет основное приложение
type Application struct {
	vi *vkindexer.VkIndexer
}

// NewApplication — это конструктор для Application, он зависит от conf.AppConfig
func NewApplication(vi *vkindexer.VkIndexer) *Application {
	return &Application{vi: vi}
}

// Run — метод запуска приложения, который использует логгер
func (a *Application) Run() {
	a.vi.GetAndIndexedPosts()
}
