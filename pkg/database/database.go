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
	file, err := os.Create(d.file_path)
	if err != nil {
		log.Printf("Error creating file: %s", err.Error())
		return
	}

	defer file.Close()

	bytes, _ := json.MarshalIndent(comics, "", " ")
	file.Write(bytes)
}

func (d *DbApi) PrintAll(n int) {
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
