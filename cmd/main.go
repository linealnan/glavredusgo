package main

import (
	"log"
	"os"

	"github.com/linealnan/glavredusgo/internal/application"
	blevesearch "github.com/linealnan/glavredusgo/internal/blevesearch"
	conf "github.com/linealnan/glavredusgo/internal/config"
	db "github.com/linealnan/glavredusgo/internal/db"
	"github.com/linealnan/glavredusgo/internal/telegrambot"
	"github.com/linealnan/glavredusgo/internal/vkclient"
	"github.com/linealnan/glavredusgo/internal/vkindexer"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
	"go.uber.org/dig"
)

func main() {
	container := dig.New()
	// TODO Переписать бороду условий
	if err := container.Provide(conf.InitWithDotEnv); err != nil {
		panic(err)
	}

	if err := container.Provide(blevesearch.NewBleaveSearch); err != nil {
		panic(err)
	}

	if err := container.Provide(vkclient.NewVkClient); err != nil {
		panic(err)
	}

	if err := container.Provide(db.NewDbConnection); err != nil {
		panic(err)
	}

	if err := container.Provide(vkindexer.NewVkIndexer); err != nil {
		panic(err)
	}

	if err := container.Provide(telegrambot.NewTelegramBot); err != nil {
		panic(err)
	}

	if err := container.Provide(application.NewApplication); err != nil {
		panic(err)
	}

	app := &cli.App{
		Name:  "glavredus",
		Usage: "Поиск по группам",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "index",
				Aliases: []string{"i"},
				Value:   "",
				Usage:   "загрузить данные групп в индекс",
			},
		},
		Action: func(c *cli.Context) error {
			if err := container.Invoke(func(app *application.Application) {
				app.Run()

			}); err != nil {
				panic(err)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
