package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/becheran/go-testreport/src/args"
	"github.com/becheran/go-testreport/src/report"
)

func main() {
	args, err := args.ParseArgs(os.Args, flag.CommandLine)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}
	defer args.OutputStream.Close()
	defer args.InputStream.Close()

	tmp, err := report.GetTemplate(args.TemplateFile)
	if err != nil {
		log.Fatalf("Invalid template. %s", err)
	}

	result, err := report.ParseTestJson(args.InputStream)
	if err != nil {
		log.Fatalf("Failed to parse test result %s", err)
	}

	result.Vars = args.EnvArgs

	if err := report.CreateReport(result, args.OutputStream, tmp); err != nil {
		log.Fatalf("Failed to create test report. %s", err)
	}

	failed := false
	for _, packRes := range result.PackageResult {
		fmt.Println(packRes)
		if packRes.PackageResult == report.FTSFail {
			failed = true
		}
	}

	if args.NonZeroExitOnFailure && failed {
		os.Exit(1)
	}
}
