package testreport_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/becheran/go-testreport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateReport(t *testing.T) {
	var suite = []struct {
		result   testreport.Result
		template string
		out      string
	}{
		{testreport.Result{}, "foo bar", "foo bar"},
		{testreport.Result{Passed: 42}, "foo bar {{.Passed}}", "foo bar 42"},
		{testreport.Result{Passed: 42, Failed: 12, Skipped: 0, Duration: 120}, `{{.Passed}} {{.Failed}} {{.Skipped}} {{.Duration | printf "%d"}}`, "42 12 0 120"},
		{testreport.Result{PackageResult: []testreport.PackageResult{
			{
				Name:          "foo",
				Duration:      time.Second * 125,
				PackageResult: testreport.FTPSSkip,
				Tests: []testreport.TestResult{
					{Name: "t1", Duration: time.Minute, TestResult: testreport.FTSFail, Output: []testreport.OutputLine{
						{Time: time.Time{}, Text: "foo"},
						{Time: time.Time{}, Text: "bar"},
					}},
				},
			}}}, `{{range .PackageResult}}Result:
name={{.Name}}
duration={{.Duration}}
res={{.PackageResult.Icon}}
{{range .Tests}}Tests:
   {{.Name}}: {{.Duration}} {{.TestResult}} {{range .Output}}{{.Time}} {{.Text}} {{end}}
{{end}}{{end}}`,
			"Result:\nname=foo\nduration=2m5s\nres=⏩\nTests:\n   t1: 1m0s fail 0001-01-01 00:00:00 +0000 UTC foo 0001-01-01 00:00:00 +0000 UTC bar \n"},
	}
	for i, s := range suite {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			buff := bytes.NewBuffer(nil)
			temp, err := template.New("test").Parse(s.template)
			require.Nil(t, err)

			assert.Nil(t, testreport.CreateReport(s.result, buff, temp))

			fmt.Println(buff.String())
			assert.Equal(t, s.out, buff.String())
		})
	}
}

func TestCreateDefaultReport(t *testing.T) {
	var suite = []struct {
		result testreport.Result
		out    string
	}{
		{testreport.Result{
			Tests:    19,
			Passed:   12,
			Skipped:  3,
			Failed:   4,
			Duration: time.Second * 124,
			PackageResult: []testreport.PackageResult{
				{
					Name:          "name/p1",
					Duration:      time.Second * 12,
					PackageResult: testreport.FTSPass,
					Succeeded:     1,
					Tests: []testreport.TestResult{
						{
							Name:       "t1",
							Duration:   time.Second,
							TestResult: testreport.FTPSSkip,
							Output: []testreport.OutputLine{
								{Time: time.Time{}, Text: "foo"},
								{Time: time.Time{}, Text: "bar"},
							},
						},
					},
				},
			},
		},
			`# Test Report

Total: 19 ✔️ Passed: 12 ⏩ Skipped: 3 ❌ Failed: 4 ⏱️ Duration: 2m4s

<details><summary>✔️ 1/1 name/<b>p1</b> 12s</summary>

<blockquote><details><summary>⏩ t1 1s</summary><blockquote>

` + "`" + `00:00:00.000` + "`" + ` foo

` + "`" + `00:00:00.000` + "`" + ` bar

</blockquote></details></blockquote></details>
`,
		},
	}
	for i, s := range suite {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			buff := bytes.NewBuffer(nil)
			temp, err := testreport.GetTemplate("")
			require.Nil(t, err)
			assert.Nil(t, testreport.CreateReport(s.result, buff, temp))

			fmt.Println(buff.String())
			assert.Equal(t, s.out, buff.String())
		})
	}
}

func TestCreateReportCustomVars(t *testing.T) {
	buff := bytes.NewBuffer(nil)
	temp, err := template.New("test").Parse(`{{.Vars.foo}} {{.Vars.bar}}`)
	require.Nil(t, err)
	res := testreport.Result{Vars: map[string]string{"foo": "one", "bar": "other"}}

	assert.Nil(t, testreport.CreateReport(res, buff, temp))

	fmt.Println(buff.String())
	assert.Equal(t, `one other`, buff.String())
}

func TestPackageName(t *testing.T) {
	var suite = []struct {
		name string
		path string
		pack string
	}{
		{"", "", ""},
		{"foo", "", "foo"},
		{"foo/bar", "foo/", "bar"},
		{"github.com/becheran/go-testreport", "github.com/becheran/", "go-testreport"},
		{"foo/", "foo/", ""},
	}
	for _, s := range suite {
		t.Run(s.name, func(t *testing.T) {
			p := testreport.PackageName(s.name)
			assert.Equal(t, s.path, p.Path())
			assert.Equal(t, s.pack, p.Package())
		})
	}
}

func TestParseTestJson(t *testing.T) {
	const timeStr = "2023-02-01T19:55:05.5952434+01:00"
	timeObj, err := time.Parse(time.RFC3339, timeStr)
	require.Nil(t, err)

	var suite = []struct {
		json   string
		result testreport.Result
		isErr  bool
	}{
		{"foo bar", testreport.Result{}, true},
		{"{}", testreport.Result{PackageResult: []testreport.PackageResult{}}, false},
		{`{"Time":"` + timeStr + `","Action":"run","Package":"github.com/becheran/go-testreport","Test":"TestIsLess"}
{"Time":"` + timeStr + `","Action":"pass","Package":"github.com/becheran/go-testreport","Test":"TestIsLess","Elapsed":0}
{"Time":"` + timeStr + `","Action":"pass","Package":"github.com/becheran/go-testreport","Elapsed":1.117}
{"Time":"` + timeStr + `","Action":"skip","Package":"github.com/becheran/foo","Elapsed":0}
`, testreport.Result{Tests: 1, Passed: 1, Duration: time.Second, PackageResult: []testreport.PackageResult{
			{
				Name:          "github.com/becheran/go-testreport",
				Duration:      1117000000,
				PackageResult: testreport.FTSPass,
				Succeeded:     1,
				Tests: []testreport.TestResult{
					{Name: "TestIsLess",
						TestResult: testreport.FTSPass,
						Output: []testreport.OutputLine{
							{Time: timeObj},
							{Time: timeObj},
						}},
				}}}}, false},
	}
	for i, s := range suite {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			res, err := testreport.ParseTestJson(strings.NewReader(s.json))
			if s.isErr {
				assert.NotNil(t, err)
				assert.Empty(t, res)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, s.result, res)
			}
		})
	}
}

func TestPackageResultString(t *testing.T) {
	var suite = []struct {
		res testreport.PackageResult
		str string
	}{
		{
			testreport.PackageResult{Name: "github.com/becheran/go-testreport/cmd/TestReport", PackageResult: testreport.FTPSSkip},
			"?       github.com/becheran/go-testreport/cmd/TestReport [no test files]",
		},
		{
			testreport.PackageResult{Name: "foo", PackageResult: testreport.FTSPass, Duration: time.Second * 130},
			"ok      foo 2m10s",
		},
		{
			testreport.PackageResult{
				Name:          "foo",
				PackageResult: testreport.FTSFail,
				Duration:      time.Minute * 2,
				Tests: []testreport.TestResult{
					{Name: "t1", Duration: time.Minute, TestResult: testreport.FTSFail, Output: []testreport.OutputLine{
						{Text: "output_1\n"},
						{Text: "output_2\n"},
					}},
				},
			},
			"output_1\noutput_2\nFAIL    foo 2m0s",
		},
	}
	for _, s := range suite {
		t.Run(s.str, func(t *testing.T) {
			assert.Equal(t, s.str, s.res.String())
		})
	}
}
