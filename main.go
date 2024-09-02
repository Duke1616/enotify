package main

import (
	"enotify/template"
	"fmt"
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
