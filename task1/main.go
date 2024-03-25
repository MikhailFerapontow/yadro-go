package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	// read string from CLI
	// var input string
	sFlag := flag.String("s", "", "flag s only takes strings!")

	flag.Parse()

	fmt.Println(*sFlag) // распарсили аргументы CLI, получили нашу строку.
	words := strings.Fields(*sFlag)
	fmt.Println(words)
	for i := range words {
		fmt.Println(words[i])
	}
}
