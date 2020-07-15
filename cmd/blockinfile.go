package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/czerasz/go-blockinfile/blockinfile"
)

var version = "v0.0.0"
var commit = "4fe5445d34637a0ace412b154eaad840686e684a"
var date = "now"
var builtBy = "local user"

const program = "blockinfile"

func customUsage() {
	fmt.Printf("Usage of %s:\n", program)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = customUsage

	pathFlag := flag.String("path", "", "The file to modify. (Required)")
	contentFlag := flag.String("content", "", "The text to insert inside the marker lines. (Required)")
	markerFlag := flag.String("marker", blockinfile.DefaultMarkerTemplate, "The marker line template.")
	versionFlag := flag.Bool("version", false, "Display version.")

	flag.Parse()

	if *versionFlag {
		osArch := runtime.GOOS + "/" + runtime.GOARCH

		fmt.Printf("%s %s (%s) %s BuildDate: %s\n BuiltBy: %s", program, version, commit, osArch, date, builtBy)
		os.Exit(0)
	}

	if *pathFlag == "" || *contentFlag == "" {
		fmt.Printf("Required arguments not provided\n\n")
		flag.Usage()
		os.Exit(1)
	}

	var content []byte

	if *contentFlag == "-" {
		info, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}

		if info.Mode()&os.ModeNamedPipe == 0 {
			fmt.Printf("The command is intended to work with pipes.\n\n")
			fmt.Printf("Usage: echo 'example' | %s -path %s -content -\n", program, *pathFlag)
			os.Exit(1)
		}

		content, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
	} else {
		content = []byte(*contentFlag)
	}

	err := blockinfile.Update(*pathFlag, *markerFlag, content)
	if err != nil {
		panic(err)
	}
}
