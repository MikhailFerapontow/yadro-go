package main

import (
	"flag"

	"github.com/MikhailFerapontow/yadro-go/internal/config"
	"github.com/MikhailFerapontow/yadro-go/pkg/database"
	"github.com/MikhailFerapontow/yadro-go/pkg/xkcd"
	"github.com/spf13/viper"
)

func main() {
	var print_output bool
	var comics_number int // используем zero value для получения всех комиксов
	flag.BoolVar(&print_output, "o", false, "flag o prints result json into terminal")
	flag.IntVar(&comics_number, "n", 0, "flag n prints n-th comic")

	/*
		ничего плохого не произойдёт из-за паники в этой функции,
		ведь работа программы ещё не начата
	*/
	config.MustLoad()

	flag.Parse()

	// тоже самое что с конфигом
	if comics_number < 0 {
		panic("n must be >= 0")
	}

	db := database.NewDbApi(viper.GetString("db_file"))
	client := xkcd.NewCLient(viper.GetString("source_url"), db)
	client.GetComics(comics_number)

	if print_output {
		client.PrintAllComics()
	}
}
