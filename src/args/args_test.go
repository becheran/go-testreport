package args_test

import (
	"flag"
	"os"
	"testing"

	"github.com/becheran/go-testreport/src/args"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseArgs_NoFile_UseStdOutAndIn(t *testing.T) {
	res, err := args.ParseArgs([]string{"exe"}, flag.NewFlagSet("test", flag.PanicOnError))

	assert.Nil(t, err)
	assert.Equal(t, os.Stdout, res.OutputStream)
	assert.Equal(t, os.Stdin, res.InputStream)
}

func TestParseArgs_Files_UseFiles(t *testing.T) {
	file, err := os.Create(t.TempDir() + "/file")
	require.Nil(t, err)
	defer file.Close()

	res, err := args.ParseArgs(
		[]string{"exe", "-input", file.Name(), "-output", file.Name()},
		flag.NewFlagSet("test", flag.PanicOnError),
	)
	assert.Nil(t, err)

	_, err = res.OutputStream.Write([]byte("test"))
	require.Nil(t, err)
	res.OutputStream.Close()
	readBytes := make([]byte, 4)
	_, err = res.InputStream.Read(readBytes)
	require.Nil(t, err)
	res.InputStream.Close()
	assert.Equal(t, "test", string(readBytes))
	assert.True(t, res.NonZeroExitOnFailure)
}

func TestParseArgs_CommaSeparatedList_ExpectedOutput(t *testing.T) {
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
	file, err := os.Create(t.TempDir() + "/file")
	require.Nil(t, err)
	defer file.Close()
	for _, s := range suite {
		t.Run(s.in, func(t *testing.T) {
			res, err := args.ParseArgs([]string{"exe", "-vars", s.in, "-output", file.Name()}, flag.NewFlagSet("test", flag.PanicOnError))
			if s.isErr {
				assert.NotNil(t, err)
				assert.Empty(t, res.EnvArgs)
			} else {
				res.OutputStream.Close()
				assert.Nil(t, err)
				assert.Equal(t, s.out, res.EnvArgs)
			}
		})
	}
}
