package testreport

import (
	"html"
	"strings"
	"time"
)

const NonBreakingSpace = "&nbsp;"

func IsLess(a_PackageResult, b_PackageResult FinalTestStatus, a_Duration, b_Duration time.Duration) bool {
	if a_PackageResult == b_PackageResult {
		return a_Duration < b_Duration
	}
	return a_PackageResult < b_PackageResult
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

func EscapeHtml(input string) (escapedHtml string) {
	return html.EscapeString(input)
}
