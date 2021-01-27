package helpers

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

const prefix = "NumberOfPages:"

func AreInputFilesExisting(args ArgumentMap) bool {
	_, errFront := os.Stat(args.FrontFile)
	_, errBack := os.Stat(args.BackFile)

	return errFront == nil && errBack == nil
}

func GeneratePagePickArguments(args ArgumentMap) []string {

	numFront := getNumberOfPages(args.FrontFile)
	numBack := getNumberOfPages(args.BackFile)

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

	return argsArray
}

func getNumberOfPages(filename string) int64 {
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
