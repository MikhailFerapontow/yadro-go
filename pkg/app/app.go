package app

import (
	"log"

	"github.com/MikhailFerapontow/yadro-go/models"
	"github.com/MikhailFerapontow/yadro-go/pkg/database"
	"github.com/MikhailFerapontow/yadro-go/pkg/words"
	"github.com/MikhailFerapontow/yadro-go/pkg/xkcd"
)

type Config struct {
	File_path string
	Url       string
}

type App struct {
	db      *database.DbApi
	client  *xkcd.Client
	stemmer *words.Stemmer
}

func InitApp(cfg Config) *App {
	db := database.NewDbApi(cfg.File_path)
	client := xkcd.NewCLient(cfg.Url)
	stemmer := words.InitStemmer()

	return &App{
		db:      db,
		client:  client,
		stemmer: stemmer,
	}
}

func (a *App) GetComics() {
	comics, err := a.client.GetComics()
	if err != nil {
		log.Printf("Error getting comics: %s", err)
		return
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

func (a *App) PrintAll(n int) {
	max_id, err := a.client.GetLastComicId()
	if err != nil {
		log.Printf("Error getting last comic id: %s", err)
		return
	}
	if n > max_id {
		n = max_id
	}

	a.db.PrintAll(n)
}
