package testreport_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/becheran/go-testreport"

	"github.com/stretchr/testify/assert"
)

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
			assert.Equal(t, s.isLess, testreport.IsLess(s.aStatus, s.bStatus, time.Duration(float64(time.Second)*s.aElapsedSec), time.Duration(float64(time.Second)*s.bElapsedSec)))
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
	}
	for _, s := range suite {
		t.Run(fmt.Sprintf("%s => %s", s.mdIn, s.escaped), func(t *testing.T) {
			assert.Equal(t, s.escaped, testreport.EscapeMarkdown(s.mdIn))
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
			assert.Equal(t, s.escaped, testreport.EscapeHtml(s.htmlIn))
		})
	}
}
