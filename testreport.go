package testreport

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/template"
	"time"
)

type FinalTestStatus uint8

const (
	FTPSSkip FinalTestStatus = iota
	FTSPass
	FTSFail
)

func (fs FinalTestStatus) String() string {
	switch fs {
	case FTSPass:
		return "pass"
	case FTSFail:
		return "fail"
	case FTPSSkip:
		return "skip"
	default:
		return ""
	}
}

func (fs FinalTestStatus) Icon() string {
	switch fs {
	case FTSPass:
		return "✔️"
	case FTSFail:
		return "❌"
	case FTPSSkip:
		return "⏩"
	default:
		return ""
	}
}

func FinalTestStatusFromAction(e TestAction) *FinalTestStatus {
	var status FinalTestStatus
	switch e {
	case TAFail:
		status = FTSFail
	case TAPass:
		status = FTSPass
	case TASkip:
		status = FTPSSkip
	default:
		return nil
	}
	return &status
}

type OutputLine struct {
	Time time.Time
	Text string
}

type TestResult struct {
	Name       string
	Duration   time.Duration
	Output     []OutputLine
	TestResult FinalTestStatus
}

type PackageName string

func (p PackageName) Package() string {
	lastIdx := strings.LastIndex(string(p), "/")
	if lastIdx > 0 {
		return string(p)[lastIdx+1:]
	}
	return string(p)
}

func (p PackageName) Path() string {
	lastIdx := strings.LastIndex(string(p), "/")
	if lastIdx > 0 {
		return string(p)[:lastIdx+1]
	}
	return ""
}

type PackageResult struct {
	Name          PackageName
	Duration      time.Duration
	PackageResult FinalTestStatus
	Succeeded     int
	Tests         []TestResult
}

func (p PackageResult) String() string {
	res := strings.Builder{}
	switch p.PackageResult {
	case FTSPass:
		res.WriteString("ok      ")
	case FTPSSkip:
		res.WriteString("?       ")
	case FTSFail:
		for _, test := range p.Tests {
			if test.TestResult == FTSFail {
				for _, line := range test.Output {
					res.WriteString(line.Text)
				}
			}
		}
		res.WriteString("FAIL    ")
	default:
		panic("BUG! Unexpected package result" + p.PackageResult.String())
	}
	res.WriteString(string(p.Name))
	res.WriteString(" ")
	if p.PackageResult == FTPSSkip && len(p.Tests) == 0 {
		res.WriteString("[no test files]")
	} else {
		res.WriteString(p.Duration.String())
	}
	return res.String()
}

type Result struct {
	Failed        uint
	Passed        uint
	Skipped       uint
	Tests         uint
	Duration      time.Duration
	PackageResult []PackageResult
	Vars          map[string]string
}

func ParseTestJson(in io.Reader) (result Result, err error) {
	packageResult := make(map[string]*PackageResult)
	testResultForPackage := make(map[string]map[string]*TestResult)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Bytes()
		var evt TestEvent
		if err := json.Unmarshal(line, &evt); err != nil {
			return Result{}, err
		}
		if _, packageExists := packageResult[evt.Package]; !packageExists {
			res := PackageResult{
				Name: PackageName(evt.Package),
			}
			packageResult[evt.Package] = &res
			testResultForPackage[evt.Package] = map[string]*TestResult{}
		}

		if evt.Test == "" {
			if status := FinalTestStatusFromAction(evt.Action); status != nil {
				packageResult[evt.Package].PackageResult = *status
				result.Duration += time.Second * time.Duration(evt.ElapsedSec)
			}
			packageResult[evt.Package].Duration = time.Duration(float64(time.Second) * evt.ElapsedSec)
		} else {
			if testRes, testExists := testResultForPackage[evt.Package][evt.Test]; testExists {
				testRes.Output = append(testRes.Output, OutputLine{Time: evt.Time, Text: evt.Output})
			} else {
				testResultForPackage[evt.Package][evt.Test] = &TestResult{
					Name:   evt.Test,
					Output: []OutputLine{{Time: evt.Time, Text: evt.Output}},
				}
			}
			if status := FinalTestStatusFromAction(evt.Action); status != nil {
				test := testResultForPackage[evt.Package][evt.Test]
				test.TestResult = *status
				test.Duration = time.Duration(float64(time.Second) * evt.ElapsedSec)
				switch *status {
				case FTSPass:
					result.Passed++
				case FTSFail:
					result.Failed++
				case FTPSSkip:
					result.Skipped++
				}
			}
		}
	}
	result.Tests = result.Skipped + result.Failed + result.Passed
	result.PackageResult = make([]PackageResult, 0, len(packageResult))
	for _, val := range packageResult {
		if val.PackageResult == FTPSSkip {
			// Ignore skipped packages in report
			continue
		}
		res := *val
		tests := testResultForPackage[string(val.Name)]
		res.Tests = make([]TestResult, 0, len(tests))
		for _, test := range tests {
			res.Tests = append(res.Tests, *test)
			if test.TestResult == FTSPass || test.TestResult == FTPSSkip {
				res.Succeeded++
			}
		}
		result.PackageResult = append(result.PackageResult, res)
	}
	sort.Slice(result.PackageResult, func(i, j int) bool {
		return !IsLess(result.PackageResult[i].PackageResult, result.PackageResult[j].PackageResult,
			result.PackageResult[i].Duration, result.PackageResult[j].Duration)
	})
	return result, nil
}

func CreateReport(result Result, out io.Writer, temp *template.Template) (err error) {
	if temp == nil {
		return fmt.Errorf("template must be defined")
	}
	return temp.Execute(out, result)
}
