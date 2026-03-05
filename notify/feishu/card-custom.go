package feishu

import (
	"encoding/json"

	"github.com/Duke1616/enotify/notify/feishu/message"
	"github.com/Duke1616/enotify/template"
	"github.com/gotomicro/ego/core/elog"
)

type cardCustom struct {
	tmpl   *template.Template
	name   string
	data   map[string]interface{}
	logger *elog.Component
}

func (c *cardCustom) Builder() (string, error) {
	// 执行模板
	execute, err := c.tmpl.Execute(c.name, c.data)
	if err != nil {
		return "", err
	}

	c.logger.Debug("feishu custom card",
		elog.String("name", c.name), elog.Any("execute", json.RawMessage(execute)))

	return execute, nil
}

func (c *cardCustom) MsgType() string {
	return "interactive"
}

func NewFeishuCustomCard(tmpl *template.Template, name string, data map[string]interface{}) message.Content {
	return &cardCustom{
		tmpl:   tmpl,
		name:   name,
		data:   data,
		logger: elog.DefaultLogger,
	}
}
