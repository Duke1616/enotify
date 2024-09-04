package main

import (
	"fmt"
	"github.com/Duke1616/enotify/template"
)

func main() {
	t()
}

func t() {
	globs, err := template.FromGlobs([]string{})
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print(globs)
}
