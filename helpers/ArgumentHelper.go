package helpers

import (
	"os"
	"os/exec"
)

type ArgumentMap struct {
	FrontFile          string
	BackFile           string
	OutputFile         string
	HasEnoughArguments bool
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
