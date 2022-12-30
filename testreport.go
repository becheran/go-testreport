package testreport

import (
	"bufio"
	"encoding/json"
	"io"
	"time"
)

type FinalTestStatus uint8

const (
	FTSPass FinalTestStatus = iota
	FTSFail
	FTPSSkip
)

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

type packageName = string

type OutputLine struct {
	Time time.Time
	Text string
}

type TestResult struct {
	Name       string
	ElapsedSec float64
	Output     []OutputLine
	TestResult FinalTestStatus
}

type testName = string

type PackageResult struct {
	Name          string
	ElapsedSec    float64
	Tests         map[testName]*TestResult
	PackageResult FinalTestStatus
}

type Result struct {
	Failed        uint
	Passed        uint
	Skipped       uint
	Duration      time.Duration
	PackageResult map[packageName]*PackageResult
}

func CreateReport(in io.Reader) (markdown []byte, err error) {
	result := Result{
		PackageResult: make(map[packageName]*PackageResult),
	}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Bytes()
		var evt TestEvent
		if err := json.Unmarshal(line, &evt); err != nil {
			return nil, err
		}
		if _, packageExists := result.PackageResult[evt.Package]; !packageExists {
			res := PackageResult{
				Name:  evt.Package,
				Tests: make(map[string]*TestResult),
			}
			result.PackageResult[evt.Package] = &res
		}

		if evt.Test == "" {
			if status := FinalTestStatusFromAction(evt.Action); status != nil {
				result.PackageResult[evt.Package].PackageResult = *status
				result.Duration += time.Second * time.Duration(evt.ElapsedSec)
			}
			result.PackageResult[evt.Package].ElapsedSec = evt.ElapsedSec
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
				test.ElapsedSec = evt.ElapsedSec
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
	return ResultToMarkdown(result), nil
}
