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
	db          *database.DbApi
	client      *xkcd.Client
	stemmer     *words.Stemmer
	max_workers int
}

func InitApp(db *database.DbApi, client *xkcd.Client, max_workers int) *App {
	stemmer := words.InitStemmer()

	return &App{
		db:          db,
		client:      client,
		stemmer:     stemmer,
		max_workers: max_workers,
	}
}

func (a *App) GetComics(ctx context.Context) {
	existing_comics := a.db.GetExisting()
	comics, err := a.client.GetComics(ctx, a.max_workers, existing_comics)
	if err != nil {
		log.Printf("Error getting comics: %s", err)
	}
	a.db.Insert(a.stem_comics(comics))
}

func (a *App) stem_comics(response_comics []models.ResponseComic) []models.DbComic {
	db_comics := make([]models.DbComic, len(response_comics))
	for i, comic := range response_comics {
		processing_text := comic.Alt + comic.Transcript
		key_words := a.stemmer.Stem(processing_text)

		db_comics[i] = models.DbComic{
			Id:       comic.Num,
			Url:      comic.Img,
			Keywords: key_words,
		}
	}

	return db_comics
}

func (a *App) PrintAll(ctx context.Context, n int) {
	max_id, err := a.client.GetLastId(ctx)
	if err != nil {
		log.Printf("Error getting last comic id: %s", err)
		return
	}
	if n > max_id {
		n = max_id
	}

	a.db.Print(n)
}
