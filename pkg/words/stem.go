package words

import (
	_ "embed"
	"log"
	"strings"

	"github.com/kljensen/snowball"
)

/*
	Оригинальный список стоп-слов из библиотеки snowball мал, поэтому решил взять данный список.
	Взято отсюда https://countwordsfree.com/stopwords
*/
//go:embed stop_words_english.txt
var stopWordsEnglish string

type Stemmer struct {
	stopWordMap map[string]bool
}

func InitStemmer() *Stemmer {
	var stopWords = make(map[string]bool)

	stopWordsList := strings.Fields(stopWordsEnglish)
	for _, elem := range stopWordsList {
		stopWords[elem] = true
	}

	return &Stemmer{stopWordMap: stopWords}
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
		if s.stopWordMap[stemmed] {
			continue
		}

		// Проверям если уже встречали это слово
		if seenWords[stemmed] {
			continue
		}

		seenWords[stemmed] = true
		ans = append(ans, stemmed)
	}

	/*
		По хорошему нужно провести статистический анализ какие слова встречаются чаще всего
		пока так сойдёт
	*/
	if len(ans) > 10 {
		ans = ans[:10]
	}
	return ans
}
