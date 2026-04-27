package main

import (
	"github.com/givensuman/gh-tag/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
