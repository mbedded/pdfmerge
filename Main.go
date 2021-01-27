package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"pdfmerge/helpers"
)

func main() {

	if helpers.IsPdftkInstalled() == false {
		log.Fatal("Error: pdftk is not installed.\nPlease install pdftk to use this tool.")
	}

	args := helpers.GetMappedArguments()
	if args.HasEnoughArguments == false {
		log.Fatal("3 arguments are required: front-file, back-file, output-file")
	}

	if helpers.AreInputFilesExisting(args) == false {
		log.Fatal("One or both input files could not be found.")
	}

	argsArray := helpers.GeneratePagePickArguments(args)

	fileA := fmt.Sprintf("A=%s", args.FrontFile)
	fileB := fmt.Sprintf("B=%s", args.BackFile)

	begin := []string{fileA, fileB, "cat"}
	end := []string{"output", args.OutputFile}

	params := append(begin, argsArray...)
	params = append(params, end...)

	cmd := exec.Command("pdftk", params...)

	var out bytes.Buffer
	var errorbuf bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &errorbuf

	err := cmd.Run()

	if err != nil {
		log.Fatalf("Error parsing pdf-file. %q", errorbuf.String())
	}

	log.Printf("File %q created successfully!", args.OutputFile)
}
