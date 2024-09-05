package feishu

type custom struct {
	Type    string
	Content string
}

func (c *custom) Builder() string {
	return c.Content
}

func (c *custom) MsgType() string {
	return c.Type
}

func NewFeishuCustom(msgType, content string) Content {
	return &custom{
		Type:    msgType,
		Content: content,
	}
}
