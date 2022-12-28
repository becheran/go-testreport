package testreport

import (
	"encoding/json"
	"time"
)

type TestAction uint8

const (
	TAUnknown TestAction = iota
	TARun                // the test has started running
	TAPause              // the test has been paused
	TACont               // the test has continued running
	TAPass               // the test passed
	TABench              // the benchmark printed log output but did not fail
	TAFail               // the test or benchmark failed
	TAOutput             // the test printed output
	TASkip               // the test was skipped or the package contained no tests
)

var taStrings = []string{"run", "pause", "cont", "pass", "bench", "fail", "output", "skip"}

func (ta TestAction) String() string {
	idx := int(ta) - 1
	if idx < 0 && idx < len(taStrings) {
		return taStrings[ta-1]
	}
	return "unknown"
}

func TestActionFromString(s string) TestAction {
	for idx, str := range taStrings {
		if str == s {
			return TestAction(idx + 1)
		}
	}
	return TAUnknown
}

func (b *TestAction) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*b = TestActionFromString(s)
	return nil
}

// Test event represents a single json test output line.
// Implements the marshaller interface
// From https://pkg.go.dev/cmd/test2json
type TestEvent struct {
	Time       time.Time  `json:"time,omitempty"` // encodes as an RFC3339-format string
	Action     TestAction `json:"action,omitempty"`
	Package    string     `json:"package,omitempty"`
	Test       string     `json:"test,omitempty"`
	ElapsedSec float64    `json:"elapsed,omitempty"` // seconds
	Output     string     `json:"output,omitempty"`
}
