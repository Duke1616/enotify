package wechat

type MarkdownMessage struct {
	Receivers
	MsgType string `json:"msgtype"`
	Content string `json:"content"`
}

func NewMarkdownMessage(builder Receivers, content string) *MarkdownMessage {
	return &MarkdownMessage{
		Receivers: builder,
		Content:   content,
	}
}

func (m *MarkdownMessage) ToJSON() (map[string]any, error) {
	j := map[string]any{
		"touser":  m.ToUser,
		"toparty": m.ToParty,
		"totag":   m.ToTag,
		"agentid": m.AgentId,
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": m.Content,
		},
	}

	return j, nil
}
