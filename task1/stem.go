package main

import (
	"log"
	"os"
	"strings"

	"github.com/kljensen/snowball"
)

type Stemmer struct {
	stopWordsList []string
}

func InitStemmer() *Stemmer {
	/*
		Оригинальный список стоп-слов из библиотеки snowball мал, поэтому решил взять данный список.
		Взято отсюда https://countwordsfree.com/stopwords
	*/
	file, err := os.ReadFile("stop_words_english.txt")

	if err != nil {
		log.Fatalf("Error reading stop words file: %s", err.Error())
	}

	stopWordsList := strings.Fields(string(file))

	// // После Split на конце строк остался '\n'. Удаляем его.
	// for i, line := range stopWordsList {
	// 	stopWordsList[i] = strings.TrimSpace(line) // наверное это можно сделать проще
	// }
	return &Stemmer{stopWordsList: stopWordsList}
}

func (s *Stemmer) checkStopWord(target string) bool {
	for _, stopWord := range s.stopWordsList {
		if target == stopWord {
			return true
		}
	}
	return false
}

func (s *Stemmer) trimPunctuation(target string) string {
	return strings.Trim(target, ",.!?:;\"'()[]{}")
}

func (s *Stemmer) Stem(initialString string) string {
	words := strings.Fields(initialString)
	ans := []string{}
	seenWords := make(map[string]bool)

	for i := range words {

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
	return strings.Join(ans, " ")
}
