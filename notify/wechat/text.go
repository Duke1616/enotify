package wechat

type TextMessage struct {
	Receivers
	MsgType string `json:"msgtype"`
	Content string `json:"content"`
}

func NewTextMessage(builder Receivers, content string) *TextMessage {
	return &TextMessage{
		Receivers: builder,
		Content:   content,
	}
}

func (m *TextMessage) ToJSON() (map[string]any, error) {
	j := map[string]any{
		"touser":  m.ToUser,
		"toparty": m.ToParty,
		"totag":   m.ToTag,
		"agentid": m.AgentId,
		"msgtype": "text",
		"text": map[string]string{
			"content": m.Content,
		},
	}

	return j, nil
}
