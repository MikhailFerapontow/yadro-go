package main

import (
	"flag"
	"fmt"
)

func main() {
	sFlag := flag.String("s", "", "flag s only takes strings!")

	flag.Parse()

	stemmer := InitStemmer() // инициализация стеммера

	ans := stemmer.Stem(*sFlag)
	fmt.Println(ans)
}
