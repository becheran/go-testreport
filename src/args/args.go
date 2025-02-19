package args

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Args struct {
	TemplateFile         string
	OutputStream         io.WriteCloser
	InputStream          io.ReadCloser
	EnvArgs              map[string]string
	NonZeroExitOnFailure bool
}

func ParseArgs(cmdArgs []string, fs *flag.FlagSet) (result Args, err error) {
	fs.Usage = func() {
		fmt.Printf("go-testreport [<options>]")
		flag.PrintDefaults()
	}

	var vars, inputFile, outputFile string
	fs.StringVar(&inputFile, "input", "", "Input json test result file. If not set, stdin will be used")
	fs.StringVar(&outputFile, "output", "", "Output result file. If not set, stdout will be used")
	fs.StringVar(&result.TemplateFile, "template", "", "Template file for the report generation. If not set, the default template will be applied")
	fs.StringVar(&vars, "vars", "", "Comma separated list of custom variables which can be used in the template. For example -vars=\"Title:Custom Title\"")

	if err := fs.Parse(cmdArgs[1:]); err != nil {
		return Args{}, err
	}
	if fs.NArg() != 0 {
		return Args{}, fmt.Errorf("unexpected arguments: %v", fs.Args())
	}

	if inputFile != "" {
		result.InputStream, err = os.Open(inputFile)
		if err != nil {
			return Args{}, fmt.Errorf("failed to open input file %s. %s", inputFile, err)
		}
		result.NonZeroExitOnFailure = true
	} else {
		result.InputStream = os.Stdin
	}

	if outputFile != "" {
		result.OutputStream, err = os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			result.InputStream.Close()
			return Args{}, fmt.Errorf("failed to open output file %s. %s", outputFile, err)
		}
	} else {
		result.OutputStream = os.Stdout
	}

	result.EnvArgs, err = parseCommaSeparatedList(vars)
	if err != nil {
		result.InputStream.Close()
		result.OutputStream.Close()
		return Args{}, err
	}

	return result, nil
}

func parseCommaSeparatedList(input string) (result map[string]string, err error) {
	result = make(map[string]string)
	args := strings.Split(input, ",")
	for _, arg := range args {
		if arg == "" {
			continue
		}
		varVal := strings.Split(arg, ":")
		if len(varVal) != 2 {
			return nil, fmt.Errorf("expected variable and value separated with an ':' sign")
		}
		if varVal[0] == "" {
			return nil, fmt.Errorf("key must not be empty")
		}
		if _, ok := result[varVal[0]]; ok {
			return nil, fmt.Errorf("variable %s can only be assigned once", varVal[0])
		}
		result[varVal[0]] = varVal[1]
	}
	return
}
