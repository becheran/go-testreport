package args_test

import (
	"testing"

	"github.com/becheran/go-testreport/internal/args"
	"github.com/stretchr/testify/assert"
)

func TestParseCommaSeparatedList(t *testing.T) {
	var suite = []struct {
		in    string
		out   map[string]string
		isErr bool
	}{
		{"", map[string]string{}, false},
		{"v:a", map[string]string{"v": "a"}, false},
		{"v:a,b:z", map[string]string{"v": "a", "b": "z"}, false},
		{"v:a,,bar:baz", map[string]string{"v": "a", "bar": "baz"}, false},
		{"a:b,b:,c:2", map[string]string{"a": "b", "b": "", "c": "2"}, false},

		{"f", nil, true},
		{":12", nil, true},
		{"f:b:a", nil, true},
		{"f:b,:1", nil, true},
		{"a:1,a:2", nil, true},
		{"a:1,a:2", nil, true},
		{"a:1,a:2", nil, true},
	}
	for _, s := range suite {
		t.Run(s.in, func(t *testing.T) {
			res, err := args.ParseCommaSeparatedList(s.in)
			if s.isErr {
				assert.NotNil(t, err)
				assert.Empty(t, res)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, s.out, res)
			}
		})
	}
}
