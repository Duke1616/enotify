package template

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"
)

type Field struct {
	IsShort bool
	Tag     string
	Content string
}

type FieldList struct {
	Fields []Field
}

func TestTemplate(t *testing.T) {
	tmpl, err := template.ParseFiles("./approval.tmpl")
	if err != nil {

	}

	// 创建一个包含多个休假申请的列表
	requestList := FieldList{
		Fields: []Field{
			{
				IsShort: false,
				Tag:     "lark_md",
				Content: "**时间：**\n2020-4-8 至 2020-4-10（共3天）",
			},
		},
	}

	var Title string = "您好"

	// 执行模板
	var out bytes.Buffer
	err = tmpl.Execute(&out, requestList)

	err = tmpl.Execute(&out, Title)

	fmt.Print(out)
}
