package feishu

import "github.com/Duke1616/enotify/notify/feishu/message"

type custom struct {
	Type    string
	Content string
}

func (c *custom) Builder() (string, error) {
	return c.Content, nil
}

func (c *custom) MsgType() string {
	return c.Type
}

func NewFeishuCustom(msgType, content string) message.Content {
	return &custom{
		Type:    msgType,
		Content: content,
	}
}
