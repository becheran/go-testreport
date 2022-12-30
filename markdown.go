package testreport

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

func IsLess(a_PackageResult, b_PackageResult FinalTestStatus, a_ElapsedSec, b_ElapsedSec float64) bool {
	if a_PackageResult == FTSPass && b_PackageResult != FTSPass {
		return true
	} else if b_PackageResult == FTSPass && a_PackageResult != FTSPass {
		return false
	}
	if a_ElapsedSec < b_ElapsedSec {
		return true
	}
	return false
}

func ResultToMarkdown(result Result) []byte {
	buf := bytes.NewBufferString("# Test Report\n\n")
	var total = result.Failed + result.Passed + result.Skipped
	buf.WriteString(fmt.Sprintf("Total: %d ✔️ Passed: %d ⏩ Skipped: %d ❌ Failed: %d ⏱️ Duration: %s\n",
		total, result.Passed, result.Skipped, result.Failed, result.Duration))
	res := maps.Values(result.PackageResult)
	sort.Slice(res, func(i, j int) bool {
		return !IsLess(res[i].PackageResult, res[j].PackageResult, res[i].ElapsedSec, res[j].ElapsedSec)
	})
	for _, packRes := range res {
		if packRes.PackageResult == FTPSSkip {
			// Do not print skipped packages
			continue
		}
		buf.WriteString("\n<details><summary>")
		buf.WriteString(fmt.Sprintf("%s %s %s %.2fs", packRes.PackageResult.Icon(), PackageTestPassRatio(packRes), packRes.Name, packRes.ElapsedSec))
		buf.WriteString("</summary>")
		tests := maps.Values(packRes.Tests)
		sort.Slice(tests, func(i, j int) bool {
			return !IsLess(tests[i].TestResult, tests[j].TestResult, tests[i].ElapsedSec, tests[j].ElapsedSec)
		})
		for _, testRes := range tests {
			buf.WriteString("<blockquote><details><summary>")
			buf.WriteString(fmt.Sprintf("%s %s %.2fs", testRes.TestResult.Icon(), testRes.Name, testRes.ElapsedSec))
			buf.WriteString("</summary><blockquote>\n\n")
			for _, outputLine := range testRes.Output {
				if outputLine.Text != "" {
					buf.WriteString(fmt.Sprintf("`%s` %s\n", outputLine.Time.Format("15:04:05.000"), EscapeMarkdown(outputLine.Text)))
				}
			}
			buf.WriteString("</blockquote></details></blockquote>")
		}
		buf.WriteString("</details>")
	}
	buf.WriteString("\n")

	return buf.Bytes()
}

func PackageTestPassRatio(res *PackageResult) string {
	passed := 0
	for _, p := range res.Tests {
		if p.TestResult != FTSFail {
			passed++
		}
	}
	return fmt.Sprintf("%0d/%0d", passed, len(res.Tests))
}

func EscapeMarkdown(input string) (escapedMarkdown string) {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"#", "\\#",
		"[", "\\[",
		"]", "\\]",
		"\\", "\\\\",
		"`", "\\`",
		" ", "&nbsp;",
		"	", "&nbsp;&nbsp;&nbsp;&nbsp;",
	)
	return replacer.Replace(input)
}
