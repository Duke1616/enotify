package feishu

import (
	"github.com/Duke1616/enotify/notify/feishu/message"
	"github.com/Duke1616/enotify/template"
)

type cardCustom struct {
	tmpl *template.Template
	name string
	data map[string]interface{}
}

func (c *cardCustom) Builder() (string, error) {
	// 执行模板
	execute, err := c.tmpl.Execute(c.name, c.data)
	if err != nil {
		return "", err
	}

	return execute, nil
}

func (c *cardCustom) MsgType() string {
	return "interactive"
}

func NewFeishuCustomCard(tmpl *template.Template, name string, data map[string]interface{}) message.Content {
	return &cardCustom{
		tmpl: tmpl,
		name: name,
		data: data,
	}
}
