# Go Test Report

[![Pipeline Status](https://github.com/becheran/go-testreport/actions/workflows/go.yml/badge.svg)](https://github.com/becheran/go-testreport/actions/workflows/go.yml)
[![Go Report Card][go-report-image]][go-report-url]
[![PRs Welcome][pr-welcome-image]][pr-welcome-url]
[![License][license-image]][license-url]
[![GHAction][gh-action-image]][gh-action-url]

[license-url]: https://github.com/becheran/go-testreport/blob/main/LICENSE
[license-image]: https://img.shields.io/badge/License-MIT-brightgreen.svg
[go-report-image]: https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat
[go-report-url]: https://goreportcard.com/report/github.com/becheran/go-testreport
[pr-welcome-image]: https://img.shields.io/badge/PRs-welcome-brightgreen.svg
[pr-welcome-url]: https://github.com/becheran/go-testreport/blob/main/CONTRIBUTING.md
[gh-action-image]: https://img.shields.io/badge/Get-GH_Action-blue
[gh-action-url]: https://github.com/marketplace/actions/golang-test-report

Generate a markdown test report from the go json test result.

Matches perfectly with [github job summaries]( https://github.blog/news-insights/product-news/supercharging-github-actions-with-job-summaries/) to visualize test results:

![ReportExample](./doc/GitHubReport.png)

The default output sorts the tests by failing and slowest execution time.

## Install

### Go

Install via the go install command:

``` sh
go install github.com/becheran/go-testreport@latest
```

### Binaries

Or use the pre-compiled binaries for Linux, Windows, and Mac OS from the [github releases page](https://github.com/becheran/go-testreport/releases).

## Usage

Run the following command to get a list of all available command line options:

``` sh
go-testreport -h
```

### Input and Output

When `-input` and `-output` is not set, the stdin stream will be used and return the result will be written to stdout:

``` sh
go test ./... -json | go-testreport > result.html
```

Use the `-input` and `-output` file to set files for the input and output:

``` sh
go-testreport -input result.json -output result.html
```

### Templates

Customize by providing a own [template file](https://pkg.go.dev/text/template). See also the [default markdown template](./src/report/templates/md.tmpl) which is used if the `-template` argument is left empty. With the `vars` options custom dynamic values can be passed to the template from the outside which can be resolved within the template:

``` sh
go test ./... -json | go-testreport -template=./html.tmpl -vars="Title:Test Report Linux" > $GITHUB_STEP_SUMMARY
```

### GitHub Actions

The [Golang Test Report](https://github.com/marketplace/actions/golang-test-report) from the marketplace can be used to integrate the go-testreport tool into an GitHub workflow:

``` yaml
- name: Test
  run: go test ./... -json > report.json
- name: Report
  uses: becheran/go-testreport@main
  with:
    input: report.json
```
