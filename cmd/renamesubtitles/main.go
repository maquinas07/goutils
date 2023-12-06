package main

import (
	"os"

	"github.com/maquinas07/golibs/getopt"
)

func main() {
	getopt.AddOption('e', "video-extension", &videoExtension, true, nil)
	getopt.AddOption('s', "sub-extension", &subExtension, true, nil)
	getopt.AddOption('l', "language", &language, true, nil)
	getopt.AddFlag('r', "reverse", &reverse)
	getopt.AddFlag('v', "verbose", &verbose)
	getopt.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		reportAndDie(err, -1)
	}

	if reverse {
		err = reverseRenameSubtitles(cwd)
	} else {
		err = renameSubtitles(cwd)
	}
	if err != nil {
		reportAndDie(err, -1)
	}
	os.Exit(0)
}
