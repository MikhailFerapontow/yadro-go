package words

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStemmer(t *testing.T) {
	stemmer := InitStemmer()

	tests := []struct {
		name          string
		initialString string
		expectedAns   string
	}{
		{
			name:          "String without stop words",
			initialString: "Follow rule",
			expectedAns:   "follow rule",
		},
		{
			name:          "String with stop words",
			initialString: "follower brings bunch of questions",
			expectedAns:   "follow bring bunch question",
		},
		{
			name:          "String with punctuation",
			initialString: "follower, follow followers!",
			expectedAns:   "follow",
		},
		{
			name:          "String from task",
			initialString: "i'll follow you as long as you are following me",
			expectedAns:   "follow long",
		},
		{
			name:          "String with all stop words",
			initialString: "me you I",
			expectedAns:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ans := stemmer.Stem(test.initialString)
			assert.Equal(t, test.expectedAns, strings.Join(ans, " "))
		})
	}
}