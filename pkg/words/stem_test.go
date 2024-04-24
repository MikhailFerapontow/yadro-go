package words

import (
	"testing"

	"github.com/MikhailFerapontow/yadro-go/models"
	"github.com/stretchr/testify/assert"
)

func TestStemmer(t *testing.T) {
	stemmer := InitStemmer()

	tests := []struct {
		name          string
		initialString string
		expectedAns   []models.WeightedWord
	}{
		{
			name:          "String without stop words",
			initialString: "Follow rule",
			expectedAns: []models.WeightedWord{
				{Word: "follow", Count: 1},
				{Word: "rule", Count: 1},
			},
		},
		{
			name:          "String with stop words",
			initialString: "Follow rule mines",
			expectedAns: []models.WeightedWord{
				{Word: "follow", Count: 1},
				{Word: "rule", Count: 1},
			},
		},
		{
			name:          "String from task",
			initialString: "i'll follow you as long as you are following me",
			expectedAns: []models.WeightedWord{
				{Word: "long", Count: 1},
				{Word: "follow", Count: 2},
			},
		},
		{
			name:          "String with punctuation",
			initialString: "follower, follow followers!",
			expectedAns: []models.WeightedWord{
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
