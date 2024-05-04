package repository

import (
	"testing"

	"github.com/MikhailFerapontow/yadro-go/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestStemmer(t *testing.T) {
	stemmer := InitStemmer()

	tests := []struct {
		name          string
		initialString string
		expectedAns   []domain.WeightedWord
	}{
		{
			name:          "String without stop words",
			initialString: "Follow rule",
			expectedAns: []domain.WeightedWord{
				{Word: "follow", Count: 1},
				{Word: "rule", Count: 1},
			},
		},
		{
			name:          "String with stop words",
			initialString: "Follow rule mines",
			expectedAns: []domain.WeightedWord{
				{Word: "follow", Count: 1},
				{Word: "rule", Count: 1},
			},
		},
		{
			name:          "String from task",
			initialString: "i'll follow you as long as you are following me",
			expectedAns: []domain.WeightedWord{
				{Word: "long", Count: 1},
				{Word: "follow", Count: 2},
			},
		},
		{
			name:          "String with punctuation",
			initialString: "follower, follow followers!",
			expectedAns: []domain.WeightedWord{
				{Word: "follow", Count: 3},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ans := stemmer.Stem(test.initialString)
			for _, word := range ans {
				assert.Contains(t, test.expectedAns, word)
			}
		})
	}
}
