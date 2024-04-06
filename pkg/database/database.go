package database

import (
	"encoding/json"
	"log"
	"os"

	"github.com/MikhailFerapontow/yadro-go/models"
	"github.com/MikhailFerapontow/yadro-go/pkg/words"
)

type DbApi struct {
	file_path string
	stemmer   *words.Stemmer
}

func NewDbApi(file_path string) *DbApi {
	stemmer := words.InitStemmer()

	return &DbApi{
		file_path: file_path,
		stemmer:   stemmer,
	}
}

func (d *DbApi) Insert(comics []models.ResponseComics) {
	db_comics := make([]models.DbComic, len(comics))
	for i, comic := range comics {
		db_comics[i] = d.process(comic)
	}

	file, err := os.Create(d.file_path)
	if err != nil {
		log.Printf("Error creating file: %s", err.Error())
		return
	}

	defer file.Close()

	bytes, _ := json.MarshalIndent(db_comics, "", " ")
	file.Write(bytes)
}

func (d *DbApi) PrintAll() {
	file, err := os.ReadFile(d.file_path)
	if err != nil {
		log.Printf("Error opening file: %s", err.Error())
		return
	}
	os.Stdout.Write(file)
}

func (d *DbApi) process(comic models.ResponseComics) models.DbComic {
	processing_text := comic.Alt + comic.Transcript

	key_words := d.stemmer.Stem(processing_text)

	return models.DbComic{
		Id:       comic.Num,
		Url:      comic.Img,
		Keywords: key_words,
	}
}
