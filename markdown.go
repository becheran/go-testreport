package testreport

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/becheran/go-testreport/internal/numhelper"
	"golang.org/x/exp/maps"
)

const NonBreakingSpace = "&nbsp;"

func IsLess(a_PackageResult, b_PackageResult FinalTestStatus, a_ElapsedSec, b_ElapsedSec float64) bool {
	if a_PackageResult == b_PackageResult {
		return a_ElapsedSec < b_ElapsedSec
	}
	return a_PackageResult < b_PackageResult
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

	digitsPackageTests := 0
	for _, packRes := range res {
		if packRes.PackageResult == FTPSSkip {
			// Do not print skipped packages
			continue
		}
		digits := numhelper.Digits(len(packRes.Tests))
		if digitsPackageTests < digits {
			digitsPackageTests = digits
		}
	}

	for _, packRes := range res {
		if packRes.PackageResult == FTPSSkip {
			// Do not print skipped packages
			continue
		}
		buf.WriteString("\n<details><summary>")

		// Highlight last part of package name
		var packageHtml string
		lastIdx := strings.LastIndex(packRes.Name, "/")
		if lastIdx > 0 {
			lastIdx++ // Exclude slash
			packageHtml = packRes.Name[:lastIdx] + "<b>" + packRes.Name[lastIdx:] + "</b>"
		} else {
			packageHtml = fmt.Sprintf("<b>%s</b>", packRes.Name)
		}

		passed := 0
		for _, p := range packRes.Tests {
			if p.TestResult != FTSFail {
				passed++
			}
		}
		buf.WriteString(fmt.Sprintf("%s %s %s %.2fs", packRes.PackageResult.Icon(), PackageTestPassRatio(passed, len(packRes.Tests), digitsPackageTests), packageHtml, packRes.ElapsedSec))
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

func PackageTestPassRatio(passed int, tests int, digitsPackageTests int) string {
	result := bytes.NewBuffer(nil)
	// Pad with whitespace
	for i := 0; i < digitsPackageTests-numhelper.Digits(passed); i++ {
		result.WriteString(NonBreakingSpace)
	}
	result.WriteString(fmt.Sprintf("%d/%d", passed, tests))
	// Pad with whitespace
	for i := 0; i < digitsPackageTests-numhelper.Digits(tests); i++ {
		result.WriteString(NonBreakingSpace)
	}

	return result.String()
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
		" ", NonBreakingSpace,
		"	", NonBreakingSpace+NonBreakingSpace+NonBreakingSpace+NonBreakingSpace,
	)
	return replacer.Replace(input)
}
