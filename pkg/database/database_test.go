package database

import (
	"testing"

	"github.com/MikhailFerapontow/yadro-go/pkg/words"
)

func BenchmarkFind(b *testing.B) {
	db := NewDbApi("database.json")
	stemmer := words.InitStemmer()
	weightedInput := stemmer.Stem("i'll follow you as long as you are following me")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Find(weightedInput)
	}
}
