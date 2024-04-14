package app

import (
	"context"
	"log"

	"github.com/MikhailFerapontow/yadro-go/models"
	"github.com/MikhailFerapontow/yadro-go/pkg/database"
	"github.com/MikhailFerapontow/yadro-go/pkg/words"
	"github.com/MikhailFerapontow/yadro-go/pkg/xkcd"
)

type App struct {
	db         *database.DbApi
	client     *xkcd.Client
	stemmer    *words.Stemmer
	maxWorkers int
}

func InitApp(db *database.DbApi, client *xkcd.Client, maxWorkers int) *App {
	stemmer := words.InitStemmer()

	return &App{
		db:         db,
		client:     client,
		stemmer:    stemmer,
		maxWorkers: maxWorkers,
	}
}

func (a *App) GetComics(ctx context.Context) {
	existingComics := a.db.GetExisting()
	comics, err := a.client.GetComics(ctx, a.maxWorkers, existingComics)
	if err != nil {
		log.Printf("Error getting comics: %s", err)
	}
	a.db.Insert(a.stem_comics(comics))
}

func (a *App) stem_comics(response_comics []models.ResponseComic) []models.DbComic {
	dbComics := make([]models.DbComic, len(response_comics))
	for i, comic := range response_comics {
		processingText := comic.Alt + comic.Transcript
		keyWords := a.stemmer.Stem(processingText)

		dbComics[i] = models.DbComic{
			Id:       comic.Num,
			Url:      comic.Img,
			Keywords: keyWords,
		}
	}

	return dbComics
}
