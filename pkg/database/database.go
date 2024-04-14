package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/MikhailFerapontow/yadro-go/models"
)

type DbApi struct {
	file_path string
}

func NewDbApi(file_path string) *DbApi {
	return &DbApi{
		file_path: file_path,
	}
}

func (d *DbApi) Insert(comics []models.DbComic) {
	op := "op.insert"

	text, err := os.ReadFile(d.file_path)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return
	}

	var dbComics []models.DbComic
	if err := json.Unmarshal(text, &dbComics); err != nil {
		log.Printf("%s: Error unmarshaling json, file empty or with errors: %s", op, err)
	}

	comics = append(comics, dbComics...)
	file, err := os.Create(d.file_path)
	if err != nil {
		log.Printf("Error creating file: %s", err)
		return
	}
	defer file.Close()

	bytes, _ := json.MarshalIndent(comics, "", " ")
	file.Write(bytes)
	log.Printf("%s: Successfully inserted comics", op)
}

func (d *DbApi) GetExisting() map[int]bool {
	op := "op.get_existing_comics"

	text, err := os.ReadFile(d.file_path)
	if err != nil {
		log.Printf("%s: Error opening file: %s", op, err)
		log.Printf("Creating file: %s", d.file_path)
		os.Create(d.file_path)
		return make(map[int]bool)
	}

	var dbComics []models.DbComic
	if err := json.Unmarshal(text, &dbComics); err != nil {
		log.Printf("%s: Error unmarshaling json, file empty or with errors: %s", op, err)
		return make(map[int]bool)
	}

	existingComics := make(map[int]bool)
	for _, comic := range dbComics {
		existingComics[comic.Id] = true
	}

	return existingComics
}

func (d *DbApi) Print(n int) {
	file, err := os.ReadFile(d.file_path)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return
	}
	var db_comics []models.DbComic
	if err := json.Unmarshal(file, &db_comics); err != nil {
		log.Printf("Error unmarshaling json: %s", err)
		return
	}

	for i := 0; i < n; i++ {
		pretty_json, err := json.MarshalIndent(db_comics[i], "", " ")
		if err != nil {
			log.Printf("Error marshaling json with id = %d: %s", db_comics[i].Id, err)
		}
		fmt.Println(string(pretty_json))
	}
}
