package main

import (
	"embed"
	"fmt"
	"github.com/Duke1616/enotify/template"
	"io/fs"
)

//go:embed template/*
var templates embed.FS

func main() {
	// 打印嵌入的文件系统内容
	fs.WalkDir(templates, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	})

	// 调用 FromGlobs 函数
	paths := []string{"template/default.tmpl", "template/email.tmpl"}
	_, err := template.FromGlobs(paths)
	if err != nil {
	}
}
