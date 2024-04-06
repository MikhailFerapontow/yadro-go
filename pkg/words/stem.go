package words

import (
	"log"
	"os"
	"strings"

	"github.com/kljensen/snowball"
)

type Stemmer struct {
	stopWordMap map[string]any
}

func InitStemmer() *Stemmer {
	/*
		Оригинальный список стоп-слов из библиотеки snowball мал, поэтому решил взять данный список.
		Взято отсюда https://countwordsfree.com/stopwords
	*/
	file, err := os.ReadFile("pkg/words/stop_words_english.txt")

	if err != nil {
		log.Fatalf("Error reading stop words file: %s", err.Error())
	}

	var stopWords = make(map[string]any)

	stopWordsList := strings.Fields(string(file))
	for i, elem := range stopWordsList {
		stopWords[elem] = i
	}

	return &Stemmer{stopWordMap: stopWords}
}

func (s *Stemmer) checkStopWord(target string) bool {
	_, ok := s.stopWordMap[target]
	return ok
}

func (s *Stemmer) trimPunctuation(target string) string {
	return strings.Trim(target, ",.!?:;\"'()[]{}#<>")
}

func (s *Stemmer) Stem(initialString string) []string {
	words := strings.Fields(initialString)
	ans := []string{}
	seenWords := make(map[string]bool)

	for i := range words {
		if len(words[i]) < 4 {
			continue
		}

		stemmed, err := snowball.Stem(s.trimPunctuation(words[i]), "english", false)

		if err != nil {
			log.Fatalf("Internal error stemming word: %s", err.Error())
		}

		// Проверяем является ли слово стоп-словом
		if s.checkStopWord(stemmed) {
			continue
		}

		// Проверям если уже встречали это слово
		if seenWords[stemmed] {
			continue
		}

		seenWords[stemmed] = true
		ans = append(ans, stemmed)
	}
	return ans
}
