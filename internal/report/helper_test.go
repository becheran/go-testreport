package report_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/becheran/go-testreport/internal/report"
	"github.com/stretchr/testify/assert"
)

func TestIsLess(t *testing.T) {
	var suite = []struct {
		aStatus     report.FinalTestStatus
		bStatus     report.FinalTestStatus
		aElapsedSec float64
		bElapsedSec float64
		isLess      bool
	}{
		{report.FTSPass, report.FTSPass, 0, 0, false},
		{report.FTSFail, report.FTSFail, 0, 0, false},
		{report.FTSFail, report.FTSPass, 0, 0, false},
		{report.FTSFail, report.FTSPass, 0, 100, false},
		{report.FTSFail, report.FTSFail, 100, 0, false},

		{report.FTSPass, report.FTSFail, 0, 0, true},
		{report.FTSPass, report.FTSFail, 100, 0, true},
		{report.FTSFail, report.FTSFail, 0, 100, true},
		{report.FTSPass, report.FTSPass, 0, 100, true},
		{report.FTPSSkip, report.FTSPass, 0, 0, true},
		{report.FTPSSkip, report.FTSPass, 100, 0, true},
		{report.FTPSSkip, report.FTSFail, 0, 0, true},
		{report.FTPSSkip, report.FTSFail, 100, 0, true},
	}
	for _, s := range suite {
		t.Run(fmt.Sprintf("A(%s %f) B(%s %f)", s.aStatus.Icon(), s.aElapsedSec, s.bStatus.Icon(), s.bElapsedSec), func(t *testing.T) {
			assert.Equal(t, s.isLess, report.IsLess(s.aStatus, s.bStatus, time.Duration(float64(time.Second)*s.aElapsedSec), time.Duration(float64(time.Second)*s.bElapsedSec)))
		})
	}
}

func TestEscapeMarkdown(t *testing.T) {
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
		{"<details></details>", "&lt;details&gt;&lt;/details&gt;"},
	}
	for _, s := range suite {
		t.Run(fmt.Sprintf("%s => %s", s.mdIn, s.escaped), func(t *testing.T) {
			assert.Equal(t, s.escaped, report.EscapeMarkdown(s.mdIn))
		})
	}
}

func TestEscapeHtml(t *testing.T) {
	var suite = []struct {
		htmlIn  string
		escaped string
	}{
		{"foo", "foo"},
		{"<script>alert(123);</script>", "&lt;script&gt;alert(123);&lt;/script&gt;"},
	}
	for _, s := range suite {
		t.Run(fmt.Sprintf("%s => %s", s.htmlIn, s.escaped), func(t *testing.T) {
			assert.Equal(t, s.escaped, report.EscapeHtml(s.htmlIn))
		})
	}
}
