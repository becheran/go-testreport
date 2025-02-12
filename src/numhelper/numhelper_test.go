package numhelper_test

import (
	"fmt"
	"testing"

	"github.com/becheran/go-testreport/src/numhelper"
	"github.com/stretchr/testify/assert"
)

func TestDigits(t *testing.T) {
	suite := []struct {
		input  int
		digits int
	}{
		{2, 1},
		{-1, 1},
		{0, 1},
		{9, 1},
		{10, 2},
		{11, 2},
		{99, 2},
		{999, 3},
		{-10, 2},
	}
	for _, test := range suite {
		t.Run(fmt.Sprintf("Digits(%d)=%d", test.input, test.digits), func(t *testing.T) {
			assert.Equal(t, test.digits, numhelper.Digits(test.input))
		})
	}
}
