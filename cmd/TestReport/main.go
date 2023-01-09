package main

import (
	"log"
	"os"

	"github.com/becheran/go-testreport"
)

func main() {
	md, err := testreport.CreateReport(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to create test report. %s", err)
	}
	_, err = os.Stdout.Write(md)
	if err != nil {
		log.Fatalf("Failed to write markdown result file. %s", err)
	}
}
