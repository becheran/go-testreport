package main

import (
	"go-testreport"
	"log"
	"os"
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
