package main

import (
	"fmt"
	"os"

	smartConfig "github.com/flw-cn/go-smartConfig"
	"github.com/flw-cn/playground"
)

type Option struct {
	Lang string `flag:"t|go|programming {language}"`
	File string `flag:"f||file {name}"`
}

func main() {
	var opt Option
	smartConfig.LoadConfig("play", "v0.1.0", &opt)

	if opt.File == "" {
		fmt.Printf("Usage: %s --help\n", os.Args[0])
		os.Exit(1)
	}

	var output string
	var err error

	switch opt.Lang {
	case "go":
		output, err = playground.PlayGoFile(opt.File)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			fmt.Fprintf(os.Stderr, "%s", output)
			return
		}
	}

	fmt.Print(output)
}
