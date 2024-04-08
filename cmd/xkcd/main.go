package main

import (
	"flag"
	"math"

	"github.com/MikhailFerapontow/yadro-go/internal/config"
	"github.com/MikhailFerapontow/yadro-go/pkg/app"
	"github.com/spf13/viper"
)

func main() {
	var print_output bool
	var comics_number int // я очень хотел использовать uint (но бесконечный каст типов)
	flag.BoolVar(&print_output, "o", false, "flag -o prints result json into terminal")
	flag.IntVar(&comics_number, "n", math.MaxInt, "flag n prints up to n-th comic, WORKS ONLY WITH -o flag")
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

	app := app.InitApp(app.Config{
		File_path: viper.GetString("db_file"),
		Url:       viper.GetString("source_url"),
	})

	app.GetComics()

	if print_output {
		app.PrintAll(comics_number)
	}
}
