package testreport

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
	Tests         map[string]*TestResult
}

type Result struct {
	Failed        uint
	Passed        uint
	Skipped       uint
	Tests         uint
	Duration      time.Duration
	PackageResult map[string]*PackageResult
	Vars          map[string]string
}

func ParseTestJson(in io.Reader) (result Result, err error) {
	result = Result{
		PackageResult: make(map[string]*PackageResult),
	}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Bytes()
		var evt TestEvent
		if err := json.Unmarshal(line, &evt); err != nil {
			return Result{}, err
		}
		if _, packageExists := result.PackageResult[evt.Package]; !packageExists {
			res := PackageResult{
				Name:  PackageName(evt.Package),
				Tests: make(map[string]*TestResult),
			}
			result.PackageResult[evt.Package] = &res
		}

		if evt.Test == "" {
			if status := FinalTestStatusFromAction(evt.Action); status != nil {
				result.PackageResult[evt.Package].PackageResult = *status
				result.Duration += time.Second * time.Duration(evt.ElapsedSec)
			}
			result.PackageResult[evt.Package].Duration = time.Duration(float64(time.Second) * evt.ElapsedSec)
		} else {
			if testRes, testExists := result.PackageResult[evt.Package].Tests[evt.Test]; testExists {
				testRes.Output = append(testRes.Output, OutputLine{Time: evt.Time, Text: evt.Output})
			} else {
				result.PackageResult[evt.Package].Tests[evt.Test] = &TestResult{
					Name:   evt.Test,
					Output: []OutputLine{{Time: evt.Time, Text: evt.Output}},
				}
			}
			if status := FinalTestStatusFromAction(evt.Action); status != nil {
				test := result.PackageResult[evt.Package].Tests[evt.Test]
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
	return result, nil
}

func CreateReport(result Result, out io.Writer, temp *template.Template) (err error) {
	if temp == nil {
		return fmt.Errorf("template must be defined")
	}
	return temp.Execute(out, result)
}
