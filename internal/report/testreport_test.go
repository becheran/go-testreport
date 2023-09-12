package report_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/becheran/go-testreport/internal/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateReport(t *testing.T) {
	var suite = []struct {
		result   report.Result
		template string
		out      string
	}{
		{report.Result{}, "foo bar", "foo bar"},
		{report.Result{Passed: 42}, "foo bar {{.Passed}}", "foo bar 42"},
		{report.Result{Passed: 42, Failed: 12, Skipped: 0, Duration: 120}, `{{.Passed}} {{.Failed}} {{.Skipped}} {{.Duration | printf "%d"}}`, "42 12 0 120"},
		{report.Result{PackageResult: []report.PackageResult{
			{
				Name:          "foo",
				Duration:      time.Second * 125,
				PackageResult: report.FTPSSkip,
				Tests: []report.TestResult{
					{Name: "t1", Duration: time.Minute, TestResult: report.FTSFail, Output: []report.OutputLine{
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

			assert.Nil(t, report.CreateReport(s.result, buff, temp))

			fmt.Println(buff.String())
			assert.Equal(t, s.out, buff.String())
		})
	}
}

func TestCreateDefaultReport(t *testing.T) {
	var suite = []struct {
		result report.Result
		out    string
	}{
		{report.Result{
			Tests:    19,
			Passed:   12,
			Skipped:  3,
			Failed:   4,
			Duration: time.Second * 124,
			PackageResult: []report.PackageResult{
				{
					Name:          "name/p1",
					Duration:      time.Second * 12,
					PackageResult: report.FTSPass,
					Succeeded:     1,
					Tests: []report.TestResult{
						{
							Name:       "t1",
							Duration:   time.Second,
							TestResult: report.FTPSSkip,
							Output: []report.OutputLine{
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

<details>
    <summary>✔️ 1/1 name/<b>p1</b> 12s</summary>
        
<blockquote>⏩ t1 1s  </blockquote>
</details>
`,
		},
	}
	for i, s := range suite {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			buff := bytes.NewBuffer(nil)
			temp, err := report.GetTemplate("")
			require.Nil(t, err)
			assert.Nil(t, report.CreateReport(s.result, buff, temp))

			assert.Equal(t, s.out, strings.ReplaceAll(buff.String(), "\r\n", "\n"))
		})
	}
}

func TestCreateReportCustomVars(t *testing.T) {
	buff := bytes.NewBuffer(nil)
	temp, err := template.New("test").Parse(`{{.Vars.foo}} {{.Vars.bar}}`)
	require.Nil(t, err)
	res := report.Result{Vars: map[string]string{"foo": "one", "bar": "other"}}

	assert.Nil(t, report.CreateReport(res, buff, temp))

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
			p := report.PackageName(s.name)
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
		result report.Result
		isErr  bool
	}{
		{"foo bar", report.Result{}, true},
		{"{}", report.Result{PackageResult: []report.PackageResult{}}, false},
		{`{"Time":"` + timeStr + `","Action":"run","Package":"github.com/becheran/go-testreport","Test":"TestIsLess"}
{"Time":"` + timeStr + `","Action":"pass","Package":"github.com/becheran/go-testreport","Test":"TestIsLess","Elapsed":0}
{"Time":"` + timeStr + `","Action":"pass","Package":"github.com/becheran/go-testreport","Elapsed":1.117}
{"Time":"` + timeStr + `","Action":"skip","Package":"github.com/becheran/foo","Elapsed":0}
`, report.Result{Tests: 1, Passed: 1, Duration: time.Second, PackageResult: []report.PackageResult{
			{
				Name:          "github.com/becheran/go-testreport",
				Duration:      1117000000,
				PackageResult: report.FTSPass,
				Succeeded:     1,
				Tests: []report.TestResult{
					{Name: "TestIsLess",
						TestResult: report.FTSPass,
						Output: []report.OutputLine{
							{Time: timeObj},
							{Time: timeObj},
						}},
				}}}}, false},
	}
	for i, s := range suite {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			res, err := report.ParseTestJson(strings.NewReader(s.json))
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
		res report.PackageResult
		str string
	}{
		{
			report.PackageResult{Name: "github.com/becheran/go-testreport/cmd/TestReport", PackageResult: report.FTPSSkip},
			"?       github.com/becheran/go-testreport/cmd/TestReport [no test files]",
		},
		{
			report.PackageResult{Name: "foo", PackageResult: report.FTSPass, Duration: time.Second * 130},
			"ok      foo 2m10s",
		},
		{
			report.PackageResult{
				Name:          "foo",
				PackageResult: report.FTSFail,
				Duration:      time.Minute * 2,
				Tests: []report.TestResult{
					{Name: "t1", Duration: time.Minute, TestResult: report.FTSFail, Output: []report.OutputLine{
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
