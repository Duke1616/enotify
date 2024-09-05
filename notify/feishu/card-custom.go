package feishu

import (
	"fmt"
	"github.com/Duke1616/enotify/template"
)

type cardCustom struct {
	tmpl *template.Template
	data map[string]interface{}
}

func (c *cardCustom) Builder() string {
	// 执行模板
	execute, err := c.tmpl.Execute("app", c.data)
	if err != nil {
		return ""
	}

	// 打印生成的 JSON 字符串
	fmt.Println(execute)

	return execute
}

func (c *cardCustom) MsgType() string {
	return "interactive"
}

func NewFeishuCustomCard(tmpl *template.Template, data map[string]interface{}) Content {
	return &cardCustom{
		tmpl: tmpl,
		data: data,
	}
}
