package database

import (
	"testing"

	"github.com/MikhailFerapontow/yadro-go/pkg/words"
)

var table = []struct {
	testName    string
	searchInput string
}{
	{testName: "Empty search input", searchInput: ""},
	{testName: "Short search input", searchInput: "word"},
	{testName: "Medium search input", searchInput: "I'll follow you as long as I can"},
	{testName: "Very long search input",
		searchInput: `Even at the hour when the grey St. Petersburg sky had quite
	dispersed, and all the official world had eaten or dined, each as he could,
	in accordance with the salary he received and his own fancy; when all were
	resting from the departmental jar of pens, running to and fro from their
	own and other people’s indispensable occupations, and from all the work that
	an uneasy man makes willingly for himself, rather than what is necessary;
	when officials hasten to dedicate to pleasure the time which is left to them,
	one bolder than the rest going to the theatre; another, into the street
	looking under all the bonnets; another wasting his evening in compliments to
	some pretty girl, the star of a small official circle; another — and this is
	the common case of all — visiting his comrades on the fourth or third floor,
	in two small rooms with an ante-room or kitchen, and some pretensions to fashion,
	such as a lamp or some other trifle which has cost many a sacrifice of dinner or
	pleasure trip; in a word, at the hour when all officials disperse among the contracted
	quarters of their friends, to play whist, as they sip their tea from glasses with a
	kopek’s worth of sugar, smoke long pipes, relate at times some bits of gossip
	which a Russian man can never, under any circumstances, refrain from, and,
	when there is nothing else to talk of, repeat eternal anecdotes about the
	commandant to whom they had sent word that the tails of the horses on the
	Falconet Monument had been cut off, when all strive to divert themselves,
	Akakiy Akakievitch indulged in no kind of diversion.`},
}

func BenchmarkFind(b *testing.B) {
	db := NewDbApi("database.json")
	stemmer := words.InitStemmer()

	b.ResetTimer()
	for _, testCase := range table {
		weightedInput := stemmer.Stem(testCase.searchInput)
		b.Run(testCase.testName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				db.Find(weightedInput)
			}
		})
	}
}

func BenchmarkFindByIndex(b *testing.B) {
	db := NewDbApi("database.json")
	stemmer := words.InitStemmer()

	b.ResetTimer()
	for _, testCase := range table {
		weightedInput := stemmer.Stem(testCase.searchInput)
		b.Run(testCase.testName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				db.FindByIndex(weightedInput)
			}
		})
	}
}
