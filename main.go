package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/becheran/go-testreport/internal/args"
	"github.com/becheran/go-testreport/internal/report"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage %s [<options>] <file>:\n", os.Args[0])
		flag.PrintDefaults()
	}

	var templateFile string
	var vars string
	flag.StringVar(&templateFile, "template", "", "Template file for the report generation. If not set, the default template will be applied")
	flag.StringVar(&vars, "vars", "", "Comma separated list of custom variables which can be used in the template. For example -vars=\"Title:Custom Title\"")
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	filepath := flag.Args()[0]
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Failed to open file %s. %s", filepath, err)
	}
	defer file.Close()

	tmp, err := report.GetTemplate(templateFile)
	if err != nil {
		log.Fatalf("Invalid template. %s", err)
	}

	result, err := report.ParseTestJson(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to parse test result %s", err)
	}

	result.Vars, err = args.ParseCommaSeparatedList(vars)
	if err != nil {
		log.Fatalf("Failed to parse variables. %s", err)
	}

	if err := report.CreateReport(result, file, tmp); err != nil {
		log.Fatalf("Failed to create test report. %s", err)
	}

	failed := false
	for _, packRes := range result.PackageResult {
		fmt.Println(packRes)
		if packRes.PackageResult == report.FTSFail {
			failed = true
		}
	}

	if failed {
		os.Exit(1)
	}
}
