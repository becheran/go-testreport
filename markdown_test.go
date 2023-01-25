package testreport_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/becheran/go-testreport"

	"github.com/stretchr/testify/assert"
)

func TestResultToMarkdown(t *testing.T) {
	res := testreport.Result{Passed: 1, Failed: 2, Skipped: 3, Duration: time.Second * 13, PackageResult: map[string]*testreport.PackageResult{}}

	defaultReport := testreport.ResultToMarkdown(res)

	assert.Equal(t, "# Test Report\n\nTotal: 6 ✔️ Passed: 1 ⏩ Skipped: 3 ❌ Failed: 2 ⏱️ Duration: 13s\n\n", string(defaultReport))
}

func TestIsLess(t *testing.T) {
	var suite = []struct {
		aStatus     testreport.FinalTestStatus
		bStatus     testreport.FinalTestStatus
		aElapsedSec float64
		bElapsedSec float64
		isLess      bool
	}{
		{testreport.FTSPass, testreport.FTSPass, 0, 0, false},
		{testreport.FTSFail, testreport.FTSFail, 0, 0, false},
		{testreport.FTSFail, testreport.FTSPass, 0, 0, false},
		{testreport.FTSFail, testreport.FTSPass, 0, 100, false},
		{testreport.FTSFail, testreport.FTSFail, 100, 0, false},

		{testreport.FTSPass, testreport.FTSFail, 0, 0, true},
		{testreport.FTSPass, testreport.FTSFail, 100, 0, true},
		{testreport.FTSFail, testreport.FTSFail, 0, 100, true},
		{testreport.FTSPass, testreport.FTSPass, 0, 100, true},
		{testreport.FTPSSkip, testreport.FTSPass, 0, 0, true},
		{testreport.FTPSSkip, testreport.FTSPass, 100, 0, true},
		{testreport.FTPSSkip, testreport.FTSFail, 0, 0, true},
		{testreport.FTPSSkip, testreport.FTSFail, 100, 0, true},
	}
	for _, s := range suite {
		t.Run(fmt.Sprintf("A(%s %f) B(%s %f)", s.aStatus.Icon(), s.aElapsedSec, s.bStatus.Icon(), s.bElapsedSec), func(t *testing.T) {
			assert.Equal(t, s.isLess, testreport.IsLess(s.aStatus, s.bStatus, s.aElapsedSec, s.bElapsedSec))
		})
	}
}

func TestReplaceMarkdown(t *testing.T) {
	var suite = []struct {
		mdIn    string
		escaped string
	}{
		{"foo", "foo"},
		{"*", "\\*"},
		{"_italic_", "\\_italic\\_"},
		{"*foo*", "\\*foo\\*"},
		{"[link](\"a lining\")", "\\[link\\](\"a&nbsp;lining\")"},
		{"This is a reals backslash: \\", "This&nbsp;is&nbsp;a&nbsp;reals&nbsp;backslash:&nbsp;\\\\"},
		{"`code`", "\\`code\\`"},
		{" ", "&nbsp;"},
		{"	", "&nbsp;&nbsp;&nbsp;&nbsp;"},
	}
	for _, s := range suite {
		t.Run(fmt.Sprintf("%s => %s", s.mdIn, s.escaped), func(t *testing.T) {
			assert.Equal(t, s.escaped, testreport.EscapeMarkdown(s.mdIn))
		})
	}
}

func TestPackageTestPassRatio(t *testing.T) {
	const s = testreport.NonBreakingSpace
	var suite = []struct {
		passed int
		tests  int
		digits int
		result string
	}{
		{0, 0, 1, "0/0"},
		{0, 1, 1, "0/1"},
		{1, 0, 1, "1/0"},
		{10, 10, 2, "10/10"},
		{0, 0, 2, s + "0/0" + s},
		{1, 10, 2, s + "1/10"},
		{10, 1, 2, "10/1" + s},
		{10, 1, 3, s + "10/1" + s + s},
		{666, 1, 3, "666/1" + s + s},
		{10, 1, 4, s + s + "10/1" + s + s + s},
	}
	for _, s := range suite {
		t.Run(fmt.Sprintf("PackageTestPassRatio(%d,%d,%d)", s.passed, s.tests, s.digits), func(t *testing.T) {
			assert.Equal(t, s.result, testreport.PackageTestPassRatio(s.passed, s.tests, s.digits))
		})
	}
}
