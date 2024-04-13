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

// func (d *DbApi) Insert(comics []models.DbComic) {

// 	file, err := os.OpenFile(d.file_path, os.O_APPEND|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Printf("Error creating file: %s", err.Error())
// 		return
// 	}
// 	defer file.Close()

// 	bytes, _ := json.MarshalIndent(comics, "", " ")
// 	file.Write(bytes)
// }

func (d *DbApi) Insert(comics []models.DbComic) {
	text, err := os.ReadFile(d.file_path)
	if err != nil {
		log.Printf("Error opening file: %s", err.Error())
		return
	}

	// getting existing comics
	var db_comics []models.DbComic
	if err := json.Unmarshal(text, &db_comics); err != nil {
		log.Printf("Error unmarshaling json: %s", err.Error())
		return
	}

	comics = append(comics, db_comics...)
	file, err := os.Create(d.file_path)
	if err != nil {
		log.Printf("Error creating file: %s", err.Error())
		return
	}

	bytes, _ := json.MarshalIndent(comics, "", " ")
	file.Write(bytes)
}

func (d *DbApi) GetExisting() map[int]bool {
	file, err := os.ReadFile(d.file_path)
	if err != nil {
		log.Fatalf("Error opening file: %s", err.Error())
	}

	var db_comics []models.DbComic
	if err := json.Unmarshal(file, &db_comics); err != nil {
		log.Printf("Error unmarshaling json: %s", err.Error())
		return make(map[int]bool)
	}

	var existing_comics = make(map[int]bool)
	for _, comic := range db_comics {
		existing_comics[comic.Id] = true
	}

	return existing_comics
}

func (d *DbApi) Print(n int) {
	file, err := os.ReadFile(d.file_path)
	if err != nil {
		log.Printf("Error opening file: %s", err.Error())
		return
	}
	var db_comics []models.DbComic
	if err := json.Unmarshal(file, &db_comics); err != nil {
		log.Printf("Error unmarshaling json: %s", err.Error())
		return
	}

	for i := 0; i < n; i++ {
		pretty_json, err := json.MarshalIndent(db_comics[i], "", " ")
		if err != nil {
			log.Printf("Error marshaling json with id = %d: %s", db_comics[i].Id, err.Error())
		}
		fmt.Println(string(pretty_json))
	}
}
