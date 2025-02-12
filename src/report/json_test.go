package report_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/becheran/go-testreport/src/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalTestAction(t *testing.T) {
	const dateStr = "2022-12-27T20:45:01.5635437+01:00"
	const in = `{"Time":"` + dateStr + `","Action":"output","Package":"github.com/becheran/wildmatch-go","Test":"TestIsMatch/___#04","Output":"    --- PASS: TestIsMatch/___#04 (0.00s)\n"}`
	var evt report.TestEvent

	err := json.Unmarshal([]byte(in), &evt)
	require.Nil(t, err)

	date, err := time.Parse(time.RFC3339, dateStr)
	require.Nil(t, err)
	assert.Equal(t, evt, report.TestEvent{
		Time:    date,
		Action:  report.TAOutput,
		Package: "github.com/becheran/wildmatch-go",
		Test:    "TestIsMatch/___#04",
		Output:  "    --- PASS: TestIsMatch/___#04 (0.00s)\n",
	})
}
