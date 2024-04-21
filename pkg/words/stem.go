package words

import (
	_ "embed"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/MikhailFerapontow/yadro-go/models"
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
	return strings.Trim(target, ",.!?:;\"'()[]{}#<>*")
}

func (s *Stemmer) Stem(initialString string) []models.WeightedWord {
	words := strings.Fields(initialString)
	seenWords := make(map[string]int)

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	wg.Add(len(words))
	for i := range words {
		go func(i int) {
			defer wg.Done()

			if len(words[i]) < 4 {
				return
			}

			stemmed, err := snowball.Stem(s.trimPunctuation(words[i]), "english", false)

			if err != nil {
				log.Fatalf("Internal error stemming word: %s", err.Error())
			}

			// Проверяем является ли слово стоп-словом
			if s.stopWordMap[stemmed] {
				return
			}

			mu.Lock()
			seenWords[stemmed] += 1
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	cnt := make([]models.WeightedWord, len(seenWords))
	i := 0
	for k, v := range seenWords {
		cnt[i] = models.WeightedWord{Word: k, Count: v}
		i++
	}

	sort.Slice(cnt, func(i, j int) bool {
		return cnt[i].Count > cnt[j].Count
	})

	ans := make([]models.WeightedWord, 0)
	n := 10
	for i, w := range cnt {
		if i > n {
			break
		}
		ans = append(ans, w)
	}

	return ans
}
