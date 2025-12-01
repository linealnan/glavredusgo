package application

import (
	"database/sql"

	"github.com/linealnan/glavredusgo/internal/db"
	vkindexer "github.com/linealnan/glavredusgo/internal/vkindexer"
)

// Application представляет основное приложение
type Application struct {
	vi     *vkindexer.VkIndexer
	dbconn *sql.DB
}

func NewApplication(vi *vkindexer.VkIndexer, dbconn *sql.DB) *Application {
	return &Application{vi: vi, dbconn: dbconn}
}

func (a *Application) Run() {
	a.vi.GetAndIndexedPosts()
	db.LoadSchema(a.dbconn)
	db.LoadSchoolVkGroups(a.dbconn)

	defer a.dbconn.Close()
}
