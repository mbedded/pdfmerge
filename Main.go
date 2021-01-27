package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type ArgumentMap struct {
	FrontFile          string
	BackFile           string
	OutputFile         string
	HasEnoughArguments bool
}

func main() {

	if IsPdftkInstalled() == false {
		log.Fatal("Error: pdftk is not installed.\nPlease install pdftk to use this tool.")
	}

	args := GetMappedArguments()
	if args.HasEnoughArguments == false {
		log.Fatal("3 arguments are required: front-file, back-file, output-file")
	}

	if AreInputFilesExisting(args) == false {
		log.Fatal("One or both input files could not be found.")
	}

	numFront := GetNumberOfPages(args.FrontFile)
	numBack := GetNumberOfPages(args.BackFile)

	if (numFront < numBack) || (numFront-numBack) > 1 {
		log.Fatalf("Front has %d pages. Back has %d pages.\nFront needs to have more pages than back and the different must not be bigger than 1.", numFront, numBack)
	}

	pageSum := numFront + numBack
	argsArray := make([]string, pageSum)

	pageTaker := 1
	for page := int64(1); page <= pageSum; page++ {
		// A = front, B = back
		var value string
		if page%2 == 1 {
			value = fmt.Sprintf("A%d", pageTaker)
		} else {
			value = fmt.Sprintf("B%d", pageTaker)
			pageTaker++
		}

		argsArray[page-1] = value
	}

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

	//dummy := out.String()
	//dummy_err := errorbuf.String()

	if err != nil {
		log.Fatal(err)
	}
}

func GetNumberOfPages(filename string) int64 {
	const prefix = "NumberOfPages:"

	var buffer bytes.Buffer
	cmd := exec.Command("pdftk", filename, "dump_data")
	cmd.Stdout = &buffer

	err := cmd.Run()

	if err != nil {
		log.Fatalf("Could not determine number if pages from '%q'", filename)
	}

	scanner := bufio.NewScanner(&buffer)

	var text string
	for scanner.Scan() {
		text = scanner.Text()

		if strings.HasPrefix(text, prefix) {
			text = strings.Replace(text, prefix, "", 1)
			break
		}
	}

	numberOfPages, err := strconv.ParseInt(strings.TrimSpace(text), 10, 0)

	if err != nil {
		log.Fatalf("Unable to determine number of pages for file %q", filename)
	}

	return numberOfPages
}

func IsPdftkInstalled() bool {
	err := exec.Command("pdftk", "-version").Run()
	return err == nil
}

func GetMappedArguments() ArgumentMap {
	hasEnoughArgs := len(os.Args) >= 4

	if hasEnoughArgs == false {
		return ArgumentMap{
			HasEnoughArguments: false,
		}
	} else {
		return ArgumentMap{
			FrontFile:          os.Args[1],
			BackFile:           os.Args[2],
			OutputFile:         os.Args[3],
			HasEnoughArguments: true,
		}
	}
}

func AreInputFilesExisting(args ArgumentMap) bool {
	_, errFront := os.Stat(args.FrontFile)
	_, errBack := os.Stat(args.BackFile)

	return errFront == nil && errBack == nil
}
