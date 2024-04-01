package main

import (
	"flag"
	"fmt"

	"github.com/MikhailFerapontow/yadro-go/internal/config"
)

func main() {
	var print_output bool
	var comic_number int
	flag.BoolVar(&print_output, "o", false, "flag o prints result json into terminal")
	flag.IntVar(&comic_number, "n", 0, "flag n prints n-th comic")

	/*
		ничего плохого не произойдёт из-за паники в этой функции,
		ведь работа программы ещё не начата
	*/
	config.MustLoad()

	flag.Parse()

	// тоже самое что с конфигом
	if comic_number < 0 {
		panic("n must be >= 0")
	}

	if print_output {
		print_to_terminal()
	}
}

func print_to_terminal() {
	fmt.Println("Bruh")
}
