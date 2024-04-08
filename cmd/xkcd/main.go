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
	var config_path string
	flag.BoolVar(&print_output, "o", false, "flag -o prints result json into terminal")
	flag.IntVar(&comics_number, "n", math.MaxInt, "flag n prints up to n-th comic, WORKS ONLY WITH -o flag")
	flag.StringVar(&config_path, "c", ".", "path to config file. Name of config filemust be config.yaml")

	flag.Parse()

	config.MustLoad(config_path)

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
