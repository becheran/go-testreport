package testreport

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

const space = "&nbsp;"

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
		digits := digits(len(packRes.Tests))
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

		buf.WriteString(fmt.Sprintf("%s %s %s %.2fs", packRes.PackageResult.Icon(), PackageTestPassRatio(packRes, digitsPackageTests), packageHtml, packRes.ElapsedSec))
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

func digits(n int) int {
	if n == 0 {
		return 0
	}
	count := 0
	for n > 0 {
		n = n / 10
		count++
	}
	return count
}

func PackageTestPassRatio(res *PackageResult, digitsPackageTests int) string {
	passed := 0
	for _, p := range res.Tests {
		if p.TestResult != FTSFail {
			passed++
		}
	}

	result := bytes.NewBuffer(nil)
	// Pad with whitespace
	for i := 0; i < digitsPackageTests-digits(passed); i++ {
		result.WriteString(space)
	}
	result.WriteString(fmt.Sprintf("%d/%d", passed, len(res.Tests)))
	// Pad with whitespace
	for i := 0; i < digitsPackageTests-digits(len(res.Tests)); i++ {
		result.WriteString(space)
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
		" ", space,
		"	", space+space+space+space,
	)
	return replacer.Replace(input)
}
