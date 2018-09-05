package structs_test

import (
	"testing"

	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

func TestNumberWithinRange(t *testing.T) {

	var max int = 100
	var distance int = 5

	tt := []struct {
		botValue       int
		nodeValue      int
		expectedResult bool
	}{
		{50, 95, false},
		{50, 5, false},
		{97, 99, true},
		{97, 2, true},
		{4, 2, true},
		{4, 99, true},
		{97, 3, false},
		{3, 97, false},
	}

	for _, tc := range tt {
		actual := structs.NumberWithinRange(tc.botValue, distance, max, tc.nodeValue)
		if tc.expectedResult != actual {
			t.Errorf("Number within range didn't return correct answer. Bot: %d nodeValue: %d Actual result: %v", tc.botValue, tc.nodeValue, actual)
		}
	}

}
