package main

import (
	"fmt"
	"os"

	smartConfig "github.com/flw-cn/go-smartConfig"
	"github.com/flw-cn/playground"
)

type Option struct {
	Lang string `flag:"t|go|programming {language}"`
	File string `flag:"f||{file} name"`
	Code string `flag:"C||code pice {file}"`
}

func main() {
	var opt Option
	smartConfig.LoadConfig("play", "v0.1.0", &opt)

	if opt.Code == "" && opt.File == "" {
		fmt.Printf("Usage: %s --help\n", os.Args[0])
		os.Exit(1)
	}

	var output string
	var err error

	switch opt.Lang {
	case "go":
		if opt.File != "" {
			output, err = playground.PlayGoFile(opt.File)
		} else if opt.Code != "" {
			output, err = playground.PlayGoCode(opt.Code)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			pe, ok := err.(*playground.PlayError)
			if ok {
				fmt.Fprintf(os.Stderr, "%s", pe.Output())
			}
			return
		}
	}

	fmt.Print(output)
}
