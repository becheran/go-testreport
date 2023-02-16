package main

import (
	"flag"
	"log"
	"os"

	"github.com/becheran/go-testreport"
	"github.com/becheran/go-testreport/internal/args"
)

func main() {
	var templateFile string
	var vars string
	flag.StringVar(&templateFile, "template", "", "Template file for the report generation. If not set, the default template will be applied")
	flag.StringVar(&vars, "vars", "", "Comma separated list of custom variables which can be used in the template. For example -vars=version:1.2.4,build:42")
	flag.Parse()

	tmp, err := testreport.GetTemplate(templateFile)
	if err != nil {
		log.Fatalf("Invalid template. %s", err)
	}

	result, err := testreport.ParseTestJson(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to parse test result %s", err)
	}

	result.Vars, err = args.ParseCommaSeparatedList(vars)
	if err != nil {
		log.Fatalf("Failed to parse variables. %s", err)
	}

	if err := testreport.CreateReport(result, os.Stdout, tmp); err != nil {
		log.Fatalf("Failed to create test report. %s", err)
	}

	for _, packRes := range result.PackageResult {
		if packRes.PackageResult == testreport.FTSFail {
			os.Exit(1)
		}
	}
}
